package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	defaultMaxBufferSize = 1024 * 1024 * 10 // 10MB
)

type Cache struct {
	Emails map[string]string `json:"emails"`
	Phones map[string]string `json:"phones"`
	sync.RWMutex
}

type MaskConfig struct {
	emailAlgorithm string
	phoneAlgorithm string
	cacheEnabled   bool
	configFile     string
}

type LogConfig struct {
	Path  string `json:"path"`
	Level string `json:"level"`
}

var (
	memoryLimit     int64
	currentMemUsage int64
	memMutex        sync.Mutex
)

type Logger struct {
	file  *os.File
	mu    sync.Mutex
	level int
}

const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
)

var logger *Logger

func (l *Logger) Check() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file == nil {
		return fmt.Errorf("log file not opened")
	}

	_, err := l.file.WriteString("\n")
	return err
}

func NewLogger(config LogConfig) (*Logger, error) {
	logPath := getDefaultLogPath(config.Path)

	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}

	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	level := LevelInfo
	switch strings.ToLower(config.Level) {
	case "debug":
		level = LevelDebug
	case "warn":
		level = LevelWarn
	case "error":
		level = LevelError
	}

	return &Logger{
		file:  file,
		level: level,
	}, nil
}

// getDefaultLogPath returns the path to the log file following this search hierarchy:
//  1. Explicitly specified path (if provided in config)
//  2. $XDG_STATE_HOME/maskdump/logs/maskdump.log (following XDG Base Directory Specification)
//  3. ~/.local/state/maskdump/logs/maskdump.log (fallback location)
//
// The function ensures the log directory structure exists before returning the path.
// Typical locations:
//   - /var/log/maskdump.log (if specified in config)
//   - /home/user/.local/state/maskdump/logs/maskdump.log (default)
//   - /home/user/.config/maskdump/logs/maskdump.log (if XDG_STATE_HOME not set)
func getDefaultLogPath(explicitPath string) string {
	if explicitPath != "" {
		return explicitPath
	}

	// XDG Base Directory Specification
	if stateHome := os.Getenv("XDG_STATE_HOME"); stateHome != "" {
		return filepath.Join(stateHome, "maskdump", "logs", "maskdump.log")
	}

	// Fallback to ~/.local/state/maskdump/logs/maskdump.log
	return filepath.Join(os.Getenv("HOME"), ".local", "state", "maskdump", "logs", "maskdump.log")
}

func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level <= LevelDebug {
		l.log("DEBUG", format, v...)
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	if l.level <= LevelInfo {
		l.log("INFO", format, v...)
	}
}

func (l *Logger) Warn(format string, v ...interface{}) {
	if l.level <= LevelWarn {
		l.log("WARN", format, v...)
	}
}

func (l *Logger) Error(format string, v ...interface{}) {
	if l.level <= LevelError {
		l.log("ERROR", format, v...)
	}
}

func (l *Logger) log(level, format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	msg := fmt.Sprintf(format, v...)
	timestamp := time.Now().Format("2006-01-02T15:04:05.000Z07:00")
	logEntry := fmt.Sprintf("%s [%s] %s\n", timestamp, level, msg)

	if _, err := l.file.Write([]byte(logEntry)); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write log: %v\n", err)
	}
}

func trackMemoryUsage() {
	for {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		memMutex.Lock()
		currentMemUsage = int64(m.Alloc)
		memMutex.Unlock()

		time.Sleep(500 * time.Millisecond)
	}
}

func checkMemoryLimit() bool {
	memMutex.Lock()
	defer memMutex.Unlock()
	return currentMemUsage > memoryLimit
}

func freeMemory(cache *Cache) {
	if cache == nil {
		return
	}

	// Flush cache to disk if possible
	if AppConfig.CachePath != "" {
		saveCache(cache)
	}

	// Clear internal caches
	cache.Lock()
	cache.Emails = make(map[string]string)
	cache.Phones = make(map[string]string)
	cache.Unlock()

	// Force garbage collection
	runtime.GC()
	debug.FreeOSMemory()
}

func loadCache() (*Cache, error) {
	cache := &Cache{
		Emails: make(map[string]string),
		Phones: make(map[string]string),
	}

	data, err := os.ReadFile(AppConfig.CachePath)
	if err != nil {
		return cache, nil
	}

	err = json.Unmarshal(data, cache)
	return cache, err
}

