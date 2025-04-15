package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	emailRegex = regexp.MustCompile(`\b[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}\b`)
	phoneRegex = regexp.MustCompile(`(?:\+7|7|8)?(?:[\s\-\(\)]*\d){10}`)
)

type MaskConfig struct {
	emailAlgorithm string
	phoneAlgorithm string
}

func parseFlags() MaskConfig {
	emailAlg := flag.String("mask-email", "", "Email masking algorithm (light-hash)")
	phoneAlg := flag.String("mask-phone", "", "Phone masking algorithm (light-mask)")
	flag.Parse()

	return MaskConfig{
		emailAlgorithm: *emailAlg,
		phoneAlgorithm: *phoneAlg,
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

func maskEmailLightHash(email string) string {
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

	return firstChar + hashedRest + "@" + domainPart
}

func maskPhoneLightMask(phone string) string {
	// Извлекаем все цифры из номера
	digits := regexp.MustCompile(`\d`).FindAllString(phone, -1)
	if len(digits) < 10 {
		return phone
	}

	// Получаем SHA256 хэш оригинального номера
	hash := sha256.Sum256([]byte(phone))
	hashStr := hex.EncodeToString(hash[:])

	// Выбираем первые 6 цифр из хэша
	hashDigits := make([]string, 0)
	for _, c := range hashStr {
		if c >= '0' && c <= '9' {
			hashDigits = append(hashDigits, string(c))
			if len(hashDigits) == 6 {
				break
			}
		}
	}

	// Заменяем только указанные позиции (2,3,5,6,8,10 цифры)
	positions := []int{1, 2, 4, 5, 7, 9} // 0-based индексы
	for i, pos := range positions {
		if pos < len(digits) && i < len(hashDigits) {
			digits[pos] = hashDigits[i]
		}
	}

	// Восстанавливаем оригинальный формат с изменёнными цифрами
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

	return result.String()
}

func processLine(line string, config MaskConfig) string {
	if config.emailAlgorithm == "light-hash" {
		line = emailRegex.ReplaceAllStringFunc(line, maskEmailLightHash)
	}
	if config.phoneAlgorithm == "light-mask" {
		line = phoneRegex.ReplaceAllStringFunc(line, maskPhoneLightMask)
	}
	return line
}

func main() {
	config := parseFlags()
	if err := validateAlgorithms(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	for scanner.Scan() {
		line := scanner.Text()
		maskedLine := processLine(line, config)
		writer.WriteString(maskedLine + "\n")
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
}
