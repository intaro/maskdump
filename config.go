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
	defaultCacheFileName   = ".maskdump_cache.json"
	defaultConfigFileName  = "maskdump.conf"
	defaultEmailRegex      = `\b[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}\b`
	defaultPhoneRegex      = `\b(?:\+7|7|8)(?:\s?\(?\d{3}\)?\s?\d{3}[\s-]?\d{2}[\s-]?\d{2}|\d{10})\b`
	defaultMemoryLimitMB   = 1024 * 4 // 4GB
	defaultCacheFlushCount = 10000
)

type MaskingRule struct {
	Target string `json:"target"`
	Value  string `json:"value"`
}

type MaskingConfig struct {
	Email MaskingRule `json:"email"`
	Phone MaskingRule `json:"phone"`
}

type Config struct {
	CachePath               string        `json:"cache_path"`
	EmailRegex              string        `json:"email_regex"`
	PhoneRegex              string        `json:"phone_regex"`
	EmailWhiteList          string        `json:"email_white_list"`
	PhoneWhiteList          string        `json:"phone_white_list"`
	MemoryLimitMB           int           `json:"memory_limit_mb"`
	CacheFlushCount         int           `json:"cache_flush_count"`
	SkipInsertIntoTableList string        `json:"skip_insert_into_table_list"`
	Masking                 MaskingConfig `json:"masking"`
}

func LoadWhiteList(path string) (map[string]struct{}, error) {
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

func LoadSkipList(path string) (map[string]struct{}, error) {
	skipList := make(map[string]struct{})

	if path == "" {
		return skipList, nil
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
			skipList[line] = struct{}{}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return skipList, nil
}

func LoadConfig(configPath string) error {
	// Set default values first
	defaultConfig := Config{
		CachePath:               filepath.Join(os.Getenv("HOME"), defaultCacheFileName),
		EmailRegex:              defaultEmailRegex,
		PhoneRegex:              defaultPhoneRegex,
		EmailWhiteList:          "",
		PhoneWhiteList:          "",
		MemoryLimitMB:           defaultMemoryLimitMB,
		CacheFlushCount:         defaultCacheFlushCount,
		SkipInsertIntoTableList: "",
		Masking: MaskingConfig{
			Email: MaskingRule{
				Target: "username:2-",
				Value:  "hash:6",
			},
			Phone: MaskingRule{
				Target: "2,3,5,6,8,10",
				Value:  "hash",
			},
		},
	}

	// Apply default values
	AppConfig = defaultConfig

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

	// Create temporary struct to unmarshal JSON
	var fileConfig Config
	if err := json.Unmarshal(data, &fileConfig); err != nil {
		return fmt.Errorf("invalid config file: %v", err)
	}

	// Override default values with non-empty values from config file
	if fileConfig.CachePath != "" {
		AppConfig.CachePath = fileConfig.CachePath
	}
	if fileConfig.EmailRegex != "" {
		AppConfig.EmailRegex = fileConfig.EmailRegex
	}
	if fileConfig.PhoneRegex != "" {
		AppConfig.PhoneRegex = fileConfig.PhoneRegex
	}
	if fileConfig.EmailWhiteList != "" {
		AppConfig.EmailWhiteList = fileConfig.EmailWhiteList
	}
	if fileConfig.PhoneWhiteList != "" {
		AppConfig.PhoneWhiteList = fileConfig.PhoneWhiteList
	}
	if fileConfig.MemoryLimitMB != 0 {
		AppConfig.MemoryLimitMB = fileConfig.MemoryLimitMB
	}
	if fileConfig.CacheFlushCount != 0 {
		AppConfig.CacheFlushCount = fileConfig.CacheFlushCount
	}
	if fileConfig.SkipInsertIntoTableList != "" {
		AppConfig.SkipInsertIntoTableList = fileConfig.SkipInsertIntoTableList
	}
	if fileConfig.Masking.Email.Target != "" {
		AppConfig.Masking.Email.Target = fileConfig.Masking.Email.Target
	}
	if fileConfig.Masking.Email.Value != "" {
		AppConfig.Masking.Email.Value = fileConfig.Masking.Email.Value
	}
	if fileConfig.Masking.Phone.Target != "" {
		AppConfig.Masking.Phone.Target = fileConfig.Masking.Phone.Target
	}
	if fileConfig.Masking.Phone.Value != "" {
		AppConfig.Masking.Phone.Value = fileConfig.Masking.Phone.Value
	}

	// Load white lists
	EmailWhiteList, err = LoadWhiteList(AppConfig.EmailWhiteList)
	if err != nil {
		return fmt.Errorf("failed to load email white list: %v", err)
	}

	PhoneWhiteList, err = LoadWhiteList(AppConfig.PhoneWhiteList)
	if err != nil {
		return fmt.Errorf("failed to load phone white list: %v", err)
	}

	// Compile regular expressions
	EmailRegex, err = regexp.Compile(AppConfig.EmailRegex)
	if err != nil {
		return fmt.Errorf("invalid email regex: %v", err)
	}

	PhoneRegex, err = regexp.Compile(AppConfig.PhoneRegex)
	if err != nil {
		return fmt.Errorf("invalid phone regex: %v", err)
	}

	SkipTableList, err = LoadSkipList(AppConfig.SkipInsertIntoTableList)
	if err != nil {
		return fmt.Errorf("failed to load skip table list: %v", err)
	}

	return nil
}