func saveCache(cache *Cache) error {
	cache.RLock()
	defer cache.RUnlock()

	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	return os.WriteFile(AppConfig.CachePath, data, 0644)
}

func parseFlags() MaskConfig {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError) // Сбрасываем флаги

	emailAlg := flag.String("mask-email", "", "Email masking algorithm (light-hash)")
	phoneAlg := flag.String("mask-phone", "", "Phone masking algorithm (light-mask)")
	noCache := flag.Bool("no-cache", false, "Disable caching")
	configFile := flag.String("config", "", "Path to config file")

	flag.Parse()

	return MaskConfig{
		emailAlgorithm: *emailAlg,
		phoneAlgorithm: *phoneAlg,
		cacheEnabled:   !*noCache,
		configFile:     *configFile,
	}
}

func validateAlgorithms(config MaskConfig) error {
	if config.emailAlgorithm != "" && config.emailAlgorithm != "light-hash" {
		return fmt.Errorf("unsupported email algorithm: %s", config.emailAlgorithm)
	}
	if config.phoneAlgorithm != "" && config.phoneAlgorithm != "light-mask" {
		return fmt.Errorf("unsupported phone algorithm: %s", config.phoneAlgorithm)
	}
	return nil
}

func parseTargetPositions(target string, length int) []int {
	var positions []int

	// Processing modifiers for email
	parts := strings.Split(target, ":")
	if len(parts) > 1 {
		target = parts[1]
	}

	// Range processing
	if strings.Contains(target, "-") {
		rangeParts := strings.Split(target, "-")
		start := 1
		end := length

		if rangeParts[0] != "" {
			start, _ = strconv.Atoi(rangeParts[0])
		}
		if rangeParts[1] != "" {
			end, _ = strconv.Atoi(rangeParts[1])
		}

		for i := start; i <= end && i <= length; i++ {
			positions = append(positions, i-1) // 0-based
		}
		return positions
	}

	// Tilde processing
	if strings.Contains(target, "~") {
		tildeParts := strings.Split(target, "~")
		keepStart := 0
		keepEnd := 0

		if tildeParts[0] != "" {
			keepStart, _ = strconv.Atoi(tildeParts[0])
		}
		if tildeParts[1] != "" {
			keepEnd, _ = strconv.Atoi(tildeParts[1])
		}

		for i := 1; i <= length; i++ {
			if (keepStart > 0 && i <= keepStart) ||
				(keepEnd > 0 && i > length-keepEnd) {
				continue
			}
			positions = append(positions, i-1)
		}
		return positions
	}

	// Item list processing
	posParts := strings.Split(target, ",")
	for _, p := range posParts {
		pos, _ := strconv.Atoi(p)
		if pos > 0 && pos <= length {
			positions = append(positions, pos-1)
		}
	}

	return positions
}

// value - string for masking
// positions - slice of positions to mask (0-based)
// maskValue - masking value (e.g. "*", "hash:6", "hash")
// typeMaskingInfo - type of masking (Email or Phone)
func applyMasking(value string, positions []int, maskValue string, typeMaskingInfo TypeMaskingInfo) string {
	runes := []rune(value)
	maskRunes := []rune{}

	// Preparing masking symbols
	var hash string
	if maskValue == "*" {
		for i := 0; i < len(positions); i++ {
			maskRunes = append(maskRunes, '*')
		}
	} else if strings.HasPrefix(maskValue, "hash") {
		hashParts := strings.Split(maskValue, ":")
		hashLen := 16 // by default
		if len(hashParts) > 1 {
			hashLen, _ = strconv.Atoi(hashParts[1])
		}

		if strings.HasPrefix(maskValue, "hash:") {
			tmpHash := md5.Sum([]byte(value))
			// Convert to hex string
			hashStr := hex.EncodeToString(tmpHash[:])
			// Take the first N characters
			hash = hashStr[:hashLen]
		} else {
			if typeMaskingInfo == Email && len(runes) > 0 {
				tmpHash := md5.Sum([]byte(value))
				hash = hex.EncodeToString(tmpHash[:])[:len(runes)]
			} else if typeMaskingInfo == Phone {
				tmpHash := sha256.Sum256([]byte(value))
				tmpHash2 := hex.EncodeToString(tmpHash[:])

				// We get only the digits for hashing
				digits := regexp.MustCompile(`\d`).FindAllString(tmpHash2, -1)
				hash = strings.Join(digits, "")
			}
		}

		maskRunes = []rune(hash)
	}

	// Application of mask
	result := ""
	if isContinuousSequence(positions) && typeMaskingInfo == Email && strings.HasPrefix(maskValue, "hash:") {
		result = replacePositions(value, positions, hash)
	} else {
		for i, pos := range positions {
			if pos >= 0 && pos < len(runes) && i < len(maskRunes) {
				runes[pos] = maskRunes[i]
			}
		}
		result = string(runes)
	}

	return result
}

