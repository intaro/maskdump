package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	defaultCacheFileName  = ".maskdump_cache.json"
	defaultConfigFileName = "maskdump.conf"
)

type Config struct {
	CachePath       string `json:"cache_path"`
	EmailRegex      string `json:"email_regex"`
	PhoneRegex      string `json:"phone_regex"`
	EmailWhiteList  string `json:"email_white_list"`
	PhoneWhiteList  string `json:"phone_white_list"`
	MemoryLimitMB   int    `json:"memory_limit_mb"`
	CacheFlushCount int    `json:"cache_flush_count"`
}

var (
	appConfig      Config
	emailRegex     *regexp.Regexp
	phoneRegex     *regexp.Regexp
	emailWhiteList map[string]struct{}
	phoneWhiteList map[string]struct{}
)

func loadWhiteList(path string) (map[string]struct{}, error) {
	whiteList := make(map[string]struct{})

	if path == "" {
		return whiteList, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			whiteList[line] = struct{}{}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return whiteList, nil
}

func loadConfig(configPath string) error {
	// Set default values
	appConfig = Config{
		CachePath:       filepath.Join(os.Getenv("HOME"), defaultCacheFileName),
		EmailRegex:      defaultEmailRegex,
		PhoneRegex:      defaultPhoneRegex,
		EmailWhiteList:  "",
		PhoneWhiteList:  "",
		MemoryLimitMB:   defaultMemoryLimitMB,
		CacheFlushCount: defaultCacheFlushCount,
	}

	if configPath == "" {
		// Try to find config near the binary
		exePath, err := os.Executable()
		if err != nil {
			return nil // Continue with default settings
		}
		configPath = filepath.Join(filepath.Dir(exePath), defaultConfigFileName)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil // File not found - use default settings
	}

	if err := json.Unmarshal(data, &appConfig); err != nil {
		return fmt.Errorf("invalid config file: %v", err)
	}

	// Load white lists
	emailWhiteList, err = loadWhiteList(appConfig.EmailWhiteList)
	if err != nil {
		return fmt.Errorf("failed to load email white list: %v", err)
	}

	phoneWhiteList, err = loadWhiteList(appConfig.PhoneWhiteList)
	if err != nil {
		return fmt.Errorf("failed to load phone white list: %v", err)
	}

	// Compile regular expressions
	emailRegex, err = regexp.Compile(appConfig.EmailRegex)
	if err != nil {
		return fmt.Errorf("invalid email regex: %v", err)
	}

	phoneRegex, err = regexp.Compile(appConfig.PhoneRegex)
	if err != nil {
		return fmt.Errorf("invalid phone regex: %v", err)
	}

	return nil
}
