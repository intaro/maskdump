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
	"regexp"
	"runtime"
	"runtime/debug"
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

var (
	memoryLimit     int64
	currentMemUsage int64
	memMutex        sync.Mutex
)

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

func maskEmailLightHash(email string, cache *Cache) string {
	// Проверяем белый список
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

	if len(localPart) == 0 {
		return email
	}

	firstChar := string(localPart[0])
	rest := localPart[1:]

	hash := md5.Sum([]byte(rest))
	hashedRest := hex.EncodeToString(hash[:])[:6]

	masked := firstChar + hashedRest + "@" + domainPart

	if cache != nil {
		cache.Lock()
		cache.Emails[email] = masked
		cache.Unlock()
	}

	return masked
}

func maskPhoneLightMask(phone string, cache *Cache) string {
	// Проверяем белый список
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

	digits := regexp.MustCompile(`\d`).FindAllString(phone, -1)
	if len(digits) < 10 {
		return phone
	}

	hash := sha256.Sum256([]byte(phone))
	hashStr := hex.EncodeToString(hash[:])

	hashDigits := make([]string, 0)
	for _, c := range hashStr {
		if c >= '0' && c <= '9' {
			hashDigits = append(hashDigits, string(c))
			if len(hashDigits) == 6 {
				break
			}
		}
	}

	positions := []int{1, 2, 4, 5, 7, 9}
	for i, pos := range positions {
		if pos < len(digits) && i < len(hashDigits) {
			digits[pos] = hashDigits[i]
		}
	}

	var result strings.Builder
	digitIndex := 0
	for _, c := range phone {
		if c >= '0' && c <= '9' {
			if digitIndex < len(digits) {
				result.WriteString(digits[digitIndex])
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

func processLine(line string, config MaskConfig, cache *Cache) string {
	if len(SkipTableList) > 0 {
		for table := range SkipTableList {
			if strings.HasPrefix(line, "INSERT INTO `"+table+"`") {
				return ""
			}
		}
	}

	if config.emailAlgorithm == "light-hash" {
		line = EmailRegex.ReplaceAllStringFunc(line, func(email string) string {
			return maskEmailLightHash(email, cache)
		})
	}
	if config.phoneAlgorithm == "light-mask" {
		line = PhoneRegex.ReplaceAllStringFunc(line, func(phone string) string {
			return maskPhoneLightMask(phone, cache)
		})
	}
	return line
}

func main() {
	config := parseFlags()
	if err := validateAlgorithms(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Load configuration
	if err := LoadConfig(config.configFile); err != nil {
		fmt.Fprintf(os.Stderr, "Config error: %v\n", err)
		os.Exit(1)
	}

	var cache *Cache
	if config.cacheEnabled {
		var err error
		cache, err = loadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cache load warning: %v\n", err)
		}
	}

	// Set memory limit
	memoryLimit = int64(AppConfig.MemoryLimitMB) * 1024 * 1024
	go trackMemoryUsage()

	reader := bufio.NewReaderSize(os.Stdin, defaultMaxBufferSize)
	writer := bufio.NewWriterSize(os.Stdout, defaultMaxBufferSize)
	defer writer.Flush()

	lineCount := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}

		maskedLine := processLine(line, config, cache)
		if strings.TrimSpace(maskedLine) != "" {
			_, err = writer.WriteString(maskedLine)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
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
			fmt.Fprintf(os.Stderr, "Cache save warning: %v\n", err)
		}
	}
}
