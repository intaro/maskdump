package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	emailRegex = regexp.MustCompile(`\b[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}\b`)
)

func maskEmail(email string) string {
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

func processLine(line string) string {
	return emailRegex.ReplaceAllStringFunc(line, maskEmail)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	for scanner.Scan() {
		line := scanner.Text()
		maskedLine := processLine(line)
		writer.WriteString(maskedLine + "\n")
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}