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

func parseTargetPositions(target string, length int) []int {
	var positions []int

	// Обработка модификаторов для email
	parts := strings.Split(target, ":")
	//modifier := ""
	if len(parts) > 1 {
		//modifier = parts[0]
		target = parts[1]
	}

	// Обработка диапазонов
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
			positions = append(positions, i-1) // переводим в 0-based
		}
		return positions
	}

	// Обработка тильды
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

	// Обработка списка позиций
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

	// Подготовка маскирующих символов
	var hash string
	if maskValue == "*" {
		for i := 0; i < len(positions); i++ {
			maskRunes = append(maskRunes, '*')
		}
	} else if strings.HasPrefix(maskValue, "hash") {
		hashParts := strings.Split(maskValue, ":")
		hashLen := 6 // по умолчанию
		if len(hashParts) > 1 {
			hashLen, _ = strconv.Atoi(hashParts[1])
		}

		if strings.HasPrefix(maskValue, "hash:") {
			tmpHash := md5.Sum([]byte(value))
			// Конвертируем в hex строку
			hashStr := hex.EncodeToString(tmpHash[:])
			// Берём первые N символов
			hash = hashStr[:hashLen]
		} else {
			if typeMaskingInfo == Email && len(runes) > 0 {
				tmpHash := md5.Sum([]byte(value))
				hash = hex.EncodeToString(tmpHash[:])[:len(runes)]
			} else if typeMaskingInfo == Phone {
				tmpHash := sha256.Sum256([]byte(value))
				tmpHash2 := hex.EncodeToString(tmpHash[:])

				// Получаем только цифры для хеширования
				digits := regexp.MustCompile(`\d`).FindAllString(tmpHash2, -1)
				hash = strings.Join(digits, "")
			}
		}

		maskRunes = []rune(hash)
	}

	// Применение маски
	if typeMaskingInfo == Email && strings.HasPrefix(maskValue, "hash:") {
		var firstSymbols []rune
		if len(runes) >= positions[0] {
			firstSymbols = runes[:positions[0]]
		} else {
			firstSymbols = runes // если меньше 2 элементов, берем все что есть
		}
		runes = firstSymbols
		runes = append(runes, maskRunes...)
	} else {
		for i, pos := range positions {
			if pos >= 0 && pos < len(runes) && i < len(maskRunes) {
				runes[pos] = maskRunes[i]
			}
		}
	}

	return string(runes)
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

	// Обработка email по правилам
	target := AppConfig.Masking.Email.Target
	value := AppConfig.Masking.Email.Value

	// Определяем какие части нужно маскировать
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

	// Обработка телефона по правилам
	target := AppConfig.Masking.Phone.Target
	value := AppConfig.Masking.Phone.Value

	// Получаем только цифры для хеширования
	digits := regexp.MustCompile(`\d`).FindAllString(phone, -1)
	digitStr := strings.Join(digits, "")

	positions := parseTargetPositions(target, len(digitStr))
	maskedDigits := applyMasking(digitStr, positions, value, Phone)

	// Восстанавливаем оригинальный формат с заменёнными цифрами
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
			return maskEmailWithRules(email, cache)
		})
	}
	if config.phoneAlgorithm == "light-mask" {
		line = PhoneRegex.ReplaceAllStringFunc(line, func(phone string) string {
			return maskPhoneWithRules(phone, cache)
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
