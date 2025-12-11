package canvas

import (
	"context"
	"fmt"
	"github.com/winezer0/codecanvas/internal/analyzer"
	"github.com/winezer0/codecanvas/internal/engine"
	"github.com/winezer0/codecanvas/internal/model"
	"time"
)

// Analyze performs a full analysis and returns a CanvasReport.
func Analyze(path string, rulesDir string, version string) (*model.CanvasReport, error) {
	ctx := context.Background()

	// Analyze code profile
	az := analyzer.NewCodeAnalyzer()
	profile, index, err := az.AnalyzeCodeProfile(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("error analyzing code profile: %v", err)
	}

	// Create rule engine
	ruleEngine, err := engine.NewCanvasEngine(rulesDir)
	if err != nil {
		return nil, fmt.Errorf("error loading rules: %v", err)
	}

	// 获取检测到的语言
	languages := make([]string, 0, len(profile.Languages))
	for _, lang := range profile.Languages {
		languages = append(languages, lang.Name)
	}

	// 扩展语言（例如 TSX -> JavaScript）以确保规则匹配
	languages = expandLanguages(languages)

	// 检测框架和组件
	detectionResult, err := ruleEngine.DetectFrameworks(ctx, index, languages)
	if err != nil {
		return nil, fmt.Errorf("error detecting frameworks and components: %v", err)
	}

	// 生成分析报告
	report := &model.CanvasReport{
		CodeProfile: *profile,
		Detection:   *detectionResult,
		Timestamp:   time.Now(),
		Version:     version,
	}
	return report, nil
}

func expandLanguages(langs []string) []string {
	seen := make(map[string]bool)
	for _, l := range langs {
		seen[l] = true
	}

	// 检查 JS 系列
	if seen["TypeScript"] || seen["TSX"] || seen["JSX"] {
		if !seen["JavaScript"] {
			langs = append(langs, "JavaScript")
			seen["JavaScript"] = true
		}
	}

	return langs
}
