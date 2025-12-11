package canvas

import (
	"sort"
	"strings"

	"github.com/winezer0/codecanvas/internal/model"
)

// AnalyzeDirectory 对指定目录进行代码画板分析，返回分类后的结果。
// rulesDir 可选，指定自定义规则目录。如果为空，则仅使用内置规则。
func AnalyzeDirectory(path string, rulesDir string) (*model.AnalysisResult, error) {
	// 调用 canvas 包的核心分析函数
	// 注意：canvas.Analyze 返回的是 CanvasReport，包含了更详细的信息
	report, err := Analyze(path, rulesDir, "0.0.0")
	if err != nil {
		return nil, err
	}

	// 转换语言列表为 Map 以便快速查找统计信息
	langStats := make(map[string]model.LanguageInfo)
	for _, l := range report.CodeProfile.Languages {
		langStats[l.Name] = l
	}

	result := &model.AnalysisResult{
		Languages:             report.CodeProfile.Languages,
		DesktopLanguages:      report.CodeProfile.DesktopLanguages,
		MainFrontendLanguages: getTopLanguages(report.CodeProfile.FrontendLanguages, langStats, nil, 3),
		FrontendLanguages:     report.CodeProfile.FrontendLanguages,
		MainBackendLanguages:  getTopLanguages(report.CodeProfile.BackendLanguages, langStats, nil, 3),
		BackendLanguages:      report.CodeProfile.BackendLanguages,
		Frameworks:            report.Detection.Frameworks,
		Components:            report.Detection.Components,
	}

	return result, nil
}

// getTopLanguages 根据代码行数和文件数对语言进行排序并返回前 N 个
func getTopLanguages(candidates []string, stats map[string]model.LanguageInfo, exclude []string, limit int) []string {
	// 过滤需要排除的语言
	excludeMap := make(map[string]bool)
	for _, e := range exclude {
		excludeMap[strings.ToLower(e)] = true
	}

	var validLangs []model.LanguageInfo
	for _, name := range candidates {
		if excludeMap[strings.ToLower(name)] {
			continue
		}
		if info, ok := stats[name]; ok {
			validLangs = append(validLangs, info)
		}
	}

	// 排序：优先代码行数，其次文件数
	sort.Slice(validLangs, func(i, j int) bool {
		if validLangs[i].CodeLines != validLangs[j].CodeLines {
			return validLangs[i].CodeLines > validLangs[j].CodeLines
		}
		return validLangs[i].Files > validLangs[j].Files
	})

	// 取前 N 个
	var result []string
	count := 0
	for _, l := range validLangs {
		if count >= limit {
			break
		}
		result = append(result, l.Name)
		count++
	}
	return result
}
