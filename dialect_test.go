package main

import (
	"strings"
	"testing"
)

// processDump feeds a whole dump through a parser line by line, imitating
// the main loop, and returns the produced output.
func processDump(t *testing.T, parser DialectParser, config MaskConfig, input string) string {
	t.Helper()

	var out strings.Builder
	for _, line := range strings.SplitAfter(input, "\n") {
		if line == "" {
			continue
		}
		res, drop := parser.ProcessLine(line, config, nil)
		if !drop {
			out.WriteString(res)
		}
	}
	if fp, ok := parser.(flushableParser); ok {
		out.WriteString(fp.Flush(config, nil))
	}
	return out.String()
}

func bothAlgorithms() MaskConfig {
	return MaskConfig{emailAlgorithm: "light-hash", phoneAlgorithm: "light-mask"}
}

func TestParseDumpDialect(t *testing.T) {
	for value, expected := range map[string]DumpDialect{
		"":           DialectAuto,
		"auto":       DialectAuto,
		"MySQL":      DialectMySQL,
		"postgresql": DialectPostgreSQL,
		"oracle":     DialectOracle,
		"mssql":      DialectMSSQL,
		"sqlite":     DialectSQLite,
		"firebird":   DialectFirebird,
	} {
		dialect, err := ParseDumpDialect(value)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", value, err)
		}
		if dialect != expected {
			t.Fatalf("expected %s for %q, got %s", expected, value, dialect)
		}
	}

	if _, err := ParseDumpDialect("dbase"); err == nil {
		t.Fatal("expected error for unsupported dialect")
	}
}

func TestNormalizeTableName(t *testing.T) {
	cases := []struct {
		raw   string
		full  string
		plain string
	}{
		{"users", "users", "users"},
		{"`users`", "users", "users"},
		{`"public"."tst_users"`, "public.tst_users", "tst_users"},
		{"[dbo].[tst_users]", "dbo.tst_users", "tst_users"},
		{`public.tst_users`, "public.tst_users", "tst_users"},
	}
	for _, tc := range cases {
		full, plain := normalizeTableName(tc.raw)
		if full != tc.full || plain != tc.plain {
			t.Fatalf("%s: expected (%s, %s), got (%s, %s)", tc.raw, tc.full, tc.plain, full, plain)
		}
	}
}

func TestMySQLDialectSkipTable(t *testing.T) {
	withTestGlobals(t, func() {
		setupMaskingDefaults(t)
		SkipTableList = map[string]struct{}{"secrets": {}}

		parser := NewDialectParser(DialectMySQL, newTestRuntime())
		out := processDump(t, parser, bothAlgorithms(),
			"INSERT INTO `secrets` VALUES (1,'x@y.com');\n"+
				"INSERT INTO `visible` VALUES (2,'keep');\n")

		if strings.Contains(out, "secrets") {
			t.Fatalf("expected skipped table to be dropped, got: %s", out)
		}
		if !strings.Contains(out, "visible") {
			t.Fatalf("expected other tables to remain, got: %s", out)
		}
	})
}

// A single pass must mask configured fields exactly once even when both
// masking algorithms are enabled (regression: the line was processed twice).
func TestMySQLDialectSelectiveMasksOnce(t *testing.T) {
	withTestGlobals(t, func() {
		setupMaskingDefaults(t)
		ProcessingTables = map[string]TableConfig{
			"users": {Email: []string{"email"}, Phone: []string{"phone"}},
		}

		dump := "CREATE TABLE `users` (\n" +
			"  `id` int,\n" +
			"  `email` varchar(255),\n" +
			"  `phone` varchar(32)\n" +
			");\n" +
			"INSERT INTO `users` VALUES (1,'test@example.com','+7 (123) 456-78-90');\n"

		parser := NewDialectParser(DialectMySQL, newTestRuntime())
		out := processDump(t, parser, bothAlgorithms(), dump)

		// md5("test")[:6] == 098f6b: a single masking application keeps
		// the first character and replaces the tail once.
		if !strings.Contains(out, "t098f6b@example.com") {
			t.Fatalf("expected email masked exactly once, got: %s", out)
		}
		if !strings.HasSuffix(out, "\n") {
			t.Fatalf("expected trailing newline to be preserved, got: %q", out)
		}
	})
}

