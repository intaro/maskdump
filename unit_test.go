package main

import (
	"regexp"
	"testing"
)

func TestParseTargetPositions(t *testing.T) {
	cases := []struct {
		name     string
		target   string
		length   int
		expected []int
	}{
		{name: "range", target: "2-4", length: 6, expected: []int{1, 2, 3}},
		{name: "tilde_keep_end", target: "~2", length: 5, expected: []int{0, 1, 2}},
		{name: "tilde_keep_start", target: "2~", length: 5, expected: []int{2, 3, 4}},
		{name: "list", target: "1,3,5", length: 5, expected: []int{0, 2, 4}},
	}

	for _, tc := range cases {
		positions := parseTargetPositions(tc.target, tc.length)
		if len(positions) != len(tc.expected) {
			t.Fatalf("%s: expected %v, got %v", tc.name, tc.expected, positions)
		}
		for i := range positions {
			if positions[i] != tc.expected[i] {
				t.Fatalf("%s: expected %v, got %v", tc.name, tc.expected, positions)
			}
		}
	}
}

func TestApplyMaskingAsterisk(t *testing.T) {
	value := "abcdef"
	positions := []int{1, 2, 4}
	masked := applyMasking(value, positions, "*", Email)
	if masked != "a**d*f" {
		t.Fatalf("expected a**d*f, got %s", masked)
	}
}

func TestParseTuple(t *testing.T) {
	tuple := "(1,'a,b','c\\'d',NULL)"
	values := parseTuple(tuple)
	if len(values) != 4 {
		t.Fatalf("expected 4 values, got %d", len(values))
	}
	if values[0] != "1" || values[1] != "'a,b'" || values[2] != "'c'd'" || values[3] != "NULL" {
		t.Fatalf("unexpected values: %v", values)
	}
}

func TestParseTableStructureAndProcessDumpLine(t *testing.T) {
	withTestGlobals(t, func() {
		setupMaskingDefaults(t)

		ProcessingTables = map[string]TableConfig{
			"users": {
				Email: []string{"EMAIL"},
				Phone: []string{"PHONE"},
			},
		}
		runtimeState := newTestRuntime()
		parser := NewTableParser(runtimeState)

		createLines := []string{
			"CREATE TABLE `users` (",
			"  `ID` int,",
			"  `EMAIL` varchar(255),",
			"  `PHONE` varchar(32)",
			");",
		}
		for _, line := range createLines {
			parser.ParseTableStructure(line)
		}

		config := MaskConfig{emailAlgorithm: "light-hash", phoneAlgorithm: "light-mask"}
		line := "INSERT INTO `users` VALUES (1, 'test@example.com', '+7 (123) 456-78-90');"
		out := parser.ProcessDumpLine(line, config, nil)

		if out == line {
			t.Fatalf("expected masked output, got unchanged line")
		}
		if !regexp.MustCompile(`t098f6b@example.com`).MatchString(out) {
			t.Fatalf("expected masked email in output, got: %s", out)
		}
		if !regexp.MustCompile(`\+7 \(\d{3}\) \d{3}-\d{2}-\d{2}`).MatchString(out) {
			t.Fatalf("expected phone in standard format, got: %s", out)
		}
		if regexp.MustCompile(`\+7 \(123\) 456-78-90`).MatchString(out) {
			t.Fatalf("expected phone to be masked, got unchanged: %s", out)
		}
	})
}
