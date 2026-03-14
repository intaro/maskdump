package main

import (
	"strings"
	"testing"
)

func BenchmarkParseTargetPositions(b *testing.B) {
	for _, tc := range []struct {
		name   string
		target string
		length int
	}{
		{name: "range", target: "2-24", length: 32},
		{name: "list", target: "1,3,5,7,9,11,13,15", length: 32},
		{name: "tilde", target: "2~4", length: 32},
	} {
		b.Run(tc.name, func(b *testing.B) {
			for b.Loop() {
				_ = parseTargetPositions(tc.target, tc.length)
			}
		})
	}
}

func BenchmarkApplyMasking(b *testing.B) {
	value := "benchmark.user@example.com"
	positions := parseTargetPositions("username:2-", len("benchmark.user"))

	for _, tc := range []struct {
		name      string
		maskValue string
		maskType  TypeMaskingInfo
	}{
		{name: "asterisk", maskValue: "*", maskType: Email},
		{name: "hash", maskValue: "hash:6", maskType: Email},
	} {
		b.Run(tc.name, func(b *testing.B) {
			for b.Loop() {
				_ = applyMasking(value, positions, tc.maskValue, tc.maskType)
			}
		})
	}
}

func BenchmarkMaskEmailWithRules(b *testing.B) {
	setupMaskingDefaultsState()
	b.Cleanup(func() {
		defaultTableParser = NewTableParser(NewRuntimeFromGlobals())
	})
	runtimeState := NewRuntimeFromGlobals()
	email := "benchmark.user@example.com"

	b.ResetTimer()
	for b.Loop() {
		_ = runtimeState.MaskEmailWithRules(email, nil)
	}
}

func BenchmarkMaskPhoneWithRules(b *testing.B) {
	setupMaskingDefaultsState()
	b.Cleanup(func() {
		defaultTableParser = NewTableParser(NewRuntimeFromGlobals())
	})
	runtimeState := NewRuntimeFromGlobals()
	phone := "+7 (900) 111-22-33"

	b.ResetTimer()
	for b.Loop() {
		_ = runtimeState.MaskPhoneWithRules(phone, nil)
	}
}

func BenchmarkProcessLine(b *testing.B) {
	setupMaskingDefaultsState()
	b.Cleanup(func() {
		defaultTableParser = NewTableParser(NewRuntimeFromGlobals())
	})
	line := strings.Repeat(
		"INSERT INTO `users` VALUES (1, 'benchmark.user@example.com', '+7 (900) 111-22-33');\n",
		4,
	)
	config := MaskConfig{emailAlgorithm: "light-hash", phoneAlgorithm: "light-mask"}
	runtimeState := NewRuntimeFromGlobals()
	parser := NewTableParser(runtimeState)

	b.ResetTimer()
	for b.Loop() {
		_ = processLine(line, config, nil, runtimeState, parser, false)
	}
}