func TestPostgresCopyMasking(t *testing.T) {
	withTestGlobals(t, func() {
		setupMaskingDefaults(t)
		ProcessingTables = map[string]TableConfig{
			"public.tst_users": {Email: []string{"email"}, Phone: []string{"phone"}},
		}

		dump := `COPY public.tst_users (id, login, email, phone, extra) FROM stdin;` + "\n" +
			"1\tignore@stay.com\ttest@example.com\t+7 (123) 456-78-90\t{\"note\": \"json stays\"}\n" +
			"2\tlogin2\t\\N\t\\N\tvalue\twith\\ttab\n" +
			`\.` + "\n" +
			"SELECT 1;\n"

		parser := NewDialectParser(DialectPostgreSQL, newTestRuntime())
		out := processDump(t, parser, bothAlgorithms(), dump)

		if !strings.Contains(out, "t098f6b@example.com") {
			t.Fatalf("expected configured email column masked, got: %s", out)
		}
		if !strings.Contains(out, "ignore@stay.com") {
			t.Fatalf("expected unconfigured column untouched, got: %s", out)
		}
		if strings.Contains(out, "+7 (123) 456-78-90") {
			t.Fatalf("expected configured phone column masked, got: %s", out)
		}
		if !strings.Contains(out, `\N`) {
			t.Fatalf("expected NULL markers preserved, got: %s", out)
		}
		if !strings.Contains(out, `{"note": "json stays"}`) {
			t.Fatalf("expected JSON value untouched, got: %s", out)
		}
		if !strings.Contains(out, `\.`+"\n") {
			t.Fatalf("expected COPY terminator preserved, got: %s", out)
		}
		if !strings.Contains(out, "SELECT 1;") {
			t.Fatalf("expected trailing statement preserved, got: %s", out)
		}
	})
}

func TestPostgresCopyMatchesPlainConfigKey(t *testing.T) {
	withTestGlobals(t, func() {
		setupMaskingDefaults(t)
		ProcessingTables = map[string]TableConfig{
			"tst_users": {Email: []string{"email"}},
		}

		dump := `COPY "public"."tst_users" ("id", "email") FROM stdin;` + "\n" +
			"1\ttest@example.com\n" +
			`\.` + "\n"

		parser := NewDialectParser(DialectPostgreSQL, newTestRuntime())
		out := processDump(t, parser, bothAlgorithms(), dump)

		if !strings.Contains(out, "t098f6b@example.com") {
			t.Fatalf("expected plain config key to match schema-qualified COPY, got: %s", out)
		}
	})
}

func TestPostgresCopySkipTable(t *testing.T) {
	withTestGlobals(t, func() {
		setupMaskingDefaults(t)
		SkipTableList = map[string]struct{}{"tst_secrets": {}}

		dump := "COPY public.tst_secrets (id, token) FROM stdin;\n" +
			"1\tsecret-token\n" +
			`\.` + "\n" +
			"COPY public.tst_keep (id) FROM stdin;\n" +
			"7\n" +
			`\.` + "\n"

		parser := NewDialectParser(DialectPostgreSQL, newTestRuntime())
		out := processDump(t, parser, bothAlgorithms(), dump)

		if strings.Contains(out, "secret-token") || strings.Contains(out, "tst_secrets") {
			t.Fatalf("expected skipped COPY block to be dropped entirely, got: %s", out)
		}
		if !strings.Contains(out, "tst_keep") || !strings.Contains(out, "7\n") {
			t.Fatalf("expected following COPY block to remain, got: %s", out)
		}
	})
}

func TestPostgresMultiLineInsertMasking(t *testing.T) {
	withTestGlobals(t, func() {
		setupMaskingDefaults(t)
		ProcessingTables = map[string]TableConfig{
			"tst_users": {Email: []string{"email"}},
		}

		dump := "INSERT INTO public.tst_users (id, login, email) VALUES\n" +
			"(1, 'login-one', 'test@example.com'),\n" +
			"(2, 'keep@untouched.com', 'test@example.com');\n" +
			"INSERT INTO public.tst_other (id, email) VALUES\n" +
			"(3, 'other@stays.com');\n"

		parser := NewDialectParser(DialectPostgreSQL, newTestRuntime())
		out := processDump(t, parser, bothAlgorithms(), dump)

		if strings.Count(out, "t098f6b@example.com") != 2 {
			t.Fatalf("expected email column masked in every row, got: %s", out)
		}
		if !strings.Contains(out, "keep@untouched.com") {
			t.Fatalf("expected unconfigured column untouched, got: %s", out)
		}
		if !strings.Contains(out, "other@stays.com") {
			t.Fatalf("expected unconfigured table untouched, got: %s", out)
		}
	})
}

