package main

import (
	"regexp"
)

var (
	// AppConfig holds the runtime configuration.
	AppConfig Config
	// EmailRegex is the compiled regex used to match emails.
	EmailRegex *regexp.Regexp
	// PhoneRegex is the compiled regex used to match phone numbers.
	PhoneRegex *regexp.Regexp
	// EmailWhiteList contains email values that must not be masked.
	EmailWhiteList map[string]struct{}
	// PhoneWhiteList contains phone values that must not be masked.
	PhoneWhiteList map[string]struct{}
	// SkipTableList contains table names to skip during processing.
	SkipTableList map[string]struct{}
	// ProcessingTables defines which tables and fields are masked in selective mode.
	ProcessingTables map[string]TableConfig
	insertRegex      = regexp.MustCompile(`INSERT INTO ` + "`" + `(.+?)` + "`" + ` VALUES (.+)`)
	tupleRegex       = regexp.MustCompile(`\((?:[^()'"\\]|'(?:\\.|[^'\\])*'|"(?:\\.|[^"\\])*"|\\.|\([^()]*\))*\)`)
)

// TypeMaskingInfo is a data type marker for masking algorithms.
type TypeMaskingInfo int

const (
	// Email indicates email masking.
	Email TypeMaskingInfo = iota + 1
	// Phone indicates phone masking.
	Phone
)

// String returns the string representation of the TypeMaskingInfo
func (s TypeMaskingInfo) String() string {
	return [...]string{"Email", "Phone"}[s-1]
}

// Index returns the index of the TypeMaskingInfo
func (s TypeMaskingInfo) Index() int {
	return int(s)
}
