package main

import (
	"regexp"
	"testing"
)

func withTestGlobals(t *testing.T, fn func()) {
	t.Helper()

	// Save current globals
	origAppConfig := AppConfig
	origEmailRegex := EmailRegex
	origPhoneRegex := PhoneRegex
	origEmailWhiteList := EmailWhiteList
	origPhoneWhiteList := PhoneWhiteList
	origSkipTableList := SkipTableList
	origProcessingTables := ProcessingTables
	origTableInfos := tableInfos
	origCurrentTable := currentTable
	origProcessingTable := processingTable

	t.Cleanup(func() {
		AppConfig = origAppConfig
		EmailRegex = origEmailRegex
		PhoneRegex = origPhoneRegex
		EmailWhiteList = origEmailWhiteList
		PhoneWhiteList = origPhoneWhiteList
		SkipTableList = origSkipTableList
		ProcessingTables = origProcessingTables
		tableInfos = origTableInfos
		currentTable = origCurrentTable
		processingTable = origProcessingTable
	})

	fn()
}

func setupMaskingDefaults(t *testing.T) {
	t.Helper()

	EmailRegex = regexp.MustCompile(defaultEmailRegex)
	PhoneRegex = regexp.MustCompile(defaultPhoneRegex)
	EmailWhiteList = map[string]struct{}{}
	PhoneWhiteList = map[string]struct{}{}
	SkipTableList = map[string]struct{}{}

	AppConfig.Masking = MaskingConfig{
		Email: MaskingRule{Target: "username:2-", Value: "hash:6"},
		Phone: MaskingRule{Target: "2,3,5,6,8,10", Value: "hash"},
	}
}

func countDigits(s string) int {
	count := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			count++
		}
	}
	return count
}

func stripDigits(s string) string {
	out := make([]rune, 0, len(s))
	for _, c := range s {
		if c < '0' || c > '9' {
			out = append(out, c)
		}
	}
	return string(out)
}

func stripDigitsAndStars(s string) string {
	out := make([]rune, 0, len(s))
	for _, c := range s {
		if (c < '0' || c > '9') && c != '*' {
			out = append(out, c)
		}
	}
	return string(out)
}
