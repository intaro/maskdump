package main

import "strings"

// mysqlDialectParser wraps the historical MySQL-specific processing:
// backtick-quoted CREATE TABLE / INSERT parsing, skip list matching by
// "INSERT INTO `table`" prefix and field-aware masking via TableParser.
type mysqlDialectParser struct {
	rt     *Runtime
	tables *TableParser
}

func newMySQLDialectParser(rt *Runtime) *mysqlDialectParser {
	return &mysqlDialectParser{
		rt:     rt,
		tables: NewTableParser(rt),
	}
}

// Dialect implements DialectParser.
func (p *mysqlDialectParser) Dialect() DumpDialect { return DialectMySQL }

// ProcessLine implements DialectParser.
func (p *mysqlDialectParser) ProcessLine(line string, config MaskConfig, cache *Cache) (string, bool) {
	if len(p.rt.SkipTableList) > 0 {
		for table := range p.rt.SkipTableList {
			if strings.HasPrefix(line, "INSERT INTO `"+table+"`") {
				return "", true
			}
		}
	}

	if len(p.rt.ProcessingTables) > 0 {
		p.tables.ParseTableStructure(line)
		body, newline := splitTrailingNewline(line)
		return p.tables.ProcessDumpLine(body, config, cache) + newline, false
	}

	return maskFullLine(p.rt, line, config, cache), false
}
