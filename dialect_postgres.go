package main

import (
	"regexp"
	"strings"
)

// COPY "schema"."table" (col1, col2, ...) FROM stdin;
var pgCopyRegex = regexp.MustCompile(`(?i)^COPY\s+([^\s(]+)\s*\(([^)]*)\)\s+FROM\s+stdin;\s*$`)

const pgCopyTerminator = `\.`

// postgresDialectParser handles pg_dump output in both COPY and INSERT
// styles. INSERT statements and CREATE TABLE parsing are shared with the
// generic SQL machinery; COPY blocks are handled here.
type postgresDialectParser struct {
	rt   *Runtime
	proc *sqlStatementProcessor

	// open COPY block state
	copyActive bool
	copyDrop   bool
	copyNoMask bool
	copyEmail  map[int]bool
	copyPhone  map[int]bool
}

func newPostgresDialectParser(rt *Runtime) *postgresDialectParser {
	return &postgresDialectParser{
		rt:   rt,
		proc: newSQLStatementProcessor(rt, false),
	}
}

// Dialect implements DialectParser.
func (p *postgresDialectParser) Dialect() DumpDialect { return DialectPostgreSQL }

// ProcessLine implements DialectParser.
func (p *postgresDialectParser) ProcessLine(line string, config MaskConfig, cache *Cache) (string, bool) {
	selective := len(p.rt.ProcessingTables) > 0
	filtering := selective || len(p.rt.SkipTableList) > 0 || len(p.rt.NoMaskTableList) > 0
	body, newline := splitTrailingNewline(line)

	// Rows inside an open COPY block.
	if p.copyActive {
		if body == pgCopyTerminator {
			drop := p.copyDrop
			p.resetCopy()
			return line, drop
		}
		if p.copyDrop {
			return "", true
		}
		if p.copyNoMask {
			return line, false
		}
		if len(p.copyEmail) > 0 || len(p.copyPhone) > 0 {
			return p.maskCopyRow(body, cache) + newline, false
		}
		if selective {
			// Selective mode masks only configured fields.
			return line, false
		}
		return maskFullLine(p.rt, body, config, cache) + newline, false
	}

	if filtering {
		if matches := pgCopyRegex.FindStringSubmatch(body); matches != nil {
			return p.startCopyBlock(matches[1], matches[2], line, config)
		}
	}

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

// startCopyBlock opens a COPY ... FROM stdin block and decides how its rows
// will be treated: dropped, masked by column position, or passed through.
func (p *postgresDialectParser) startCopyBlock(table, columnList, line string, config MaskConfig) (string, bool) {
	p.copyActive = true

	if isSkippedTable(p.rt, table, p.proc.fold) {
		p.copyDrop = true
		return "", true
	}

	if isNoMaskTable(p.rt, table, p.proc.fold) {
		p.copyNoMask = true
		return line, false
	}

	if tableConfig, ok := lookupProcessingTable(p.rt, table, p.proc.fold); ok {
		columns := splitColumnList(columnList)
		if len(columns) == 0 {
			// Safety rule: no confident column positions, no masking.
			if logger != nil {
				logger.Warn("cannot parse COPY column list for table %s: rows pass through unmasked", table)
			}
		} else {
			p.copyEmail, p.copyPhone = fieldPositions(tableConfig, columns, config, p.proc.fold)
		}
	}
	return line, false
}

// maskCopyRow masks configured columns of one tab-separated COPY data row.
// Literal tabs inside values are escaped as "\t" by pg_dump, so splitting on
// the tab character is unambiguous. "\N" marks NULL and is left untouched.
func (p *postgresDialectParser) maskCopyRow(body string, cache *Cache) string {
	values := strings.Split(body, "\t")
	for pos := range values {
		if values[pos] == `\N` {
			continue
		}
		values[pos] = maskValueAt(p.rt, values[pos], pos, p.copyEmail, p.copyPhone, cache)
	}
	return strings.Join(values, "\t")
}

func (p *postgresDialectParser) resetCopy() {
	p.copyActive = false
	p.copyDrop = false
	p.copyNoMask = false
	p.copyEmail = nil
	p.copyPhone = nil
}