func TestPostgresMultiLineInsertSkip(t *testing.T) {
	withTestGlobals(t, func() {
		setupMaskingDefaults(t)
		SkipTableList = map[string]struct{}{"tst_secrets": {}}

		dump := "INSERT INTO public.tst_secrets (id, token) VALUES\n" +
			"(1, 'secret-a'),\n" +
			"(2, 'secret-b');\n" +
			"SELECT 'after';\n"

		parser := NewDialectParser(DialectPostgreSQL, newTestRuntime())
		out := processDump(t, parser, bothAlgorithms(), dump)

		if strings.Contains(out, "secret") {
			t.Fatalf("expected skipped INSERT statement dropped, got: %s", out)
		}
		if !strings.Contains(out, "SELECT 'after';") {
			t.Fatalf("expected the following line to remain, got: %s", out)
		}
	})
}

func TestOracleInsertMasking(t *testing.T) {
	withTestGlobals(t, func() {
		setupMaskingDefaults(t)
		ProcessingTables = map[string]TableConfig{
			"tst_users": {Email: []string{"email"}, Phone: []string{"phone"}},
		}

		dump := "INSERT INTO tst_users (id, login, email, phone) VALUES (1, 'test@example.com', 'test@example.com', '+7 (123) 456-78-90');\n"

		parser := NewDialectParser(DialectOracle, newTestRuntime())
		out := processDump(t, parser, bothAlgorithms(), dump)

		if strings.Count(out, "t098f6b@example.com") != 1 {
			t.Fatalf("expected only the configured email column masked, got: %s", out)
		}
		if !strings.Contains(out, "'test@example.com'") {
			t.Fatalf("expected login column untouched, got: %s", out)
		}
		if strings.Contains(out, "+7 (123) 456-78-90") {
			t.Fatalf("expected phone column masked, got: %s", out)
		}
	})
}

func TestMSSQLInsertMaskingAndSkip(t *testing.T) {
	withTestGlobals(t, func() {
		setupMaskingDefaults(t)
		ProcessingTables = map[string]TableConfig{
			"tst_users": {Email: []string{"email"}},
		}
		SkipTableList = map[string]struct{}{"tst_secrets": {}}

		dump := "SET ANSI_NULLS ON\n" +
			"GO\n" +
			"INSERT INTO [dbo].[tst_users] ([id], [email]) VALUES (1, N'test@example.com')\n" +
			"INSERT INTO [dbo].[tst_secrets] ([id], [token]) VALUES (1, N'secret')\n" +
			"GO\n"

		parser := NewDialectParser(DialectMSSQL, newTestRuntime())
		out := processDump(t, parser, bothAlgorithms(), dump)

		if !strings.Contains(out, "N't098f6b@example.com'") {
			t.Fatalf("expected email masked inside N'...' literal, got: %s", out)
		}
		if strings.Contains(out, "secret") {
			t.Fatalf("expected skipped table INSERT dropped, got: %s", out)
		}
		if strings.Count(out, "GO\n") != 2 {
			t.Fatalf("expected GO separators preserved, got: %s", out)
		}
	})
}

func TestSQLiteInsertWithoutColumnList(t *testing.T) {
	withTestGlobals(t, func() {
		setupMaskingDefaults(t)
		ProcessingTables = map[string]TableConfig{
			"users": {Email: []string{"email"}},
		}

		dump := "PRAGMA foreign_keys=OFF;\n" +
			"CREATE TABLE \"users\" (id INTEGER PRIMARY KEY, email TEXT, note TEXT);\n" +
			"INSERT INTO \"users\" VALUES(1,'test@example.com','note keep@here.com');\n"

		parser := NewDialectParser(DialectSQLite, newTestRuntime())
		out := processDump(t, parser, bothAlgorithms(), dump)

		if !strings.Contains(out, "t098f6b@example.com") {
			t.Fatalf("expected email masked via CREATE TABLE positions, got: %s", out)
		}
		if !strings.Contains(out, "keep@here.com") {
			t.Fatalf("expected unconfigured column untouched, got: %s", out)
		}
	})
}

