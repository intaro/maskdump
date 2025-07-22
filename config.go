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
	defaultConfigName      = "config"
	defaultConfigDir       = "maskdump"
	defaultCacheFileName   = ".maskdump_cache.json"
	defaultConfigFileName  = "maskdump.conf"
	defaultEmailRegex      = `\b[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}\b`
	defaultPhoneRegex      = `\b(?:\+7|7|8)(?:[\s-]?\(?\d{3}\)?[\s-]?\d{3}[\s-]?\d{2}[\s-]?\d{2}|\d{10})\b`
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

// Структура для хранения конфигурации таблиц
type TableConfig struct {
	Email []string `json:"email"`
	Phone []string `json:"phone"`
}

type Config struct {
	CachePath               string                 `json:"cache_path"`
	EmailRegex              string                 `json:"email_regex"`
	PhoneRegex              string                 `json:"phone_regex"`
	EmailWhiteList          string                 `json:"email_white_list"`
	PhoneWhiteList          string                 `json:"phone_white_list"`
	MemoryLimitMB           int                    `json:"memory_limit_mb"`
	CacheFlushCount         int                    `json:"cache_flush_count"`
	SkipInsertIntoTableList string                 `json:"skip_insert_into_table_list"`
	Masking                 MaskingConfig          `json:"masking"`
	ProcessingTables        map[string]TableConfig `json:"processing_tables"`
}

func getDefaultConfigPaths() []string {
	paths := []string{
		"./maskdump.conf", // 1. Current directory
	}

	// 2. XDG_CONFIG_HOME (~/.config/maskdump/config)
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		paths = append(paths, filepath.Join(xdgConfig, defaultConfigDir, defaultConfigName))
	} else {
		paths = append(paths, filepath.Join(os.Getenv("HOME"), ".config", defaultConfigDir, defaultConfigName))
	}

	// 3. Global configuration
	paths = append(paths, "/etc/maskdump.conf")

	return paths
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

// LoadConfig loads the application configuration following this search hierarchy (in order of priority):
//  1. Explicit config path (-config flag, highest priority)
//  2. ./maskdump.conf (current directory)
//  3. $XDG_CONFIG_HOME/maskdump/config (~/.config/maskdump/config)
//  4. /etc/maskdump.conf (global config)
//  5. Built-in default values (lowest priority)
//
// If an explicit config path is provided but doesn't exist, returns an error.
// For other paths, falls through to next location in hierarchy if file not found.
// Returns nil if no config file is found (using defaults).
func LoadConfig(explicitPath string) error {
	// 1. Set default values first
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

	var configPath string
	var found bool

	// 2. If the config is explicitly specified, but does not exist, we return an error.
	if explicitPath != "" {
		if _, err := os.Stat(explicitPath); err == nil {
			configPath = explicitPath
			found = true
		} else {
			return fmt.Errorf("the specified config does not exist: %s", explicitPath)
		}
	} else {
		// 3. Check standard paths
		for _, path := range getDefaultConfigPaths() {
			if _, err := os.Stat(path); err == nil {
				configPath = path
				found = true
				break
			}
		}
	}

	// 4. If config found - load it
	if found {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return fmt.Errorf("error reading the config %s: %v", configPath, err)
		}

		var fileConfig Config
		if err := json.Unmarshal(data, &fileConfig); err != nil {
			return fmt.Errorf("invalid config file %s: %v", configPath, err)
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

		// ProcessingTables handling
		if len(fileConfig.ProcessingTables) > 0 {
			AppConfig.ProcessingTables = fileConfig.ProcessingTables
			ProcessingTables = fileConfig.ProcessingTables
		}
	}

	// 5. Validate all configurations
	if err := validateConfig(); err != nil {
		return fmt.Errorf("config validation failed: %v", err)
	}

	return nil
}

func validateConfig() error {
	// For the cache, we check the directory's availability and write permissions
	cacheDir := filepath.Dir(AppConfig.CachePath)
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create a cache directory %s: %v", cacheDir, err)
	}

	// Checking the ability to write to the cache file
	if err := checkFileAccess(AppConfig.CachePath, true); err != nil {
		return fmt.Errorf("cache access error: %v", err)
	}

	// Check white list files
	if AppConfig.EmailWhiteList != "" {
		if err := checkFileAccess(AppConfig.EmailWhiteList, false); err != nil {
			return fmt.Errorf("email white list error: %v", err)
		}
	}
	if AppConfig.PhoneWhiteList != "" {
		if err := checkFileAccess(AppConfig.PhoneWhiteList, false); err != nil {
			return fmt.Errorf("phone white list error: %v", err)
		}
	}

	// Check skip table list file
	if AppConfig.SkipInsertIntoTableList != "" {
		if err := checkFileAccess(AppConfig.SkipInsertIntoTableList, false); err != nil {
			return fmt.Errorf("skip table list error: %v", err)
		}
	}

	// Compile regular expressions
	var err error
	EmailRegex, err = regexp.Compile(AppConfig.EmailRegex)
	if err != nil {
		return fmt.Errorf("invalid email regex: %v", err)
	}

	PhoneRegex, err = regexp.Compile(AppConfig.PhoneRegex)
	if err != nil {
		return fmt.Errorf("invalid phone regex: %v", err)
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

	// Load skip table list
	SkipTableList, err = LoadSkipList(AppConfig.SkipInsertIntoTableList)
	if err != nil {
		return fmt.Errorf("failed to load skip table list: %v", err)
	}

	return nil
}

func checkFileAccess(path string, checkWrite bool) error {
	dir := filepath.Dir(path)

	// Checking/creating a directory
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("couldn't create a directory %s: %v", dir, err)
		}
	} else if err != nil {
		return fmt.Errorf("error checking the directory %s: %v", dir, err)
	}

	if checkWrite {
		// Checking the possibility of writing to the directory
		if err := os.WriteFile(path, []byte(""), 0644); err != nil {
			return fmt.Errorf("there is no write access to %s: %v", path, err)
		}
		// Delete the temporary file if it was created
		if info, err := os.Stat(path); err == nil && info.Size() == 0 {
			os.Remove(path)
		}
	} else if path != "" {
		// We only check the existence of files.
		if _, err := os.Stat(path); err != nil {
			return fmt.Errorf("the file does not exist or there is no access: %s: %v", path, err)
		}
	}

	return nil
}
