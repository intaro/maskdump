package main

import (
	"regexp"
)

var (
	AppConfig      Config
	EmailRegex     *regexp.Regexp
	PhoneRegex     *regexp.Regexp
	EmailWhiteList map[string]struct{}
	PhoneWhiteList map[string]struct{}
	SkipTableList  map[string]struct{}
)

type TypeMaskingInfo int

const (
	Email TypeMaskingInfo = iota + 1
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