func TestFirebirdInsertMasking(t *testing.T) {
	withTestGlobals(t, func() {
		setupMaskingDefaults(t)
		ProcessingTables = map[string]TableConfig{
			"USERS": {Email: []string{"EMAIL"}},
		}

		dump := "SET TERM ^ ;\n" +
			"INSERT INTO \"USERS\" (\"ID\", \"EMAIL\") VALUES (1, 'test@example.com');\n"

		parser := NewDialectParser(DialectFirebird, newTestRuntime())
		out := processDump(t, parser, bothAlgorithms(), dump)

		if !strings.Contains(out, "t098f6b@example.com") {
			t.Fatalf("expected email masked, got: %s", out)
		}
		if !strings.Contains(out, "SET TERM ^ ;") {
			t.Fatalf("expected SET TERM line preserved, got: %s", out)
		}
	})
}

func TestSelectiveModeWithoutColumnInfoLeavesLineUntouched(t *testing.T) {
	withTestGlobals(t, func() {
		setupMaskingDefaults(t)
		ProcessingTables = map[string]TableConfig{
			"users": {Email: []string{"email"}},
		}

		// No CREATE TABLE seen and no column list in the INSERT: the safe
		// behavior is to pass the line through unmasked.
		dump := "INSERT INTO users VALUES (1,'test@example.com');\n"

		parser := NewDialectParser(DialectOracle, newTestRuntime())
		out := processDump(t, parser, bothAlgorithms(), dump)

		if out != dump {
			t.Fatalf("expected line unchanged without column info, got: %s", out)
		}
	})
}

func TestDialectAutoDetection(t *testing.T) {
	cases := []struct {
		name      string
		firstLine string
		expected  DumpDialect
	}{
		{"mysql", "-- MySQL dump 10.13  Distrib 8.0.36\n", DialectMySQL},
		{"postgresql", "-- PostgreSQL database dump\n", DialectPostgreSQL},
		{"oracle", "-- Oracle Database dump\n", DialectOracle},
		{"mssql", "USE [maskdump]\n", DialectMSSQL},
		{"sqlite", "PRAGMA foreign_keys=OFF;\n", DialectSQLite},
		{"firebird", "SET TERM ^ ;\n", DialectFirebird},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			withTestGlobals(t, func() {
				setupMaskingDefaults(t)
				parser := NewDialectParser(DialectAuto, newTestRuntime())
				out := processDump(t, parser, bothAlgorithms(), tc.firstLine)

				if parser.Dialect() != tc.expected {
					t.Fatalf("expected %s, got %s", tc.expected, parser.Dialect())
				}
				if out != tc.firstLine {
					t.Fatalf("expected marker line passed through, got: %q", out)
				}
			})
		})
	}
}

func TestDialectAutoDetectionBuffersAndReplays(t *testing.T) {
	withTestGlobals(t, func() {
		setupMaskingDefaults(t)
		SkipTableList = map[string]struct{}{"secrets": {}}

		dump := "-- comment first\n" +
			"-- PostgreSQL database dump\n" +
			"COPY secrets (id) FROM stdin;\n" +
			"1\n" +
			`\.` + "\n"

		parser := NewDialectParser(DialectAuto, newTestRuntime())
		out := processDump(t, parser, bothAlgorithms(), dump)

		if !strings.Contains(out, "-- comment first") {
			t.Fatalf("expected buffered pre-marker line replayed, got: %q", out)
		}
		if strings.Contains(out, "secrets") {
			t.Fatalf("expected skip list applied after detection, got: %q", out)
		}
	})
}

func TestDialectAutoDetectionFallbackFlush(t *testing.T) {
	withTestGlobals(t, func() {
		setupMaskingDefaults(t)

		dump := "no markers here test@example.com\n" +
			"just plain text\n"

		parser := NewDialectParser(DialectAuto, newTestRuntime())
		out := processDump(t, parser, bothAlgorithms(), dump)

		if !strings.Contains(out, "t098f6b@example.com") {
			t.Fatalf("expected generic masking applied on flush, got: %q", out)
		}
		if !strings.Contains(out, "just plain text") {
			t.Fatalf("expected all buffered lines flushed, got: %q", out)
		}
	})
}

