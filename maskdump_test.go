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

	// Example #7 target and value
	email7 := "test@example.com"
	target7 := "username:2~"
	value7 := "hash:6"

	parts7 := strings.Split(email7, "@")
	if len(parts7) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email7)
	}

	localPart7 := parts7[0]
	domainPart7 := parts7[1]

	positions7 := parseTargetPositions(strings.TrimPrefix(target7, "username:"), len(localPart7))
	maskLocalPart7 := applyMasking(localPart7, positions7, value7, typeMaskingInfo)

	maskedEmail7 := maskLocalPart7 + "@" + domainPart7
	if maskedEmail7 != "te098f6b@example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email7, target7, value7, maskedEmail7)
	}

	// Example #8 target and value
	email8 := "test@example.com"
	target8 := "domain:2~"
	value8 := "hash:6"

	parts8 := strings.Split(email8, "@")
	if len(parts8) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email8)
	}

	localPart8 := parts8[0]
	domainPart8 := parts8[1]

	positions8 := parseTargetPositions(strings.TrimPrefix(target8, "domain:"), len(domainPart8))
	maskDomainPart8 := applyMasking(domainPart8, positions8, value8, typeMaskingInfo)

	maskedEmail8 := localPart8 + "@" + maskDomainPart8
	if maskedEmail8 != "test@ex5ababd" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email8, target8, value8, maskedEmail8)
	}

	// Example #9 target and value
	email9 := "test@example.com"
	target9 := "2~"
	value9 := "hash:6"

	positions9 := parseTargetPositions(target9, len(email9))
	maskedEmail9 := applyMasking(email9, positions9, value9, typeMaskingInfo)

	if maskedEmail9 != "te55502f" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email9, target9, value9, maskedEmail9)
	}

	// Example #10 target and value
	email10 := "test@example.com"
	target10 := "username:2~"
	value10 := "hash"

	parts10 := strings.Split(email10, "@")
	if len(parts10) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email10)
	}

	localPart10 := parts10[0]
	domainPart10 := parts10[1]

	positions10 := parseTargetPositions(strings.TrimPrefix(target10, "username:"), len(localPart10))
	maskLocalPart10 := applyMasking(localPart10, positions10, value10, typeMaskingInfo)

	maskedEmail10 := maskLocalPart10 + "@" + domainPart10
	if maskedEmail10 != "te09@example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email10, target10, value10, maskedEmail10)
	}

	// Example #11 target and value
	email11 := "test@example.com"
	target11 := "domain:2~"
	value11 := "hash"

	parts11 := strings.Split(email11, "@")
	if len(parts11) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email11)
	}

	localPart11 := parts11[0]
	domainPart11 := parts11[1]

	positions11 := parseTargetPositions(strings.TrimPrefix(target11, "domain:"), len(domainPart11))
	maskDomainPart11 := applyMasking(domainPart11, positions11, value11, typeMaskingInfo)

	maskedEmail11 := localPart11 + "@" + maskDomainPart11
	if maskedEmail11 != "test@ex5ababd603" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email11, target11, value11, maskedEmail11)
	}

	// Example #12 target and value
	email12 := "test@example.com"
	target12 := "2~"
	value12 := "hash"

	positions12 := parseTargetPositions(target12, len(email12))
	maskedEmail12 := applyMasking(email12, positions12, value12, typeMaskingInfo)

	if maskedEmail12 != "te55502f40dc8b7c" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email12, target12, value12, maskedEmail12)
	}

	// Example #13 target and value
	email13 := "test@example.com"
	target13 := "username:-2"
	value13 := "hash:6"

	parts13 := strings.Split(email13, "@")
	if len(parts13) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email13)
	}

	localPart13 := parts13[0]
	domainPart13 := parts13[1]

	positions13 := parseTargetPositions(strings.TrimPrefix(target13, "username:"), len(localPart13))
	maskLocalPart13 := applyMasking(localPart13, positions13, value13, typeMaskingInfo)

	maskedEmail13 := maskLocalPart13 + "@" + domainPart13
	if maskedEmail13 != "098f6bst@example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email13, target13, value13, maskedEmail13)
	}

	// Example #14 target and value
	email14 := "test@example.com"
	target14 := "domain:-2"
	value14 := "hash:6"

	parts14 := strings.Split(email14, "@")
	if len(parts14) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email14)
	}

	localPart14 := parts14[0]
	domainPart14 := parts14[1]

	positions14 := parseTargetPositions(strings.TrimPrefix(target14, "domain:"), len(domainPart14))
	maskDomainPart14 := applyMasking(domainPart14, positions14, value14, typeMaskingInfo)

	maskedEmail14 := localPart14 + "@" + maskDomainPart14
	if maskedEmail14 != "test@5ababdample.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email14, target14, value14, maskedEmail14)
	}

	// Example #15 target and value
	email15 := "test@example.com"
	target15 := "-2"
	value15 := "hash:6"

	positions15 := parseTargetPositions(target15, len(email15))
	maskedEmail15 := applyMasking(email15, positions15, value15, typeMaskingInfo)

	if maskedEmail15 != "55502fst@example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email15, target15, value15, maskedEmail15)
	}

	// Example #16 target and value
	email16 := "test@example.com"
	target16 := "username:-2"
	value16 := "hash"

	parts16 := strings.Split(email16, "@")
	if len(parts16) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email16)
	}

	localPart16 := parts16[0]
	domainPart16 := parts16[1]

	positions16 := parseTargetPositions(strings.TrimPrefix(target16, "username:"), len(localPart16))
	maskLocalPart16 := applyMasking(localPart16, positions16, value16, typeMaskingInfo)

	maskedEmail16 := maskLocalPart16 + "@" + domainPart16
	if maskedEmail16 != "09st@example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email16, target16, value16, maskedEmail16)
	}

	// Example #17 target and value
	email17 := "test@example.com"
	target17 := "domain:-2"
	value17 := "hash"

	parts17 := strings.Split(email17, "@")
	if len(parts17) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email17)
	}

	localPart17 := parts17[0]
	domainPart17 := parts17[1]

	positions17 := parseTargetPositions(strings.TrimPrefix(target17, "domain:"), len(domainPart17))
	maskDomainPart17 := applyMasking(domainPart17, positions17, value17, typeMaskingInfo)

	maskedEmail17 := localPart17 + "@" + maskDomainPart17
	if maskedEmail17 != "test@5aample.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email17, target17, value17, maskedEmail17)
	}

	// Example #18 target and value
	email18 := "test@example.com"
	target18 := "-2"
	value18 := "hash"

	positions18 := parseTargetPositions(target18, len(email18))
	maskedEmail18 := applyMasking(email18, positions18, value18, typeMaskingInfo)

	if maskedEmail18 != "55st@example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email18, target18, value18, maskedEmail18)
	}

	// Example #19 target and value
	email19 := "test@example.com"
	target19 := "username:~2"
	value19 := "hash:6"

	parts19 := strings.Split(email19, "@")
	if len(parts19) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email19)
	}

	localPart19 := parts19[0]
	domainPart19 := parts19[1]

	positions19 := parseTargetPositions(strings.TrimPrefix(target19, "username:"), len(localPart19))
	maskLocalPart19 := applyMasking(localPart19, positions19, value19, typeMaskingInfo)

	maskedEmail19 := maskLocalPart19 + "@" + domainPart19
	if maskedEmail19 != "098f6bst@example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email19, target19, value19, maskedEmail19)
	}

	// Example #20 target and value
	email20 := "test@example.com"
	target20 := "domain:~2"
	value20 := "hash:6"

	parts20 := strings.Split(email20, "@")
	if len(parts20) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email20)
	}

	localPart20 := parts20[0]
	domainPart20 := parts20[1]

	positions20 := parseTargetPositions(strings.TrimPrefix(target20, "domain:"), len(domainPart20))
	maskDomainPart20 := applyMasking(domainPart20, positions20, value20, typeMaskingInfo)

	maskedEmail20 := localPart20 + "@" + maskDomainPart20
	if maskedEmail20 != "test@5ababdom" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email20, target20, value20, maskedEmail20)
	}

	// Example #21 target and value
	email21 := "test@example.com"
	target21 := "~2"
	value21 := "hash:6"

	positions21 := parseTargetPositions(target21, len(email21))
	maskedEmail21 := applyMasking(email21, positions21, value21, typeMaskingInfo)

	if maskedEmail21 != "55502fom" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email21, target21, value21, maskedEmail21)
	}

	// Example #22 target and value
	email22 := "test@example.com"
	target22 := "username:~2"
	value22 := "hash"

	parts22 := strings.Split(email22, "@")
	if len(parts22) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email22)
	}

	localPart22 := parts22[0]
	domainPart22 := parts22[1]

	positions22 := parseTargetPositions(strings.TrimPrefix(target22, "username:"), len(localPart22))
	maskLocalPart22 := applyMasking(localPart22, positions22, value22, typeMaskingInfo)

	maskedEmail22 := maskLocalPart22 + "@" + domainPart22
	if maskedEmail22 != "09st@example.com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email22, target22, value22, maskedEmail22)
	}

	// Example #23 target and value
	email23 := "test@example.com"
	target23 := "domain:~2"
	value23 := "hash"

	parts23 := strings.Split(email23, "@")
	if len(parts23) != 2 {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s", email23)
	}

	localPart23 := parts23[0]
	domainPart23 := parts23[1]

	positions23 := parseTargetPositions(strings.TrimPrefix(target23, "domain:"), len(domainPart23))
	maskDomainPart23 := applyMasking(domainPart23, positions23, value23, typeMaskingInfo)

	maskedEmail23 := localPart23 + "@" + maskDomainPart23
	if maskedEmail23 != "test@5ababd603om" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email23, target23, value23, maskedEmail23)
	}

	// Example #24 target and value
	email24 := "test@example.com"
	target24 := "~2"
	value24 := "hash"

	positions24 := parseTargetPositions(target24, len(email24))
	maskedEmail24 := applyMasking(email24, positions24, value24, typeMaskingInfo)

	if maskedEmail24 != "55502f40dc8b7com" {
		t.Errorf("maskEmailWithRules not working correctly. Input: %s, mask target: %s, mask value: %s, Output: %s", email24, target24, value24, maskedEmail24)
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
