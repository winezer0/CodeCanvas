package engine

import (
	"context"
	"fmt"
	"github.com/winezer0/codecanvas/internal/model"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestEmbeddedRules verifies that every embedded rule can be triggered by a minimal test case.
func TestEmbeddedRules(t *testing.T) {
	// Initialize engine with embedded rules
	e, err := NewCanvasEngine("")
	if err != nil {
		t.Fatalf("Failed to initialize engine: %v", err)
	}

	if len(e.rules) == 0 {
		t.Fatal("No rules loaded from embedded assets")
	}

	for _, rule := range e.rules {
		t.Run(fmt.Sprintf("%s-%s", rule.Language, rule.Name), func(t *testing.T) {
			// Create a temp directory for this rule
			tmpDir, err := os.MkdirTemp("", "rule_test_*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			// Setup test environment based on the rule
			if !setupTestEnv(t, tmpDir, rule) {
				t.Skipf("Skipping rule %s: could not setup test environment (complex rule?)", rule.Name)
			}

			// Run detection
			ctx := context.Background()
			// We must provide the language the rule expects, otherwise it might skip
			languages := []string{rule.Language}

			// Build file index
			index, err := buildTestIndex(tmpDir)
			if err != nil {
				t.Fatalf("Failed to build file index: %v", err)
			}

			result, err := e.DetectFrameworks(ctx, index, languages)
			if err != nil {
				t.Fatalf("DetectFrameworks failed: %v", err)
			}

			// Verify detection
			found := false
			var items []model.DetectedItem
			if rule.Type == "framework" {
				items = result.Frameworks
			} else {
				items = result.Components
			}

			for _, item := range items {
				if item.Name == rule.Name {
					found = true
					break
				}
			}

			if !found {
				// 重新提取 level 用于日志记录
				var level *model.RuleSet
				if l, ok := rule.Levels["L1"]; ok {
					level = l
				} else if l, ok := rule.Levels["L2"]; ok {
					level = l
				} else if l, ok := rule.Levels["L3"]; ok {
					level = l
				}

				paths := []string{}
				contains := []string{}
				if level != nil {
					paths = level.Paths
					contains = level.Contains
				}

				t.Errorf("Rule %s was not detected. Created files in %s. Paths: %v. Contains: %v", rule.Name, tmpDir, paths, contains)
			}
		})
	}
}

// setupTestEnv creates files in the temp dir to satisfy the rule.
// Returns true if setup was successful, false if the rule is too complex to auto-mock.
func setupTestEnv(t *testing.T, dir string, rule *model.FrameworkRuleDefinition) bool {
	// Prefer L1, then L2, then L3
	var level *model.RuleSet
	if l, ok := rule.Levels["L1"]; ok {
		level = l
	} else if l, ok := rule.Levels["L2"]; ok {
		level = l
	} else if l, ok := rule.Levels["L3"]; ok {
		level = l
	} else {
		return false
	}

	if len(level.Paths) == 0 {
		return false
	}

	// Pick the first path
	pathPattern := level.Paths[0]

	// Handle globs roughly - just pick a simple filename that matches
	filename := pathPattern

	// 如果模式以 / 结尾，它期望目录存在。
	// 在其中创建一个文件以便被索引。
	if strings.HasSuffix(pathPattern, "/") {
		filename = pathPattern + "index.php"
	} else if pathPattern == "*.go" {
		filename = "main.go"
	} else if pathPattern == "*.js" {
		filename = "index.js"
	} else if pathPattern == "*.json" {
		filename = "package.json"
	} else if pathPattern == "*.php" {
		filename = "index.php"
	} else if pathPattern == "*.py" {
		filename = "app.py"
	} else if pathPattern == "pom.xml" {
		filename = "pom.xml"
	} else if pathPattern == "go.mod" {
		filename = "go.mod"
	} else if strings.Contains(pathPattern, "*") {
		// Replace * with something concrete
		filename = strings.ReplaceAll(pathPattern, "*", "test_file")
	}

	// Construct content
	content := ""
	if len(level.Contains) > 0 {
		// Just concatenate all required strings
		for _, s := range level.Contains {
			content += s + "\n"
		}
	} else if level.ExtractVersionFromText != nil && level.ExtractVersionFromText.Pattern != "" {
		// This is harder to mock automatically without regex reversing
		// But usually Contains is also present or we can try to guess.
		// For now, if Contains is empty but regex is there, we might fail.
		// But most rules have Contains.
	}

	// Create the file
	fullPath := filepath.Join(dir, filename)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		t.Errorf("Failed to create dirs: %v", err)
		return false
	}

	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		t.Errorf("Failed to write file: %v", err)
		return false
	}

	return true
}