func TestBlankLinesArePreserved(t *testing.T) {
	withTestGlobals(t, func() {
		setupMaskingDefaults(t)

		parser := NewDialectParser(DialectMySQL, newTestRuntime())
		out := processDump(t, parser, bothAlgorithms(), "line one\n\nline two\n")

		if out != "line one\n\nline two\n" {
			t.Fatalf("expected blank line preserved, got: %q", out)
		}
	})
}

// Oracle folds unquoted identifiers to upper case, so lower-case config keys
// must match upper-case table and column names from the dump.
func TestOracleCaseInsensitiveIdentifiers(t *testing.T) {
	withTestGlobals(t, func() {
		setupMaskingDefaults(t)
		ProcessingTables = map[string]TableConfig{
			"customers": {Email: []string{"email"}},
		}
		SkipTableList = map[string]struct{}{"audit_log": {}}

		dump := "INSERT INTO CUSTOMERS (ID, EMAIL) VALUES (1, 'test@example.com');\n" +
			"INSERT INTO AUDIT_LOG (ID, EMAIL) VALUES (7, 'gone@example.com');\n"

		parser := NewDialectParser(DialectOracle, newTestRuntime())
		out := processDump(t, parser, bothAlgorithms(), dump)

		if !strings.Contains(out, "t098f6b@example.com") {
			t.Fatalf("expected upper-case EMAIL column masked via lower-case config, got: %q", out)
		}
		if strings.Contains(out, "AUDIT_LOG") {
			t.Fatalf("expected upper-case AUDIT_LOG dropped via lower-case skip entry, got: %q", out)
		}
	})
}

// Oracle CREATE TABLE column info must also match case-insensitively when
// the INSERT carries no column list.
func TestOracleCaseInsensitiveCreateTableColumns(t *testing.T) {
	withTestGlobals(t, func() {
		setupMaskingDefaults(t)
		ProcessingTables = map[string]TableConfig{
			"customers": {Email: []string{"email"}},
		}

		dump := "CREATE TABLE CUSTOMERS (ID NUMBER, EMAIL VARCHAR2(255));\n" +
			"INSERT INTO CUSTOMERS VALUES (1, 'test@example.com');\n"

		parser := NewDialectParser(DialectOracle, newTestRuntime())
		out := processDump(t, parser, bothAlgorithms(), dump)

		if !strings.Contains(out, "t098f6b@example.com") {
			t.Fatalf("expected column positions from CREATE TABLE matched case-insensitively, got: %q", out)
		}
	})
}

// Identifier matching stays case-sensitive for non-Oracle dialects.
func TestMySQLIdentifiersStayCaseSensitive(t *testing.T) {
	withTestGlobals(t, func() {
		setupMaskingDefaults(t)
		ProcessingTables = map[string]TableConfig{
			"users": {Email: []string{"email"}},
		}

		dump := "CREATE TABLE `Users` (\n" +
			"  `id` int,\n" +
			"  `email` varchar(255)\n" +
			");\n" +
			"INSERT INTO `Users` VALUES (1,'test@example.com');\n"

		parser := NewDialectParser(DialectMySQL, newTestRuntime())
		out := processDump(t, parser, bothAlgorithms(), dump)

		if !strings.Contains(out, "test@example.com") {
			t.Fatalf("expected case-mismatched table to stay unmasked, got: %q", out)
		}
	})
}

func TestNoMaskTableListMySQL(t *testing.T) {
	withTestGlobals(t, func() {
		setupMaskingDefaults(t)
		NoMaskTableList = map[string]struct{}{"raw_data": {}}

		parser := NewDialectParser(DialectMySQL, newTestRuntime())
		out := processDump(t, parser, bothAlgorithms(),
			"INSERT INTO `raw_data` VALUES (1,'keep@asis.com');\n"+
				"INSERT INTO `other` VALUES (2,'mask@me.com');\n")

		if !strings.Contains(out, "keep@asis.com") {
			t.Fatalf("expected no-mask table left untouched, got: %q", out)
		}
		if strings.Contains(out, "mask@me.com") {
			t.Fatalf("expected other tables full-line masked, got: %q", out)
		}
	})
}

