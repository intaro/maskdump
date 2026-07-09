package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeConfigFixture creates a config file plus the auxiliary files it
// references inside a temp dir and returns the config path.
func writeConfigFixture(t *testing.T, body string) string {
	t.Helper()

	dir := t.TempDir()
	for name, content := range map[string]string{
		"skip.txt":   "secrets\n",
		"nomask.txt": "raw_data\n",
	} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
			t.Fatalf("failed to write %s: %v", name, err)
		}
	}

	replacer := strings.NewReplacer(
		"__CACHE__", filepath.Join(dir, "cache.json"),
		"__SKIP__", filepath.Join(dir, "skip.txt"),
		"__NOMASK__", filepath.Join(dir, "nomask.txt"),
	)
	configPath := filepath.Join(dir, "maskdump.conf")
	if err := os.WriteFile(configPath, []byte(replacer.Replace(body)), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	return configPath
}

func TestLoadConfigNewTableKeys(t *testing.T) {
	withTestGlobals(t, func() {
		configPath := writeConfigFixture(t, `{
			"cache_path": "__CACHE__",
			"skip_table_data_list": "__SKIP__",
			"no_masking_table_list": "__NOMASK__",
			"masking_tables": {"users": {"email": ["email"]}}
		}`)

		if err := LoadConfig(configPath); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := SkipTableList["secrets"]; !ok {
			t.Fatalf("expected skip_table_data_list loaded, got: %v", SkipTableList)
		}
		if _, ok := NoMaskTableList["raw_data"]; !ok {
			t.Fatalf("expected no_masking_table_list loaded, got: %v", NoMaskTableList)
		}
		if cfg, ok := ProcessingTables["users"]; !ok || len(cfg.Email) != 1 || cfg.Email[0] != "email" {
			t.Fatalf("expected masking_tables loaded, got: %v", ProcessingTables)
		}
	})
}

func TestLoadConfigDeprecatedTableKeys(t *testing.T) {
	withTestGlobals(t, func() {
		configPath := writeConfigFixture(t, `{
			"cache_path": "__CACHE__",
			"skip_insert_into_table_list": "__SKIP__",
			"processing_tables": {"users": {"email": ["email"]}}
		}`)

		if err := LoadConfig(configPath); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := SkipTableList["secrets"]; !ok {
			t.Fatalf("expected deprecated skip_insert_into_table_list still loaded, got: %v", SkipTableList)
		}
		if _, ok := ProcessingTables["users"]; !ok {
			t.Fatalf("expected deprecated processing_tables still loaded, got: %v", ProcessingTables)
		}
	})
}

func TestLoadConfigRejectsAliasConflicts(t *testing.T) {
	cases := map[string]string{
		"skip lists": `{
			"cache_path": "__CACHE__",
			"skip_insert_into_table_list": "__SKIP__",
			"skip_table_data_list": "__SKIP__"
		}`,
		"masking tables": `{
			"cache_path": "__CACHE__",
			"processing_tables": {"users": {"email": ["email"]}},
			"masking_tables": {"users": {"email": ["email"]}}
		}`,
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			withTestGlobals(t, func() {
				configPath := writeConfigFixture(t, body)
				if err := LoadConfig(configPath); err == nil {
					t.Fatal("expected error when both the new key and its deprecated alias are set")
				}
			})
		})
	}
}
