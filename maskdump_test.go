package main

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestMaskEmailWithRules(t *testing.T) {
	withTestGlobals(t, func() {
		setupMaskingDefaults(t)

		testCases := []struct {
			email  string
			target string
			value  string
		}{
			{email: "test@example.com", target: "username:2-", value: "hash:6"},
			{email: "test@example.com", target: "domain:2-", value: "hash:6"},
			{email: "test@example.com", target: "2-", value: "hash"},
			{email: "test@example.com", target: "username:2~", value: "hash"},
			{email: "test@example.com", target: "domain:2~", value: "hash"},
			{email: "test@example.com", target: "2~1", value: "*"},
			{email: "a@b.co", target: "username:2-", value: "*"},
		}

		for _, tc := range testCases {
			t.Run(tc.email+"_"+tc.target+"_"+tc.value, func(t *testing.T) {
				AppConfig.Masking.Email.Target = tc.target
				AppConfig.Masking.Email.Value = tc.value

				masked := maskEmailWithRules(tc.email, nil)

				parts := strings.Split(tc.email, "@")
				if len(parts) != 2 {
					t.Fatalf("invalid test email: %s", tc.email)
				}
				origDomain := parts[1]
				positionsCount := 0

				if strings.HasPrefix(tc.target, "username:") {
					positions := parseTargetPositions(strings.TrimPrefix(tc.target, "username:"), len(parts[0]))
					positionsCount = len(positions)
					if len(positions) > 0 && masked == tc.email {
						t.Fatalf("expected email to be masked: %s", tc.email)
					}
					if !strings.HasSuffix(masked, origDomain) {
						t.Fatalf("expected domain to be preserved: %s", masked)
					}
					if !strings.Contains(masked, "@") {
						t.Fatalf("expected @ to remain, got %s", masked)
					}
				} else if strings.HasPrefix(tc.target, "domain:") {
					positions := parseTargetPositions(strings.TrimPrefix(tc.target, "domain:"), len(origDomain))
					positionsCount = len(positions)
					if len(positions) > 0 && masked == tc.email {
						t.Fatalf("expected email to be masked: %s", tc.email)
					}
					if !strings.Contains(masked, "@") {
						t.Fatalf("expected @ to remain, got %s", masked)
					}
				} else {
					positions := parseTargetPositions(tc.target, len(tc.email))
					positionsCount = len(positions)
					if len(positions) > 0 && masked == tc.email {
						t.Fatalf("expected email to be masked: %s", tc.email)
					}
				}

				if tc.value == "*" {
					if positionsCount > 0 && !strings.Contains(masked, "*") {
						t.Fatalf("expected '*' masking for %s", tc.email)
					}
					if len(masked) != len(tc.email) {
						t.Fatalf("expected same length, got %d vs %d", len(masked), len(tc.email))
					}
				}
			})
		}

		EmailWhiteList["keep@example.com"] = struct{}{}
		if got := maskEmailWithRules("keep@example.com", nil); got != "keep@example.com" {
			t.Fatalf("expected whitelisted email to be unchanged, got %s", got)
		}

		if got := maskEmailWithRules("not-an-email", nil); got != "not-an-email" {
			t.Fatalf("expected invalid email to be unchanged, got %s", got)
		}

		cache := &Cache{Emails: map[string]string{}, Phones: map[string]string{}}
		first := maskEmailWithRules("cache@example.com", cache)
		second := maskEmailWithRules("cache@example.com", cache)
		if first != second {
			t.Fatalf("expected cached email to be stable, got %s vs %s", first, second)
		}
	})
}

func TestMaskPhoneWithRules(t *testing.T) {
	withTestGlobals(t, func() {
		setupMaskingDefaults(t)

		testCases := []struct {
			phone  string
			target string
			value  string
		}{
			{phone: "+7 (900) 111-22-33", target: "2-6", value: "hash"},
			{phone: "+7 (900) 111-22-33", target: "2-6", value: "*"},
			{phone: "+79001112233", target: "2-", value: "hash"},
			{phone: "+79001112233", target: "2-", value: "*"},
			{phone: "79001112233", target: "~3", value: "hash"},
			{phone: "79001112233", target: "~3", value: "*"},
			{phone: "8 (900) 111-22-33", target: "1,3,5,7,9", value: "hash"},
			{phone: "8 (900) 111-22-33", target: "1,3,5,7,9", value: "*"},
			{phone: "8-900-111-22-33", target: "2~2", value: "hash"},
			{phone: "8-900-111-22-33", target: "2~2", value: "*"},
			{phone: "+7(900)111-22-33", target: "2~", value: "hash"},
			{phone: "+7(900)111-22-33", target: "2~", value: "*"},
			{phone: "+7 900 111-22-33", target: "1,2,3", value: "hash"},
			{phone: "+7 900 111-22-33", target: "1,2,3", value: "*"},
		}

		for _, tc := range testCases {
			t.Run(tc.phone+"_"+tc.target+"_"+tc.value, func(t *testing.T) {
				AppConfig.Masking.Phone.Target = tc.target
				AppConfig.Masking.Phone.Value = tc.value

				masked := maskPhoneWithRules(tc.phone, nil)

				if len(masked) != len(tc.phone) {
					t.Fatalf("expected same length, got %d vs %d", len(masked), len(tc.phone))
				}
				if stripDigits(tc.phone) != stripDigitsAndStars(masked) {
					t.Fatalf("expected formatting to be preserved: %s -> %s", tc.phone, masked)
				}
				if countDigits(masked)+strings.Count(masked, "*") != countDigits(tc.phone) {
					t.Fatalf("unexpected masked digit count: %s -> %s", tc.phone, masked)
				}
				if tc.value == "*" {
					if !strings.Contains(masked, "*") {
						t.Fatalf("expected '*' masking for %s", tc.phone)
					}
					if masked == tc.phone {
						t.Fatalf("expected phone to be masked, got unchanged: %s", tc.phone)
					}
				}
			})
		}
	})
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
	t.Cleanup(func() {
		if err := os.Remove(AppConfig.CachePath); err != nil && !os.IsNotExist(err) {
			t.Fatalf("failed to remove cache file %s: %v", AppConfig.CachePath, err)
		}
	})

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
