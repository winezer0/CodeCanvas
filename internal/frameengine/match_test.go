package frameengine

import (
	"testing"
)

// TestExtractVersion tests the extractVersion function with various scenarios
func TestExtractVersion(t *testing.T) {
	testCases := []struct {
		name     string
		content  string
		patterns []string
		expected string
	}{
		{
			name:     "Single valid regex",
			content:  `{"version": "1.2.3"}`,
			patterns: []string{`"version"\s*:\s*"([^"]+)"`},
			expected: "1.2.3",
		},
		{
			name:    "Multiple regexes with priority",
			content: `version="2.0.0"\nappVersion: 1.0.0`,
			patterns: []string{
				`version=\"([\d.]+)\"`,
				`appVersion:\s*([\d.]+)`,
			},
			expected: "2.0.0", // First match wins
		},
		{
			name:     "Semantic version with prerelease",
			content:  `version: "3.1.0-alpha.1"`,
			patterns: []string{`version:\s*"([\w.-]+)"`},
			expected: "3.1.0-alpha.1",
		},
		{
			name:     "Major.minor version format",
			content:  `VERSION=4.5`,
			patterns: []string{`VERSION=(\d+\.\d+)`},
			expected: "4.5",
		},
		{
			name:    "Invalid regex should be skipped",
			content: `version: "5.6.7"`,
			patterns: []string{
				`[invalid-regex`, // This is invalid
				`version:\s*"([\d.]+)"`,
			},
			expected: "5.6.7", // Should skip invalid regex and use next one
		},
		{
			name:     "No matching regex",
			content:  `no version here`,
			patterns: []string{`version:\s*"([\d.]+)"`},
			expected: "", // Should return empty string
		},
		{
			name:     "Case insensitive match",
			content:  `Version: "6.7.8"`,
			patterns: []string{`(?i:version):\s*"([\d.]+)"`},
			expected: "6.7.8",
		},
		{
			name:     "XML format version",
			content:  `<project><version>7.8.9</version></project>`,
			patterns: []string{`<version>([^<]+)</version>`},
			expected: "7.8.9",
		},
		{
			name:     "YAML format with comments",
			content:  `# This is a comment\nversion: 8.9.0 # Another comment`,
			patterns: []string{`version:\s*([\d.]+)`},
			expected: "8.9.0",
		},
		{
			name:     "Multiple version occurrences",
			content:  `oldVersion: 1.0.0\nnewVersion: 9.0.0\nlatestVersion: 10.0.0`,
			patterns: []string{`newVersion:\s*([\d.]+)`},
			expected: "9.0.0", // Should match the first occurrence of the pattern
		},
		{
			name:     "Version in URL path",
			content:  `dependency: "https://example.com/lib/11.2.3/file.jar"`,
			patterns: []string{`lib/(\d+\.\d+\.\d+)/`},
			expected: "11.2.3",
		},
		{
			name:     "Version with build metadata",
			content:  `version: "12.3.4+build.123"`,
			patterns: []string{`version:\s*"([\d.]+)`}, // Note: This regex won't capture the build metadata
			expected: "12.3.4",
		},
		{
			name:     "Empty content",
			content:  ``,
			patterns: []string{`version:\s*"([\d.]+)"`},
			expected: "",
		},
		{
			name:     "Content with no version info",
			content:  `This file has no version information at all`,
			patterns: []string{`version:\s*"([\d.]+)"`},
			expected: "",
		},
		{
			name:     "Complex version string with letters",
			content:  `release: "v13.5.7"`,
			patterns: []string{`release:\s*"v?([\d.]+)"`},
			expected: "13.5.7",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := extractVersion([]byte(tc.content), tc.patterns)
			if result != tc.expected {
				t.Errorf("extractVersion() = %q, want %q\nContent: %q\nPatterns: %v", result, tc.expected, tc.content, tc.patterns)
			}
		})
	}
}

// TestExtractVersionPerformance tests the performance of extractVersion function
// with a large number of patterns to ensure it handles them efficiently
func TestExtractVersionPerformance(t *testing.T) {
	// Skip this test in short mode
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	content := `version: "1.2.3"`
	// Create 100 patterns, most of which are invalid or won't match
	patterns := make([]string, 100)
	for i := 0; i < 99; i++ {
		if i%2 == 0 {
			patterns[i] = `invalid[regex` // Invalid regex
		} else {
			patterns[i] = `unlikelypattern\d+` // Unlikely to match
		}
	}
	patterns[99] = `version:\s*"([\d.]+)"` // Should match

	result := extractVersion([]byte(content), patterns)
	if result != "1.2.3" {
		t.Errorf("Expected version 1.2.3, got %s", result)
	}
}
