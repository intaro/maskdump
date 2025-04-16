package main

import (
	"os"
	"reflect"
	"testing"
)

func TestMaskEmailLightHash(t *testing.T) {
	cache := &Cache{
		Emails: make(map[string]string),
		Phones: make(map[string]string),
	}

	email := "test@example.com"
	maskedEmail := maskEmailLightHash(email, cache)
	if maskedEmail != "t1c52bd@example.com" {
		t.Errorf("maskEmailLightHash not working correctly. Input: %s, Output: %s", email, maskedEmail)
	}

	email2 := "t@example.com"
	maskedEmail2 := maskEmailLightHash(email2, cache)
	if maskedEmail2 != "td41d8c@example.com" {
		t.Errorf("maskEmailLightHash not working correctly. Input: %s, Output: %s", email2, maskedEmail2)
	}

	// Test with caching
	email3 := "test3@example.com"
	maskedEmail3 := maskEmailLightHash(email3, cache)
	maskedEmail3Again := maskEmailLightHash(email3, cache)
	if maskedEmail3 != maskedEmail3Again {
		t.Errorf("Cache is not working for emails")
	}
}

func TestMaskPhoneLightMask(t *testing.T) {
	cache := &Cache{
		Emails: make(map[string]string),
		Phones: make(map[string]string),
	}

	phone := "+79001112233"
	maskedPhone := maskPhoneLightMask(phone, cache)
	if maskedPhone != "+73703214293" {
		t.Errorf("maskPhoneLightMask not working correctly. Input: %s, Output: %s", phone, maskedPhone)
	}

	phone2 := "8-900-111-22-33"
	maskedPhone2 := maskPhoneLightMask(phone2, cache)
	if maskedPhone2 != "8-530-231-82-83" {
		t.Errorf("maskPhoneLightMask not working correctly. Input: %s, Output: %s", phone2, maskedPhone2)
	}

	phone3 := "9001112233"
	maskedPhone3 := maskPhoneLightMask(phone3, cache)
	if maskedPhone3 != "9411282433" {
		t.Errorf("maskPhoneLightMask not working correctly. Input: %s, Output: %s", phone3, maskedPhone3)
	}

	// Test with caching
	phone4 := "89001112233"
	maskedPhone4 := maskPhoneLightMask(phone4, cache)
	maskedPhone4Again := maskPhoneLightMask(phone4, cache)
	if maskedPhone4 != maskedPhone4Again {
		t.Errorf("Cache is not working for phones")
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
