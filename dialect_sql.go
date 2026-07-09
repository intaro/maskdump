package main

import (
	"regexp"
	"strings"
)

// Shared statement-level parsing for dialects whose data lines are standard
// SQL INSERT statements: PostgreSQL (--inserts), Oracle, MS SQL Server,
// SQLite and Firebird. Identifiers may be bare, double-quoted or bracketed.
var (
	// INSERT INTO <table> [(col, ...)] VALUES <rest>
	sqlInsertRegex = regexp.MustCompile(`(?i)^\s*INSERT\s+INTO\s+([` + "`" + `"\[\]\w$.]+)\s*(?:\(([^)]*)\)\s*)?VALUES\s*(.*)$`)
	// CREATE TABLE [IF NOT EXISTS] <table> [(...]
	sqlCreateTableRegex = regexp.MustCompile(`(?i)^\s*CREATE\s+TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?([` + "`" + `"\[\]\w$.]+)\s*(\(?)(.*)$`)
	// A continuation line of a multi-row VALUES list: "(...)," or "(...);"
	sqlTupleLineRegex = regexp.MustCompile(`^\s*\(.*\)\s*[,;]?\s*$`)
	// First identifier of a column definition line inside CREATE TABLE.
	sqlColumnDefRegex = regexp.MustCompile("^\\s*([`\"\\[]?)([\\w$]+)[`\"\\]]?\\s+")
)

// sqlConstraintKeywords are line starters inside CREATE TABLE that do not
// declare a column.
var sqlConstraintKeywords = map[string]struct{}{
	"PRIMARY": {}, "UNIQUE": {}, "KEY": {}, "CONSTRAINT": {}, "FOREIGN": {},
	"CHECK": {}, "INDEX": {}, "FULLTEXT": {}, "SPATIAL": {}, "EXCLUDE": {},
	"LIKE": {}, "PERIOD": {},
}

// sqlStatementProcessor holds the streaming state shared by SQL-insert
// dialects: collected table structures and the currently open multi-line
// INSERT or CREATE TABLE statement.
type sqlStatementProcessor struct {
	rt *Runtime
	// fold enables case-insensitive identifier matching (Oracle folds
	// unquoted identifiers to upper case).
	fold bool
	// tables collects column order per table (normalized full and plain
	// names both point at the same entry).
	tables map[string][]string

	// open CREATE TABLE statement state
	creatingTable   string
	creatingColumns []string

	// open INSERT statement state
	insertActive bool
	insertDrop   bool
	insertNoMask bool
	emailPos     map[int]bool
	phonePos     map[int]bool
}

func newSQLStatementProcessor(rt *Runtime, fold bool) *sqlStatementProcessor {
	return &sqlStatementProcessor{
		rt:     rt,
		fold:   fold,
		tables: make(map[string][]string),
	}
}

// tableKey normalizes a table map key according to the fold mode.
func (p *sqlStatementProcessor) tableKey(name string) string {
	if p.fold {
		return strings.ToLower(name)
	}
	return name
}

// rememberTable stores the column order for both schema-qualified and plain
// table names.
func (p *sqlStatementProcessor) rememberTable(rawTable string, columns []string) {
	if len(columns) == 0 {
		return
	}
	full, plain := normalizeTableName(rawTable)
	p.tables[p.tableKey(full)] = columns
	p.tables[p.tableKey(plain)] = columns
}

// columnsFor returns the known column order for a table reference.
func (p *sqlStatementProcessor) columnsFor(rawTable string) ([]string, bool) {
	for _, name := range tableNameCandidates(rawTable) {
		if cols, ok := p.tables[p.tableKey(name)]; ok {
			return cols, true
		}
	}
	return nil, false
}

// splitColumnList splits a "col1, col2, ..." list and normalizes identifiers.
func splitColumnList(list string) []string {
	parts := strings.Split(list, ",")
	columns := make([]string, 0, len(parts))
	for _, part := range parts {
		name := normalizeIdentifier(part)
		if name == "" {
			return nil
		}
		columns = append(columns, name)
	}
	return columns
}

// statementTerminated reports whether an INSERT VALUES fragment ends the
// statement on this line.
func statementTerminated(rest string) bool {
	trimmed := strings.TrimRight(rest, " \t\r\n")
	return strings.HasSuffix(trimmed, ";")
}

