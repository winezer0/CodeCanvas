package embeds

import (
	"testing"
)

func TestEmbeddedRulesVerification(t *testing.T) {
	rules := EmbeddedFrameRules()

	if len(rules) == 0 {
		t.Fatal("Expected embedded rules to be loaded, got 0")
	}

	// Map to track unique rules
	// Key: Name + Type + Language
	seen := make(map[string]bool)

	// List of expected frameworks to verify they are loaded
	expectedFrameworks := map[string]bool{
		"ThinkPHP":    false,
		"Laravel":     false,
		"Yii":         false,
		"Django":      false,
		"FastAPI":     false,
		"Flask":       false,
		"React":       false,
		"Vue.js":      false,
		"Spring Boot": false,
		"Quarkus":     false,
	}

	for _, rule := range rules {
		// Verify basic fields
		if rule.Name == "" {
			t.Errorf("Rule with empty name found")
		}
		if rule.Type == "" {
			t.Errorf("Rule '%s' has empty type", rule.Name)
		}
		if rule.Language == "" {
			t.Errorf("Rule '%s' has empty language", rule.Name)
		}

		// Check for duplicates
		key := rule.Name + "|" + rule.Type + "|" + rule.Language
		if seen[key] {
			t.Errorf("Duplicate rule found: %s", key)
		}
		seen[key] = true

		// Mark as found if expected
		if _, ok := expectedFrameworks[rule.Name]; ok {
			expectedFrameworks[rule.Name] = true
		}
	}

	// Verify all expected frameworks were found
	for name, found := range expectedFrameworks {
		if !found {
			t.Errorf("Expected framework '%s' not found in embedded rules", name)
		}
	}

	t.Logf("Successfully verified %d embedded rules", len(rules))
}
