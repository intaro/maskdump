package main

import (
	"fmt"
	"strings"
)

// DumpDialect identifies the SQL dump format being processed.
type DumpDialect string

const (
	// DialectAuto selects the dialect by content-based detection.
	DialectAuto DumpDialect = "auto"
	// DialectMySQL handles MySQL / MariaDB dumps.
	DialectMySQL DumpDialect = "mysql"
	// DialectPostgreSQL handles PostgreSQL dumps (COPY and INSERT styles).
	DialectPostgreSQL DumpDialect = "postgresql"
	// DialectOracle handles Oracle dumps.
	DialectOracle DumpDialect = "oracle"
	// DialectMSSQL handles MS SQL Server dumps.
	DialectMSSQL DumpDialect = "mssql"
	// DialectSQLite handles SQLite dumps.
	DialectSQLite DumpDialect = "sqlite"
	// DialectFirebird handles Firebird dumps.
	DialectFirebird DumpDialect = "firebird"
	// DialectGeneric is the internal fallback when no dialect-specific
	// parsing is available: full-line regex masking only.
	DialectGeneric DumpDialect = "generic"
)

// ParseDumpDialect validates a user-supplied db-format value.
func ParseDumpDialect(value string) (DumpDialect, error) {
	switch DumpDialect(strings.ToLower(strings.TrimSpace(value))) {
	case "", DialectAuto:
		return DialectAuto, nil
	case DialectMySQL:
		return DialectMySQL, nil
	case DialectPostgreSQL:
		return DialectPostgreSQL, nil
	case DialectOracle:
		return DialectOracle, nil
	case DialectMSSQL:
		return DialectMSSQL, nil
	case DialectSQLite:
		return DialectSQLite, nil
	case DialectFirebird:
		return DialectFirebird, nil
	default:
		return "", fmt.Errorf("unsupported db-format %q (expected auto|mysql|postgresql|oracle|mssql|sqlite|firebird)", value)
	}
}

// DialectParser processes dump lines for one SQL dialect.
//
// ProcessLine receives one input line (with its trailing newline when
// present) and returns the transformed output plus a drop flag. When drop is
// true the line is removed from the output entirely. The returned string may
// contain several lines when the parser flushes buffered input.
type DialectParser interface {
	Dialect() DumpDialect
	ProcessLine(line string, config MaskConfig, cache *Cache) (string, bool)
}

// NewDialectParser builds the parser implementation for the requested dialect.
func NewDialectParser(dialect DumpDialect, rt *Runtime) DialectParser {
	switch dialect {
	case DialectMySQL:
		return newMySQLDialectParser(rt)
	case DialectPostgreSQL:
		return newPostgresDialectParser(rt)
	case DialectOracle:
		return newSQLInsertDialectParser(rt, DialectOracle)
	case DialectMSSQL:
		return newSQLInsertDialectParser(rt, DialectMSSQL)
	case DialectSQLite:
		return newSQLInsertDialectParser(rt, DialectSQLite)
	case DialectFirebird:
		return newSQLInsertDialectParser(rt, DialectFirebird)
	case DialectAuto:
		return newDetectingDialectParser(rt)
	default:
		return newGenericDialectParser(rt)
	}
}

// genericDialectParser is the safe fallback: it never skips lines and never
// guesses field positions; it only applies full-line regex masking, matching
// the historical behavior for dumps without table awareness.
type genericDialectParser struct {
	rt     *Runtime
	warned bool
}

func newGenericDialectParser(rt *Runtime) *genericDialectParser {
	return &genericDialectParser{rt: rt}
}

// Dialect implements DialectParser.
func (p *genericDialectParser) Dialect() DumpDialect { return DialectGeneric }

// ProcessLine implements DialectParser.
func (p *genericDialectParser) ProcessLine(line string, config MaskConfig, cache *Cache) (string, bool) {
	if !p.warned {
		p.warned = true
		if (len(p.rt.SkipTableList) > 0 || len(p.rt.NoMaskTableList) > 0 || len(p.rt.ProcessingTables) > 0) && logger != nil {
			logger.Warn("selective table filtering is disabled: dump dialect is unknown, applying full-line masking only")
		}
	}
	return maskFullLine(p.rt, line, config, cache), false
}

// maskFullLine applies the configured regex masking to a whole line without
// any table or field awareness.
func maskFullLine(rt *Runtime, line string, config MaskConfig, cache *Cache) string {
	if config.emailAlgorithm == "light-hash" && rt.EmailRegex != nil {
		line = rt.EmailRegex.ReplaceAllStringFunc(line, func(email string) string {
			return rt.MaskEmailWithRules(email, cache)
		})
	}
	if config.phoneAlgorithm == "light-mask" && rt.PhoneRegex != nil {
		line = rt.PhoneRegex.ReplaceAllStringFunc(line, func(phone string) string {
			return rt.MaskPhoneWithRules(phone, cache)
		})
	}
	return line
}

// normalizeIdentifier strips dialect quoting from one SQL identifier:
// `name` (MySQL), "name" (standard), [name] (MSSQL) and surrounding spaces.
func normalizeIdentifier(raw string) string {
	s := strings.TrimSpace(raw)
	for len(s) >= 2 {
		first := s[0]
		last := s[len(s)-1]
		if (first == '`' && last == '`') ||
			(first == '"' && last == '"') ||
			(first == '[' && last == ']') {
			s = s[1 : len(s)-1]
			continue
		}
		break
	}
	return s
}

