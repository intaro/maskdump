package main

import (
	"regexp"
	"strings"
)

// detectMaxBufferedLines bounds how much input the auto-detector may buffer
// before giving up and falling back, keeping the pipeline streaming-friendly.
const detectMaxBufferedLines = 500

// Decisive single-line markers per dialect. The first match wins.
var dialectMarkers = []struct {
	dialect DumpDialect
	pattern *regexp.Regexp
}{
	{DialectMySQL, regexp.MustCompile(`^-- MySQL dump|^/\*!\d{5}|INSERT INTO ` + "`" + `|CREATE TABLE ` + "`")},
	{DialectPostgreSQL, regexp.MustCompile(`^-- PostgreSQL database dump|FROM stdin;\s*$|^SET statement_timeout|^SELECT pg_catalog\.|^\\connect\s`)},
	{DialectMSSQL, regexp.MustCompile(`^USE \[|^SET ANSI_NULLS|IDENTITY_INSERT|\[dbo\]\.`)},
	{DialectSQLite, regexp.MustCompile(`^PRAGMA foreign_keys\s*=|sqlite_sequence|^-- SQLite dump`)},
	{DialectOracle, regexp.MustCompile(`^-- Oracle Database dump|^SET DEFINE OFF|EXECUTE IMMEDIATE|VARCHAR2\(`)},
	{DialectFirebird, regexp.MustCompile(`^SET TERM\s|RDB\$|^-- Firebird`)},
}

// detectingDialectParser implements --db-format=auto: it buffers input lines
// until a decisive dialect marker appears, then replays the buffer through
// the selected parser and delegates the rest of the stream to it.
type detectingDialectParser struct {
	rt       *Runtime
	buffer   []string
	delegate DialectParser
}

func newDetectingDialectParser(rt *Runtime) *detectingDialectParser {
	return &detectingDialectParser{rt: rt}
}

// Dialect implements DialectParser.
func (p *detectingDialectParser) Dialect() DumpDialect {
	if p.delegate != nil {
		return p.delegate.Dialect()
	}
	return DialectAuto
}

// ProcessLine implements DialectParser.
func (p *detectingDialectParser) ProcessLine(line string, config MaskConfig, cache *Cache) (string, bool) {
	if p.delegate != nil {
		return p.delegate.ProcessLine(line, config, cache)
	}

	dialect, ok := detectDialectLine(line)
	if !ok {
		p.buffer = append(p.buffer, line)
		if len(p.buffer) < detectMaxBufferedLines {
			return "", true // buffered, emitted later on flush
		}
		// No decisive marker within the buffer window: degrade safely to
		// generic full-line masking.
		if logger != nil {
			logger.Warn("could not detect dump dialect within %d lines: falling back to generic masking", detectMaxBufferedLines)
		}
		p.selectDialect(DialectGeneric)
		return p.flushBuffer(config, cache)
	}
	p.selectDialect(dialect)

	if logger != nil {
		logger.Info("detected dump dialect: %s", p.delegate.Dialect())
	}
	out, drop := p.flushBuffer(config, cache)
	tail, tailDrop := p.delegate.ProcessLine(line, config, cache)
	if tailDrop {
		if out == "" {
			return "", drop
		}
		return out, false
	}
	return out + tail, false
}

// Flush emits any input still buffered at end of stream through the generic
// fallback (detection never became confident).
func (p *detectingDialectParser) Flush(config MaskConfig, cache *Cache) string {
	if p.delegate != nil || len(p.buffer) == 0 {
		return ""
	}
	if logger != nil {
		logger.Warn("could not detect dump dialect: falling back to generic masking")
	}
	p.selectDialect(DialectGeneric)
	out, _ := p.flushBuffer(config, cache)
	return out
}

func (p *detectingDialectParser) selectDialect(dialect DumpDialect) {
	if dialect == DialectGeneric {
		p.delegate = newGenericDialectParser(p.rt)
		return
	}
	p.delegate = NewDialectParser(dialect, p.rt)
}

// flushBuffer replays buffered lines through the chosen delegate.
func (p *detectingDialectParser) flushBuffer(config MaskConfig, cache *Cache) (string, bool) {
	var out strings.Builder
	dropped := true
	for _, buffered := range p.buffer {
		processed, drop := p.delegate.ProcessLine(buffered, config, cache)
		if !drop {
			out.WriteString(processed)
			dropped = false
		}
	}
	p.buffer = nil
	if out.Len() == 0 {
		return "", dropped
	}
	return out.String(), false
}

// detectDialectLine checks one line against the decisive dialect markers.
func detectDialectLine(line string) (DumpDialect, bool) {
	for _, marker := range dialectMarkers {
		if marker.pattern.MatchString(line) {
			return marker.dialect, true
		}
	}
	return "", false
}
