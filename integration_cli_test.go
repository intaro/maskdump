package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
)

type integrationFixture struct {
	name     string
	input    string
	expected string
}

func TestCLIIntegrationFixtures(t *testing.T) {
	fixtures := collectIntegrationFixtures(t)
	binaryPath := buildMaskdumpBinary(t)

	for _, fixture := range fixtures {
		fixture := fixture
		t.Run(fixture.name, func(t *testing.T) {
			runtimeDir := t.TempDir()
			configPath := filepath.Join(runtimeDir, "integration.conf")
			writeIntegrationConfig(t, configPath, runtimeDir)

			input, err := os.ReadFile(fixture.input)
			if err != nil {
				t.Fatalf("failed to read fixture %s: %v", fixture.input, err)
			}

			expected, err := os.ReadFile(fixture.expected)
			if err != nil {
				t.Fatalf("failed to read expected fixture %s: %v", fixture.expected, err)
			}

			cmd := exec.Command(
				binaryPath,
				"--config", configPath,
				"--mask-email=light-hash",
				"--mask-phone=light-mask",
				"--no-cache",
			)
			cmd.Dir = repoRoot(t)
			cmd.Env = append(
				os.Environ(),
				"HOME="+runtimeDir,
				"XDG_STATE_HOME="+filepath.Join(runtimeDir, "state"),
				"XDG_CONFIG_HOME="+filepath.Join(runtimeDir, "config"),
			)
			cmd.Stdin = bytes.NewReader(input)

			var stdout bytes.Buffer
			var stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			if err := cmd.Run(); err != nil {
				t.Fatalf("maskdump failed: %v\nstderr:\n%s", err, stderr.String())
			}

			if diff := compareFixtureBytes(expected, stdout.Bytes()); diff != "" {
				t.Fatalf("unexpected CLI output for %s:\n%s", fixture.name, diff)
			}
		})
	}
}

func buildMaskdumpBinary(t *testing.T) string {
	t.Helper()

	binaryPath := filepath.Join(t.TempDir(), "maskdump-test-bin")
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	cmd.Dir = repoRoot(t)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to build maskdump binary: %v\n%s", err, string(output))
	}

	return binaryPath
}

func collectIntegrationFixtures(t *testing.T) []integrationFixture {
	t.Helper()

	root := repoRoot(t)
	var fixtures []integrationFixture

	for _, relativeDir := range []string{
		filepath.Join("testdata", "dump"),
		filepath.Join("testdata", "csv"),
	} {
		baseDir := filepath.Join(root, relativeDir)
		err := filepath.WalkDir(baseDir, func(path string, d os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if d.IsDir() {
				return nil
			}

			relPath, err := filepath.Rel(filepath.Join(root, "testdata"), path)
			if err != nil {
				return err
			}

			expectedPath := filepath.Join(root, "testdata", "expected", relPath)
			if _, err := os.Stat(expectedPath); err != nil {
				return err
			}

			fixtures = append(fixtures, integrationFixture{
				name:     filepath.ToSlash(relPath),
				input:    path,
				expected: expectedPath,
			})
			return nil
		})
		if err != nil {
			t.Fatalf("failed to collect fixtures from %s: %v", baseDir, err)
		}
	}

	sort.Slice(fixtures, func(i, j int) bool {
		return fixtures[i].name < fixtures[j].name
	})

	if len(fixtures) == 0 {
		t.Fatal("no integration fixtures found")
	}

	return fixtures
}

func writeIntegrationConfig(t *testing.T, configPath, runtimeDir string) {
	t.Helper()

	templatePath := filepath.Join(repoRoot(t), "testdata", "config", "integration.conf.tmpl")
	content, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("failed to read config template: %v", err)
	}

	rendered := strings.ReplaceAll(string(content), "__CACHE_PATH__", filepath.ToSlash(filepath.Join(runtimeDir, "cache", "cache.json")))
	rendered = strings.ReplaceAll(rendered, "__LOG_PATH__", filepath.ToSlash(filepath.Join(runtimeDir, "logs", "maskdump.log")))

	if err := os.WriteFile(configPath, []byte(rendered), 0644); err != nil {
		t.Fatalf("failed to write config %s: %v", configPath, err)
	}
}

func compareFixtureBytes(expected, actual []byte) string {
	if bytes.Equal(expected, actual) {
		return ""
	}

	expectedLines := strings.Split(string(expected), "\n")
	actualLines := strings.Split(string(actual), "\n")
	limit := len(expectedLines)
	if len(actualLines) < limit {
		limit = len(actualLines)
	}

	for i := 0; i < limit; i++ {
		if expectedLines[i] != actualLines[i] {
			return "line " + strconv.Itoa(i+1) + "\nexpected: " + expectedLines[i] + "\nactual:   " + actualLines[i]
		}
	}

	return "different output size"
}

func repoRoot(t *testing.T) string {
	t.Helper()

	root, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	return root
}