// normalizeTableName strips quoting from a possibly schema-qualified table
// reference and returns the normalized full name plus the plain table name
// without the schema prefix ("public.users" -> "public.users", "users").
func normalizeTableName(raw string) (full string, plain string) {
	parts := strings.Split(strings.TrimSpace(raw), ".")
	normalized := make([]string, 0, len(parts))
	for _, part := range parts {
		normalized = append(normalized, normalizeIdentifier(part))
	}
	full = strings.Join(normalized, ".")
	plain = normalized[len(normalized)-1]
	return full, plain
}

// tableNameCandidates returns config lookup keys for a table reference:
// the schema-qualified form first, then the plain name.
func tableNameCandidates(raw string) []string {
	full, plain := normalizeTableName(raw)
	if full == plain {
		return []string{plain}
	}
	return []string{full, plain}
}

// lookupProcessingTable finds the masking config for a table reference,
// accepting both schema-qualified and plain config keys. With fold the match
// is case-insensitive (Oracle folds unquoted identifiers to upper case, so
// dump and config casing routinely differ).
func lookupProcessingTable(rt *Runtime, rawTable string, fold bool) (TableConfig, bool) {
	for _, name := range tableNameCandidates(rawTable) {
		if cfg, ok := rt.ProcessingTables[name]; ok {
			return cfg, true
		}
		if fold {
			for key, cfg := range rt.ProcessingTables {
				if strings.EqualFold(key, name) {
					return cfg, true
				}
			}
		}
	}
	return TableConfig{}, false
}

// tableInList reports whether a table reference matches a configured table
// list, accepting both schema-qualified and plain config keys.
func tableInList(list map[string]struct{}, rawTable string, fold bool) bool {
	for _, name := range tableNameCandidates(rawTable) {
		if _, ok := list[name]; ok {
			return true
		}
		if fold {
			for key := range list {
				if strings.EqualFold(key, name) {
					return true
				}
			}
		}
	}
	return false
}

// isSkippedTable reports whether a table's data rows must be dropped.
func isSkippedTable(rt *Runtime, rawTable string, fold bool) bool {
	return tableInList(rt.SkipTableList, rawTable, fold)
}

// isNoMaskTable reports whether a table's data rows must pass through
// without any masking.
func isNoMaskTable(rt *Runtime, rawTable string, fold bool) bool {
	return tableInList(rt.NoMaskTableList, rawTable, fold)
}

// fieldPositions resolves configured email/phone column names to 0-based
// positions using an ordered column list. Unknown names are ignored. With
// fold the column names match case-insensitively.
func fieldPositions(tableConfig TableConfig, columns []string, config MaskConfig, fold bool) (emailPos, phonePos map[int]bool) {
	emailPos = make(map[int]bool)
	phonePos = make(map[int]bool)

	key := func(name string) string {
		if fold {
			return strings.ToLower(name)
		}
		return name
	}

	index := make(map[string]int, len(columns))
	for i, col := range columns {
		index[key(normalizeIdentifier(col))] = i
	}

	if config.emailAlgorithm == "light-hash" {
		for _, name := range tableConfig.Email {
			if i, ok := index[key(name)]; ok {
				emailPos[i] = true
			}
		}
	}
	if config.phoneAlgorithm == "light-mask" {
		for _, name := range tableConfig.Phone {
			if i, ok := index[key(name)]; ok {
				phonePos[i] = true
			}
		}
	}
	return emailPos, phonePos
}

// maskValueAt applies email/phone masking to a single raw SQL value by its
// column position. It returns the possibly modified value.
func maskValueAt(rt *Runtime, value string, pos int, emailPos, phonePos map[int]bool, cache *Cache) string {
	if value == "" || value == "NULL" {
		return value
	}
	if emailPos[pos] && rt.EmailRegex != nil {
		value = rt.EmailRegex.ReplaceAllStringFunc(value, func(email string) string {
			return rt.MaskEmailWithRules(email, cache)
		})
	}
	if phonePos[pos] && rt.PhoneRegex != nil {
		value = rt.PhoneRegex.ReplaceAllStringFunc(value, func(phone string) string {
			return rt.MaskPhoneWithRules(phone, cache)
		})
	}
	return value
}

// maskTuples masks configured columns inside every (...) tuple found in s.
func maskTuples(rt *Runtime, s string, emailPos, phonePos map[int]bool, cache *Cache) string {
	if len(emailPos) == 0 && len(phonePos) == 0 {
		return s
	}
	return tupleRegex.ReplaceAllStringFunc(s, func(tuple string) string {
		values := parseTuple(tuple)
		if len(values) == 0 {
			return tuple
		}
		modified := false
		for pos := range values {
			masked := maskValueAt(rt, values[pos], pos, emailPos, phonePos, cache)
			if masked != values[pos] {
				values[pos] = masked
				modified = true
			}
		}
		if !modified {
			return tuple
		}
		return "(" + strings.Join(values, ",") + ")"
	})
}
