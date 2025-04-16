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
)