// processCreateTableLine consumes DDL lines. It returns true when the line
// was part of a CREATE TABLE statement.
func (p *sqlStatementProcessor) processCreateTableLine(line string) bool {
	if p.creatingTable == "" {
		matches := sqlCreateTableRegex.FindStringSubmatch(line)
		if matches == nil {
			return false
		}
		table := matches[1]
		body := matches[3]
		if matches[2] == "" && !strings.HasPrefix(strings.TrimSpace(body), "(") {
			// CREATE TABLE without a column list on this line: the "("
			// may follow on the next line; treat conservatively as a
			// one-line statement we cannot use.
			if !strings.Contains(body, "(") {
				return true
			}
			body = body[strings.Index(body, "(")+1:]
		}
		// Single-line CREATE TABLE "t" (a int, b text);
		if strings.Contains(body, ")") && strings.HasSuffix(strings.TrimRight(body, " \t\r\n"), ";") {
			p.rememberTable(table, columnsFromDefinitionList(body[:strings.LastIndex(body, ")")]))
			return true
		}
		p.creatingTable = table
		p.creatingColumns = nil
		if rest := strings.TrimSpace(body); rest != "" {
			if col, ok := columnFromDefinitionLine(rest); ok {
				p.creatingColumns = append(p.creatingColumns, col)
			}
		}
		return true
	}

	// Inside an open CREATE TABLE block. A ")" line (with or without a
	// trailing ";") closes the definition; MSSQL emits it without ";".
	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, ")") || endTableRegex.MatchString(line) {
		p.rememberTable(p.creatingTable, p.creatingColumns)
		p.creatingTable = ""
		p.creatingColumns = nil
		return true
	}
	// An INSERT reaching this point means the CREATE TABLE block never
	// closed as expected: abandon DDL state instead of consuming data.
	if sqlInsertRegex.MatchString(trimmed) {
		p.creatingTable = ""
		p.creatingColumns = nil
		return false
	}
	if col, ok := columnFromDefinitionLine(line); ok {
		p.creatingColumns = append(p.creatingColumns, col)
	}
	return true
}

// columnsFromDefinitionList extracts column names from a single-line
// "a int, b text, PRIMARY KEY (a)" definition body.
func columnsFromDefinitionList(body string) []string {
	var columns []string
	depth := 0
	start := 0
	flush := func(end int) {
		part := body[start:end]
		if col, ok := columnFromDefinitionLine(part); ok {
			columns = append(columns, col)
		}
	}
	for i, c := range body {
		switch c {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				flush(i)
				start = i + 1
			}
		}
	}
	flush(len(body))
	return columns
}

// columnFromDefinitionLine extracts the column name from one column
// definition fragment, rejecting constraint clauses.
func columnFromDefinitionLine(line string) (string, bool) {
	matches := sqlColumnDefRegex.FindStringSubmatch(line)
	if matches == nil {
		return "", false
	}
	name := matches[2]
	if matches[1] == "" {
		if _, ok := sqlConstraintKeywords[strings.ToUpper(name)]; ok {
			return "", false
		}
	}
	return name, true
}

// insertAction tells the dialect parser what to do with a line examined by
// the INSERT machinery.
type insertAction int

const (
	// insertNotHandled marks a line that is not part of an INSERT statement.
	insertNotHandled insertAction = iota
	// insertHandled marks INSERT output that full-line masking may still be
	// applied to when selective mode is off.
	insertHandled
	// insertHandledRaw marks output of a no-mask table: never mask it.
	insertHandledRaw
	// insertDropped marks a line that is removed from the output.
	insertDropped
)

