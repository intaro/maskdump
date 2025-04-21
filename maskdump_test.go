package main

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestMaskEmailWithRules(t *testing.T) {
	typeMaskingInfo := Email

	// Example #1 target and value
	email1 := "test@example.com"
	target1 := "username:2-"
	value1 := "hash:6"

	parts1 := strings.Split(email1, "@")
	if len(parts1) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email1)
	}

	localPart1 := parts1[0]
	domainPart1 := parts1[1]

	positions1 := parseTargetPositions(strings.TrimPrefix(target1, "username:"), len(localPart1))
	maskLocalPart1 := applyMasking(localPart1, positions1, value1, typeMaskingInfo)

	maskedEmail1 := maskLocalPart1 + "@" + domainPart1
	if maskedEmail1 != "t098f6b@example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email1, target1, value1, maskedEmail1)
	}

	// Example #2 target and value
	email2 := "test@example.com"
	target2 := "domain:2-"
	value2 := "hash:6"

	parts2 := strings.Split(email2, "@")
	if len(parts2) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email2)
	}

	localPart2 := parts2[0]
	domainPart2 := parts2[1]

	positions2 := parseTargetPositions(strings.TrimPrefix(target2, "domain:"), len(domainPart2))
	maskDomainPart2 := applyMasking(domainPart2, positions2, value2, typeMaskingInfo)

	maskedEmail2 := localPart2 + "@" + maskDomainPart2
	if maskedEmail2 != "test@e5ababd" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email2, target2, value2, maskedEmail2)
	}

	// Example #3 target and value
	email3 := "test@example.com"
	target3 := "2-"
	value3 := "hash:6"

	positions3 := parseTargetPositions(target3, len(email3))
	maskedEmail3 := applyMasking(email3, positions3, value3, typeMaskingInfo)

	if maskedEmail3 != "t55502f" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email3, target3, value3, maskedEmail3)
	}

	// Example #4 target and value
	email4 := "test@example.com"
	target4 := "username:2-"
	value4 := "hash"

	parts4 := strings.Split(email4, "@")
	if len(parts4) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email4)
	}

	localPart4 := parts4[0]
	domainPart4 := parts4[1]

	positions4 := parseTargetPositions(strings.TrimPrefix(target4, "username:"), len(localPart4))
	maskLocalPart4 := applyMasking(localPart4, positions4, value4, typeMaskingInfo)

	maskedEmail4 := maskLocalPart4 + "@" + domainPart4
	if maskedEmail4 != "t098@example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email4, target4, value4, maskedEmail4)
	}

	// Example #5 target and value
	email5 := "test@example.com"
	target5 := "domain:2-"
	value5 := "hash"

	parts5 := strings.Split(email5, "@")
	if len(parts5) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email5)
	}

	localPart5 := parts5[0]
	domainPart5 := parts5[1]

	positions5 := parseTargetPositions(strings.TrimPrefix(target5, "domain:"), len(domainPart5))
	maskDomainPart5 := applyMasking(domainPart5, positions5, value5, typeMaskingInfo)

	maskedEmail5 := localPart5 + "@" + maskDomainPart5
	if maskedEmail5 != "test@e5ababd603b" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email5, target5, value5, maskedEmail5)
	}

	// Example #6 target and value
	email6 := "test@example.com"
	target6 := "2-"
	value6 := "hash"

	positions6 := parseTargetPositions(target6, len(email6))
	maskedEmail6 := applyMasking(email6, positions6, value6, typeMaskingInfo)

	if maskedEmail6 != "t55502f40dc8b7c7" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email6, target6, value6, maskedEmail6)
	}
}

func TestValidateAlgorithms(t *testing.T) {
	config1 := MaskConfig{
		emailAlgorithm: "light-hash",
		phoneAlgorithm: "light-mask",
	}
	err1 := validateAlgorithms(config1)
	if err1 != nil {
		t.Errorf("validateAlgorithms failed for valid config: %v", err1)
	}

	config2 := MaskConfig{
		emailAlgorithm: "invalid",
		phoneAlgorithm: "light-mask",
	}
	err2 := validateAlgorithms(config2)
	if err2 == nil {
		t.Errorf("validateAlgorithms should have failed for invalid email algorithm")
	}

	config3 := MaskConfig{
		emailAlgorithm: "light-hash",
		phoneAlgorithm: "invalid",
	}
	err3 := validateAlgorithms(config3)
	if err3 == nil {
		t.Errorf("validateAlgorithms should have failed for invalid phone algorithm")
	}
}

func TestLoadSaveCache(t *testing.T) {
	// Setup
	AppConfig.CachePath = "test_cache.json"
	defer os.Remove(AppConfig.CachePath) // Clean up after test

	cache := &Cache{
		Emails: map[string]string{"test@example.com": "masked@example.com"},
		Phones: map[string]string{"123-456-7890": "masked"},
	}

	// Save cache
	err := saveCache(cache)
	if err != nil {
		t.Fatalf("Error saving cache: %v", err)
	}

	// Load cache
	loadedCache, err := loadCache()
	if err != nil {
		t.Fatalf("Error loading cache: %v", err)
	}

	// Compare
	if !reflect.DeepEqual(cache.Emails, loadedCache.Emails) || !reflect.DeepEqual(cache.Phones, loadedCache.Phones) {
		t.Errorf("Cache data mismatch. Expected: %v, Got: %v", cache, loadedCache)
	}
}
