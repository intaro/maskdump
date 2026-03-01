package main

import "testing"

func TestMaskEmailWithRulesEdgeCases(t *testing.T) {
	withTestGlobals(t, func() {
		setupMaskingDefaults(t)

		// Invalid email should be returned as-is.
		invalid := "not-an-email"
		if got := maskEmailWithRules(invalid, nil); got != invalid {
			t.Fatalf("expected invalid email to be unchanged, got %s", got)
		}

		// Empty string should be returned as-is.
		if got := maskEmailWithRules("", nil); got != "" {
			t.Fatalf("expected empty email to be unchanged, got %s", got)
		}

		// Whitelist should bypass masking.
		EmailWhiteList["keep@example.com"] = struct{}{}
		if got := maskEmailWithRules("keep@example.com", nil); got != "keep@example.com" {
			t.Fatalf("expected whitelisted email to be unchanged, got %s", got)
		}

		// Username-only masking should preserve domain.
		AppConfig.Masking.Email.Target = "username:2-"
		AppConfig.Masking.Email.Value = "hash:6"
		masked := maskEmailWithRules("user@example.com", nil)
		if len(masked) == 0 || masked == "user@example.com" {
			t.Fatalf("expected masked email, got %s", masked)
		}
		if masked[len(masked)-len("example.com"):] != "example.com" {
			t.Fatalf("expected domain to be preserved, got %s", masked)
		}

		// Cache should return the same masked value for the same input.
		cache := &Cache{Emails: map[string]string{}, Phones: map[string]string{}}
		first := maskEmailWithRules("cache@example.com", cache)
		second := maskEmailWithRules("cache@example.com", cache)
		if first != second {
			t.Fatalf("expected cached email to be stable, got %s vs %s", first, second)
		}
	})
}

func TestMaskPhoneWithRulesEdgeCases(t *testing.T) {
	withTestGlobals(t, func() {
		setupMaskingDefaults(t)

		// Non-digit string should be returned as-is.
		invalid := "no-phone-here"
		if got := maskPhoneWithRules(invalid, nil); got != invalid {
			t.Fatalf("expected invalid phone to be unchanged, got %s", got)
		}

		// Whitelist should bypass masking.
		PhoneWhiteList["+7 (900) 111-22-33"] = struct{}{}
		if got := maskPhoneWithRules("+7 (900) 111-22-33", nil); got != "+7 (900) 111-22-33" {
			t.Fatalf("expected whitelisted phone to be unchanged, got %s", got)
		}

		// Target outside digit length should keep digits unchanged.
		AppConfig.Masking.Phone.Target = "10-"
		AppConfig.Masking.Phone.Value = "hash"
		short := "123"
		if got := maskPhoneWithRules(short, nil); got != short {
			t.Fatalf("expected short phone to be unchanged, got %s", got)
		}

		// Cache should return the same masked value for the same input.
		AppConfig.Masking.Phone.Target = "2-"
		AppConfig.Masking.Phone.Value = "hash"
		cache := &Cache{Emails: map[string]string{}, Phones: map[string]string{}}
		first := maskPhoneWithRules("+7 (900) 111-22-33", cache)
		second := maskPhoneWithRules("+7 (900) 111-22-33", cache)
		if first != second {
			t.Fatalf("expected cached phone to be stable, got %s vs %s", first, second)
		}
	})
}
