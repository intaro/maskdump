package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

const cacheFileName = ".maskdump_cache.json"

var (
	emailRegex = regexp.MustCompile(`\b[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}\b`)
	phoneRegex = regexp.MustCompile(`(?:\+7|7|8)?(?:[\s\-\(\)]*\d){10}`)
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
}

func loadCache() (*Cache, error) {
	cache := &Cache{
		Emails: make(map[string]string),
		Phones: make(map[string]string),
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return cache, nil
	}

	cachePath := filepath.Join(homeDir, cacheFileName)
	data, err := os.ReadFile(cachePath)
	if err != nil {
		return cache, nil
	}

	err = json.Unmarshal(data, cache)
	return cache, err
}

func saveCache(cache *Cache) error {
	cache.RLock()
	defer cache.RUnlock()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	cachePath := filepath.Join(homeDir, cacheFileName)
	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	return os.WriteFile(cachePath, data, 0644)
}

func parseFlags() MaskConfig {
	emailAlg := flag.String("mask-email", "", "Email masking algorithm (light-hash)")
	phoneAlg := flag.String("mask-phone", "", "Phone masking algorithm (light-mask)")
	noCache := flag.Bool("no-cache", false, "Disable caching")
	flag.Parse()

	return MaskConfig{
		emailAlgorithm: *emailAlg,
		phoneAlgorithm: *phoneAlg,
		cacheEnabled:   !*noCache,
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
	if config.emailAlgorithm == "light-hash" {
		line = emailRegex.ReplaceAllStringFunc(line, func(email string) string {
			return maskEmailLightHash(email, cache)
		})
	}
	if config.phoneAlgorithm == "light-mask" {
		line = phoneRegex.ReplaceAllStringFunc(line, func(phone string) string {
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

	var cache *Cache
	if config.cacheEnabled {
		var err error
		cache, err = loadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cache load warning: %v\n", err)
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	for scanner.Scan() {
		line := scanner.Text()
		maskedLine := processLine(line, config, cache)
		writer.WriteString(maskedLine + "\n")
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	if config.cacheEnabled && cache != nil {
		if err := saveCache(cache); err != nil {
			fmt.Fprintf(os.Stderr, "Cache save warning: %v\n", err)
		}
	}
}