func TestNoMaskTableListPostgresCopy(t *testing.T) {
	withTestGlobals(t, func() {
		setupMaskingDefaults(t)
		NoMaskTableList = map[string]struct{}{"raw_data": {}}

		dump := "COPY public.raw_data (id, email) FROM stdin;\n" +
			"1\tkeep@asis.com\n" +
			"\\.\n" +
			"COPY public.other (id, email) FROM stdin;\n" +
			"2\tmask@me.com\n" +
			"\\.\n"

		parser := NewDialectParser(DialectPostgreSQL, newTestRuntime())
		out := processDump(t, parser, bothAlgorithms(), dump)

		if !strings.Contains(out, "keep@asis.com") {
			t.Fatalf("expected no-mask COPY block left untouched, got: %q", out)
		}
		if strings.Contains(out, "mask@me.com") {
			t.Fatalf("expected other COPY rows full-line masked, got: %q", out)
		}
	})
}

// A multi-line VALUES list of a no-mask table passes through untouched, and
// case-insensitive matching applies for Oracle.
func TestNoMaskTableListOracleMultiLineInsert(t *testing.T) {
	withTestGlobals(t, func() {
		setupMaskingDefaults(t)
		NoMaskTableList = map[string]struct{}{"raw_data": {}}

		dump := "INSERT INTO RAW_DATA (ID, EMAIL) VALUES\n" +
			"(1, 'keep@asis.com'),\n" +
			"(2, 'also@keep.com');\n" +
			"INSERT INTO OTHER (ID, EMAIL) VALUES (3, 'mask@me.com');\n"

		parser := NewDialectParser(DialectOracle, newTestRuntime())
		out := processDump(t, parser, bothAlgorithms(), dump)

		if !strings.Contains(out, "keep@asis.com") || !strings.Contains(out, "also@keep.com") {
			t.Fatalf("expected multi-line no-mask insert left untouched, got: %q", out)
		}
		if strings.Contains(out, "mask@me.com") {
			t.Fatalf("expected other tables full-line masked, got: %q", out)
		}
	})
}

// Regression: with only a skip list configured (no selective masking), data
// of non-skipped tables must still get full-line masking. Before the fix
// INSERT statements and COPY rows bypassed the regex path entirely.
func TestFullLineMaskingAppliesWithSkipListOnly(t *testing.T) {
	withTestGlobals(t, func() {
		setupMaskingDefaults(t)
		SkipTableList = map[string]struct{}{"secrets": {}}

		dump := "INSERT INTO [dbo].[secrets] (id, email) VALUES (1, 'gone@example.com')\n" +
			"INSERT INTO [dbo].[other] (id, email) VALUES (2, 'mask@me.com')\n"

		parser := NewDialectParser(DialectMSSQL, newTestRuntime())
		out := processDump(t, parser, bothAlgorithms(), dump)

		if strings.Contains(out, "gone@example.com") {
			t.Fatalf("expected skipped table dropped, got: %q", out)
		}
		if strings.Contains(out, "mask@me.com") {
			t.Fatalf("expected non-skipped INSERT full-line masked, got: %q", out)
		}
	})
}

func TestFullLineMaskingAppliesWithSkipListOnlyPostgresCopy(t *testing.T) {
	withTestGlobals(t, func() {
		setupMaskingDefaults(t)
		SkipTableList = map[string]struct{}{"secrets": {}}

		dump := "COPY public.secrets (id, email) FROM stdin;\n" +
			"1\tgone@example.com\n" +
			"\\.\n" +
			"COPY public.other (id, email) FROM stdin;\n" +
			"2\tmask@me.com\n" +
			"\\.\n"

		parser := NewDialectParser(DialectPostgreSQL, newTestRuntime())
		out := processDump(t, parser, bothAlgorithms(), dump)

		if strings.Contains(out, "gone@example.com") {
			t.Fatalf("expected skipped COPY block dropped, got: %q", out)
		}
		if strings.Contains(out, "mask@me.com") {
			t.Fatalf("expected non-skipped COPY rows full-line masked, got: %q", out)
		}
	})
}