func maskEmailWithRules(email string, cache *Cache) string {
	if _, ok := EmailWhiteList[email]; ok {
		return email
	}

	if cache != nil {
		cache.RLock()
		if masked, exists := cache.Emails[email]; exists {
			cache.RUnlock()
			return masked
		}
		cache.RUnlock()
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}

	localPart := parts[0]
	domainPart := parts[1]

	// Handling email by rule
	target := AppConfig.Masking.Email.Target
	value := AppConfig.Masking.Email.Value

	// Determine which parts need to be masked
	var positions []int
	typeMaskingInfo := Email
	if strings.Contains(target, "username:") {
		positions = parseTargetPositions(strings.TrimPrefix(target, "username:"), len(localPart))
		localPart = applyMasking(localPart, positions, value, typeMaskingInfo)
	} else if strings.Contains(target, "domain:") {
		positions = parseTargetPositions(strings.TrimPrefix(target, "domain:"), len(domainPart))
		domainPart = applyMasking(domainPart, positions, value, typeMaskingInfo)
	} else {
		positions = parseTargetPositions(target, len(email))
		masked := applyMasking(email, positions, value, typeMaskingInfo)
		return masked
	}

	masked := localPart + "@" + domainPart

	if cache != nil {
		cache.Lock()
		cache.Emails[email] = masked
		cache.Unlock()
	}

	return masked
}

func maskPhoneWithRules(phone string, cache *Cache) string {
	if _, ok := PhoneWhiteList[phone]; ok {
		return phone
	}

	if cache != nil {
		cache.RLock()
		if masked, exists := cache.Phones[phone]; exists {
			cache.RUnlock()
			return masked
		}
		cache.RUnlock()
	}

	// Handling phones according to the rules
	target := AppConfig.Masking.Phone.Target
	value := AppConfig.Masking.Phone.Value

	// We get only the digits for hashing
	digits := regexp.MustCompile(`\d`).FindAllString(phone, -1)
	digitStr := strings.Join(digits, "")

	positions := parseTargetPositions(target, len(digitStr))
	maskedDigits := applyMasking(digitStr, positions, value, Phone)

	// Restore the original format with replaced digits
	var result strings.Builder
	digitIndex := 0
	for _, c := range phone {
		if c >= '0' && c <= '9' {
			if digitIndex < len(maskedDigits) {
				result.WriteByte(maskedDigits[digitIndex])
				digitIndex++
			}
		} else {
			result.WriteRune(c)
		}
	}

	masked := result.String()

	if cache != nil {
		cache.Lock()
		cache.Phones[phone] = masked
		cache.Unlock()
	}

	return masked
}

// The function checks the character indices to be replaced.
// It concludes whether all listed indices are a continuous sequence.
func isContinuousSequence(positions []int) bool {
	if len(positions) <= 1 {
		// For 0 or 1 element, we consider the sequence to be continuous
		return true
	}

	// Check that all numbers are consecutive
	for i := 1; i < len(positions); i++ {
		if positions[i] != positions[i-1]+1 {
			return false
		}
	}

	return true
}

func replacePositions(value string, positions []int, hash string) string {
	if len(positions) == 0 {
		return value
	}

	// Convert string to runes for correct work with Unicode
	runes := []rune(value)
	var result []rune

	prev := 0
	for _, pos := range positions {
		if pos < 0 || pos >= len(runes) {
			continue // Skip invalid indexes
		}

		// Add a part of the string to the current position
		result = append(result, runes[prev:pos]...)

		// Update the previous position (skip the character to be deleted)
		prev = pos + 1
	}

	// Add the rest of the string
	result = append(result, runes[prev:]...)

	// Insert the replacement string before the first deleted character
	insertPos := positions[0]
	if insertPos < 0 {
		insertPos = 0
	} else if insertPos > len(runes) {
		insertPos = len(runes)
	}

	// Putting together the final line
	final := make([]rune, 0, len(result)+6)
	final = append(final, result[:insertPos]...)
	final = append(final, []rune(hash)...)
	final = append(final, result[insertPos:]...)

	return string(final)
}

