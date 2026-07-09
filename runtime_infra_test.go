package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoggerLevelFiltering(t *testing.T) {
	path := filepath.Join(t.TempDir(), "app.log")
	l, err := NewLogger(LogConfig{Path: path, Level: "warn"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	l.Debug("debug %d", 1)
	l.Info("info %d", 2)
	l.Warn("warn %s", "x")
	l.Error("error %s", "y")

	if err := l.Check(); err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if err := l.Close(); err != nil {
		t.Fatalf("close failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read log: %v", err)
	}
	out := string(data)
	if strings.Contains(out, "DEBUG") || strings.Contains(out, "INFO") {
		t.Fatalf("expected debug/info filtered at warn level, got: %q", out)
	}
	if !strings.Contains(out, "[WARN] warn x") || !strings.Contains(out, "[ERROR] error y") {
		t.Fatalf("expected warn and error entries, got: %q", out)
	}
}

func TestLoggerDebugLevelLogsEverything(t *testing.T) {
	path := filepath.Join(t.TempDir(), "app.log")
	l, err := NewLogger(LogConfig{Path: path, Level: "debug"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	l.Debug("d")
	l.Info("i")
	l.Warn("w")
	l.Error("e")
	if err := l.Close(); err != nil {
		t.Fatalf("close failed: %v", err)
	}

	data, _ := os.ReadFile(path)
	for _, level := range []string{"[DEBUG]", "[INFO]", "[WARN]", "[ERROR]"} {
		if !strings.Contains(string(data), level) {
			t.Fatalf("expected %s entry at debug level, got: %q", level, data)
		}
	}
}

func TestNewLoggerCreatesNestedDirectory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "a", "b", "app.log")
	l, err := NewLogger(LogConfig{Path: path, Level: "error"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := l.Close(); err != nil {
		t.Fatalf("close failed: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected log file created: %v", err)
	}
}

func TestNewLoggerFailsWhenDirIsAFile(t *testing.T) {
	blocker := filepath.Join(t.TempDir(), "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0644); err != nil {
		t.Fatalf("failed to create blocker: %v", err)
	}
	if _, err := NewLogger(LogConfig{Path: filepath.Join(blocker, "app.log")}); err == nil {
		t.Fatal("expected error when the log directory path is a file")
	}
}

func TestCheckOnClosedLogger(t *testing.T) {
	l := &Logger{}
	if err := l.Check(); err == nil {
		t.Fatal("expected error for logger without a file")
	}
}

func TestParseFlags(t *testing.T) {
	origArgs := os.Args
	t.Cleanup(func() {
		os.Args = origArgs
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	})

	os.Args = []string{"maskdump", "--mask-email=light-hash", "--mask-phone=light-mask",
		"--no-cache", "--config=/x.conf", "--db-format=mysql"}
	cfg := parseFlags()

	if cfg.emailAlgorithm != "light-hash" || cfg.phoneAlgorithm != "light-mask" {
		t.Fatalf("unexpected algorithms: %+v", cfg)
	}
	if cfg.cacheEnabled {
		t.Fatal("expected cache disabled by --no-cache")
	}
	if cfg.configFile != "/x.conf" || cfg.dbFormat != "mysql" {
		t.Fatalf("unexpected config/db-format: %+v", cfg)
	}
}

func TestResolveDialect(t *testing.T) {
	withTestGlobals(t, func() {
		AppConfig.DBFormat = "postgresql"

		if d, err := resolveDialect(MaskConfig{dbFormat: "mysql"}); err != nil || d != DialectMySQL {
			t.Fatalf("expected CLI flag to win, got %s, %v", d, err)
		}
		if d, err := resolveDialect(MaskConfig{}); err != nil || d != DialectPostgreSQL {
			t.Fatalf("expected config fallback, got %s, %v", d, err)
		}
		if _, err := resolveDialect(MaskConfig{dbFormat: "dbase"}); err == nil {
			t.Fatal("expected error for unsupported db-format")
		}
	})
}

func TestProcessLineDelegatesToParser(t *testing.T) {
	withTestGlobals(t, func() {
		setupMaskingDefaults(t)

		parser := NewDialectParser(DialectGeneric, newTestRuntime())
		if parser.Dialect() != DialectGeneric {
			t.Fatalf("expected generic dialect, got %s", parser.Dialect())
		}

		out, drop := processLine("mask@me.com\n", bothAlgorithms(), nil, parser)
		if drop || strings.Contains(out, "mask@me.com") {
			t.Fatalf("expected full-line masking via generic parser, got: %q (drop=%v)", out, drop)
		}
	})
}

func TestGetDefaultConfigPaths(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "/xdg")
	paths := getDefaultConfigPaths()
	if paths[1] != filepath.Join("/xdg", defaultConfigDir, defaultConfigName) {
		t.Fatalf("expected XDG_CONFIG_HOME to be used, got: %v", paths)
	}

	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("HOME", "/home/x")
	paths = getDefaultConfigPaths()
	if paths[1] != filepath.Join("/home/x", ".config", defaultConfigDir, defaultConfigName) {
		t.Fatalf("expected HOME fallback, got: %v", paths)
	}
	if paths[0] != "./maskdump.conf" || paths[2] != "/etc/maskdump.conf" {
		t.Fatalf("unexpected search order: %v", paths)
	}
}

func TestGetDefaultLogPath(t *testing.T) {
	if got := getDefaultLogPath("/var/log/custom.log"); got != "/var/log/custom.log" {
		t.Fatalf("expected explicit path returned, got: %s", got)
	}

	t.Setenv("XDG_STATE_HOME", "/state")
	if got := getDefaultLogPath(""); got != filepath.Join("/state", "maskdump", "logs", "maskdump.log") {
		t.Fatalf("expected XDG_STATE_HOME path, got: %s", got)
	}

	t.Setenv("XDG_STATE_HOME", "")
	t.Setenv("HOME", "/home/x")
	if got := getDefaultLogPath(""); got != filepath.Join("/home/x", ".local", "state", "maskdump", "logs", "maskdump.log") {
		t.Fatalf("expected HOME fallback path, got: %s", got)
	}
}

func TestLoadWhiteList(t *testing.T) {
	path := filepath.Join(t.TempDir(), "wl.txt")
	if err := os.WriteFile(path, []byte("a@b.com\n\n  c@d.com  \n"), 0644); err != nil {
		t.Fatalf("failed to write fixture: %v", err)
	}

	wl, err := LoadWhiteList(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(wl) != 2 {
		t.Fatalf("expected 2 entries (blank lines skipped, values trimmed), got: %v", wl)
	}
	if _, ok := wl["c@d.com"]; !ok {
		t.Fatalf("expected trimmed entry present, got: %v", wl)
	}

	if wl, err := LoadWhiteList(""); err != nil || len(wl) != 0 {
		t.Fatalf("expected empty list for empty path, got: %v, %v", wl, err)
	}
	if _, err := LoadWhiteList(filepath.Join(t.TempDir(), "missing.txt")); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestMemoryLimitCheck(t *testing.T) {
	origLimit, origUsage := memoryLimit, currentMemUsage
	t.Cleanup(func() { memoryLimit, currentMemUsage = origLimit, origUsage })

	memoryLimit = 100
	currentMemUsage = 200
	if !checkMemoryLimit() {
		t.Fatal("expected limit exceeded")
	}
	currentMemUsage = 50
	if checkMemoryLimit() {
		t.Fatal("expected limit not exceeded")
	}
}

func TestFreeMemoryFlushesAndClearsCache(t *testing.T) {
	withTestGlobals(t, func() {
		AppConfig.CachePath = filepath.Join(t.TempDir(), "cache.json")

		freeMemory(nil) // must not panic

		cache := &Cache{
			Emails: map[string]string{"a@b.com": "x@b.com"},
			Phones: map[string]string{},
		}
		freeMemory(cache)

		if len(cache.Emails) != 0 {
			t.Fatalf("expected in-memory cache cleared, got: %v", cache.Emails)
		}
		data, err := os.ReadFile(AppConfig.CachePath)
		if err != nil {
			t.Fatalf("expected cache flushed to disk: %v", err)
		}
		if !strings.Contains(string(data), "a@b.com") {
			t.Fatalf("expected flushed cache to keep entries, got: %q", data)
		}
	})
}

func TestTypeMaskingInfo(t *testing.T) {
	if Email.String() != "Email" || Phone.String() != "Phone" {
		t.Fatalf("unexpected names: %s, %s", Email, Phone)
	}
	if Email.Index() != 1 || Phone.Index() != 2 {
		t.Fatalf("unexpected indexes: %d, %d", Email.Index(), Phone.Index())
	}
}

func TestLegacyPackageLevelTableAPI(t *testing.T) {
	withTestGlobals(t, func() {
		setupMaskingDefaults(t)
		ProcessingTables = map[string]TableConfig{
			"users": {Email: []string{"email"}},
		}
		defaultTableParser = NewTableParser(NewRuntimeFromGlobals())

		ParseTableStructure("CREATE TABLE `users` (")
		ParseTableStructure("  `id` int,")
		ParseTableStructure("  `email` varchar(255)")
		ParseTableStructure(");")

		info, ok := GetTableInfo("users")
		if !ok || len(info.Fields) != 2 {
			t.Fatalf("expected users table with 2 fields, got: %+v (ok=%v)", info, ok)
		}
		if all := GetAllTables(); len(all) != 1 {
			t.Fatalf("expected 1 table, got: %v", all)
		}

		out := ProcessDumpLine("INSERT INTO `users` VALUES (1,'test@example.com')", bothAlgorithms(), nil)
		if !strings.Contains(out, "t098f6b@example.com") {
			t.Fatalf("expected email masked via legacy API, got: %q", out)
		}
	})
}

func TestMaskPhoneWhiteListAndCache(t *testing.T) {
	withTestGlobals(t, func() {
		setupMaskingDefaults(t)
		PhoneWhiteList = map[string]struct{}{"+79001234567": {}}
		rt := newTestRuntime()

		if got := rt.MaskPhoneWithRules("+79001234567", nil); got != "+79001234567" {
			t.Fatalf("expected whitelisted phone untouched, got: %s", got)
		}

		cache := &Cache{
			Emails: map[string]string{},
			Phones: map[string]string{"+70000000000": "+71111111111"},
		}
		if got := rt.MaskPhoneWithRules("+70000000000", cache); got != "+71111111111" {
			t.Fatalf("expected cached value returned, got: %s", got)
		}
	})
}

func TestSplitTrailingNewlineCRLF(t *testing.T) {
	body, newline := splitTrailingNewline("line\r\n")
	if body != "line" || newline != "\r\n" {
		t.Fatalf("expected CRLF split, got: %q + %q", body, newline)
	}
	body, newline = splitTrailingNewline("line")
	if body != "line" || newline != "" {
		t.Fatalf("expected no newline, got: %q + %q", body, newline)
	}
}