// processInsertLine handles INSERT statements, multi-line VALUES lists and
// skip/no-mask-listed tables. It returns the transformed line and the action
// the caller must take.
func (p *sqlStatementProcessor) processInsertLine(line string, config MaskConfig, cache *Cache) (string, insertAction) {
	// Continuation of an open multi-line VALUES list.
	if p.insertActive {
		if sqlTupleLineRegex.MatchString(line) {
			drop, noMask := p.insertDrop, p.insertNoMask
			emailPos, phonePos := p.emailPos, p.phonePos
			if statementTerminated(line) {
				p.resetInsert()
			}
			switch {
			case drop:
				return "", insertDropped
			case noMask:
				return line, insertHandledRaw
			case len(emailPos) == 0 && len(phonePos) == 0:
				return line, insertHandled
			default:
				return maskTuples(p.rt, line, emailPos, phonePos, cache), insertHandled
			}
		}
		// The line does not look like a tuple: the statement ended
		// implicitly. Safety rule: never consume lines we are not sure
		// about; fall through to regular processing.
		p.resetInsert()
	}

	matches := sqlInsertRegex.FindStringSubmatch(line)
	if matches == nil {
		return line, insertNotHandled
	}
	table := matches[1]
	columnList := matches[2]
	rest := matches[3]

	// Multi-line statements keep state until the closing ";".
	multiLine := !statementTerminated(rest)

	if isSkippedTable(p.rt, table, p.fold) {
		if multiLine {
			p.insertActive = true
			p.insertDrop = true
		}
		return "", insertDropped
	}

	if isNoMaskTable(p.rt, table, p.fold) {
		if multiLine {
			p.insertActive = true
			p.insertNoMask = true
		}
		return line, insertHandledRaw
	}

	tableConfig, ok := lookupProcessingTable(p.rt, table, p.fold)
	if !ok {
		if multiLine && strings.TrimSpace(rest) == "" {
			// Open multi-line VALUES list of an unconfigured table: keep
			// passing tuple lines through with no field awareness.
			p.insertActive = true
			p.insertDrop = false
			p.emailPos = nil
			p.phonePos = nil
		}
		return line, insertHandled
	}

	var columns []string
	if strings.TrimSpace(columnList) != "" {
		columns = splitColumnList(columnList)
	} else {
		columns, _ = p.columnsFor(table)
	}
	if len(columns) == 0 {
		// Safety rule: without confident column positions no field-aware
		// masking is applied.
		if logger != nil {
			logger.Warn("no column information for table %s: leaving INSERT unmasked", table)
		}
		return line, insertHandled
	}

	emailPos, phonePos := fieldPositions(tableConfig, columns, config, p.fold)
	if multiLine {
		p.insertActive = true
		p.insertDrop = false
		p.emailPos = emailPos
		p.phonePos = phonePos
	}
	if strings.TrimSpace(rest) == "" {
		return line, insertHandled
	}
	masked := maskTuples(p.rt, rest, emailPos, phonePos, cache)
	if masked == rest {
		return line, insertHandled
	}
	return line[:len(line)-len(rest)] + masked, insertHandled
}

// resetInsert clears the multi-line INSERT statement state.
func (p *sqlStatementProcessor) resetInsert() {
	p.insertActive = false
	p.insertDrop = false
	p.insertNoMask = false
	p.emailPos = nil
	p.phonePos = nil
}

// sqlInsertDialectParser adapts sqlStatementProcessor to the DialectParser
// interface for Oracle, MSSQL, SQLite and Firebird dumps.
type sqlInsertDialectParser struct {
	dialect DumpDialect
	rt      *Runtime
	proc    *sqlStatementProcessor
}

func newSQLInsertDialectParser(rt *Runtime, dialect DumpDialect) *sqlInsertDialectParser {
	return &sqlInsertDialectParser{
		dialect: dialect,
		rt:      rt,
		proc:    newSQLStatementProcessor(rt, dialect == DialectOracle),
	}
}

// Dialect implements DialectParser.
func (p *sqlInsertDialectParser) Dialect() DumpDialect { return p.dialect }

// ProcessLine implements DialectParser.
func (p *sqlInsertDialectParser) ProcessLine(line string, config MaskConfig, cache *Cache) (string, bool) {
	selective := len(p.rt.ProcessingTables) > 0
	filtering := selective || len(p.rt.SkipTableList) > 0 || len(p.rt.NoMaskTableList) > 0
	body, newline := splitTrailingNewline(line)

	if selective && !p.proc.insertActive && p.proc.processCreateTableLine(body) {
		return line, false
	}

	if filtering {
		out, action := p.proc.processInsertLine(body, config, cache)
		switch action {
		case insertDropped:
			return "", true
		case insertHandledRaw:
			return out + newline, false
		case insertHandled:
			if selective {
				// Selective mode masks only configured fields.
				return out + newline, false
			}
			return maskFullLine(p.rt, out, config, cache) + newline, false
		}
	}

	if selective {
		// Selective mode masks only configured fields; unrelated lines
		// pass through unchanged.
		return line, false
	}
	return maskFullLine(p.rt, line, config, cache), false
}

// splitTrailingNewline separates the line body from its trailing newline so
// anchored regexes can match the body. The newline is re-appended on output.
func splitTrailingNewline(line string) (body, newline string) {
	if strings.HasSuffix(line, "\r\n") {
		return line[:len(line)-2], "\r\n"
	}
	if strings.HasSuffix(line, "\n") {
		return line[:len(line)-1], "\n"
	}
	return line, ""
}