// It's a basic function. It processes incoming strings.
// It starts the necessary masking functions according to the program settings.
func processLine(line string, config MaskConfig, cache *Cache, hasProcessingTables bool) string {
	if len(SkipTableList) > 0 {
		for table := range SkipTableList {
			if strings.HasPrefix(line, "INSERT INTO `"+table+"`") {
				return ""
			}
		}
	}

	if hasProcessingTables {
		ParseTableStructure(line)
	}

	if config.emailAlgorithm == "light-hash" {
		if hasProcessingTables {
			line = ProcessDumpLine(line, config, cache)
		} else if EmailRegex != nil {
			line = EmailRegex.ReplaceAllStringFunc(line, func(email string) string {
				return maskEmailWithRules(email, cache)
			})
		}
	}

	if config.phoneAlgorithm == "light-mask" {
		if hasProcessingTables {
			line = ProcessDumpLine(line, config, cache)
		} else if PhoneRegex != nil {
			line = PhoneRegex.ReplaceAllStringFunc(line, func(phone string) string {
				return maskPhoneWithRules(phone, cache)
			})
		}
	}

	return line
}

// The function prepares the required values of the setting variables.
// Keeps track of memory and cache. Reads the input buffer, starts processing of incoming strings.
// Outputs to the output buffer the result after masking and ignoring the specified tables.
func main() {
	config := parseFlags()

	// Temporary logger for initialization errors
	initLogger, err := NewLogger(LogConfig{
		Path:  "/tmp/maskdump_init.log", // temporary file
		Level: "error",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create init logger: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		initLogger.Close()
		os.Remove("/tmp/maskdump_init.log")
	}()

	// Load configuration
	if err := LoadConfig(config.configFile); err != nil {
		initLogger.Error("Config error: %v", err)
		fmt.Fprintf(os.Stderr, "Config error: %v\n", err)
		os.Exit(1)
	}

	// Checking if the log directory exists
	if err := os.MkdirAll(filepath.Dir(AppConfig.Logging.Path), 0755); err != nil {
		initLogger.Error("Failed to create log directory: %v", err)
		fmt.Fprintf(os.Stderr, "Failed to create log directory: %v\n", err)
		os.Exit(1)
	}

	// Now we initialize the main logger from the config
	logger, err = NewLogger(AppConfig.Logging)
	if err != nil {
		initLogger.Error("Failed to initialize main logger: %v", err)
		fmt.Fprintf(os.Stderr, "Failed to initialize main logger: %v\n", err)
		os.Exit(1)
	}

	if err := logger.Check(); err != nil {
		fmt.Fprintf(os.Stderr, "Logger check failed: %v\n", err)
		os.Exit(1)
	}

	defer logger.Close()

	logger.Info("Starting maskdump with config from %s", config.configFile)
	logger.Debug("Config settings: %+v", AppConfig)

	if err := validateAlgorithms(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	var cache *Cache
	if config.cacheEnabled {
		var err error
		cache, err = loadCache()
		if err != nil {
			logger.Warn("Cache load warning: %v", err)
		}
	}

	// Set memory limit
	memoryLimit = int64(AppConfig.MemoryLimitMB) * 1024 * 1024
	go trackMemoryUsage()

	reader := bufio.NewReaderSize(os.Stdin, defaultMaxBufferSize)
	writer := bufio.NewWriterSize(os.Stdout, defaultMaxBufferSize)
	defer writer.Flush()

	// Checking if there are any processing tables
	hasProcessingTables := len(ProcessingTables) > 0

	lineCount := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			logger.Error("Error reading input: %v", err)
			os.Exit(1)
		}

		maskedLine := processLine(line, config, cache, hasProcessingTables)
		if strings.TrimSpace(maskedLine) != "" {
			_, err = writer.WriteString(maskedLine)
			if err != nil {
				logger.Error("Error writing output: %v", err)
				os.Exit(1)
			}
		}

		lineCount++
		if lineCount%AppConfig.CacheFlushCount == 0 {
			if checkMemoryLimit() {
				freeMemory(cache)
			}
		}
	}

	if config.cacheEnabled && cache != nil {
		if err := saveCache(cache); err != nil {
			logger.Warn("Cache save warning: %v", err)
		}
	}

	logger.Info("Processing completed, processed %d lines", lineCount)
}
