// Package engine 提供了 CodeCanvas 的框架和组件检测功能。
package engine

import (
	"context"
	"fmt"
	"github.com/winezer0/codecanvas/internal/embeds"
	"github.com/winezer0/codecanvas/internal/model"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// CanvasEngine 实现框架和组件检测功能。
type CanvasEngine struct {
	rules          []*model.FrameworkRuleDefinition
	frameworkRules map[string]*model.FrameworkRuleDefinition
	componentRules map[string]*model.FrameworkRuleDefinition
}

// NewCanvasEngine 创建一个新的规则引擎实例，默认加载嵌入式规则。
// 如果提供了规则目录（rulesDir），则会从该目录加载规则，并将其与嵌入式规则合并，
// 其中用户定义的规则将覆盖具有相同名称的嵌入式规则。
func NewCanvasEngine(rulesDir string) (*CanvasEngine, error) {
	engine := &CanvasEngine{
		rules:          []*model.FrameworkRuleDefinition{},
		frameworkRules: make(map[string]*model.FrameworkRuleDefinition),
		componentRules: make(map[string]*model.FrameworkRuleDefinition),
	}

	// 首先加载嵌入式规则
	engine.loadEmbeddedRules()

	// 如果提供了规则目录，则加载用户定义的规则并与嵌入式规则合并
	if rulesDir != "" {
		err := engine.loadRulesFromDirectory(rulesDir)
		if err != nil {
			return nil, err
		}
	}

	return engine, nil
}

// loadEmbeddedRules 将默认的嵌入式规则加载到规则引擎中。
func (e *CanvasEngine) loadEmbeddedRules() {
	embeddedRules := embeds.EmbeddedCanvas()

	// 将嵌入式规则添加到引擎中
	for _, rule := range embeddedRules {
		e.addRule(rule)
	}
}

// addRule 向引擎添加单个规则，替换具有相同名称的任何现有规则。
func (e *CanvasEngine) addRule(rule *model.FrameworkRuleDefinition) {
	// 检查规则是否已存在
	existingIndex := -1
	for i, r := range e.rules {
		if r.Name == rule.Name && r.Type == rule.Type && r.Language == rule.Language {
			existingIndex = i
			break
		}
	}

	// 如果规则存在，则替换它
	if existingIndex != -1 {
		e.rules[existingIndex] = rule
	} else {
		// 否则，添加新规则
		e.rules = append(e.rules, rule)
	}

	// 更新规则映射
	switch rule.Type {
	case model.RuleTypeFramework:
		e.frameworkRules[rule.Name] = rule
	case model.RuleTypeComponent:
		e.componentRules[rule.Name] = rule
	}
}

// loadRulesFromDirectory 从给定目录加载所有 YAML 规则文件。
func (e *CanvasEngine) loadRulesFromDirectory(rulesDir string) error {
	// 读取目录中的所有 YAML 文件
	yamlFiles, err := filepath.Glob(filepath.Join(rulesDir, "*.yml"))
	if err != nil {
		return err
	}

	for _, yamlFile := range yamlFiles {
		rules, err := e.loadRulesFromFile(yamlFile)
		if err != nil {
			return err
		}

		// 使用 addRule 方法将规则添加到引擎中以处理规则合并
		for _, rule := range rules {
			e.addRule(rule)
		}
	}

	return nil
}

// loadRulesFromFile 从单个 YAML 文件加载规则，支持单文档数组格式和多文档格式。
func (e *CanvasEngine) loadRulesFromFile(filePath string) ([]*model.FrameworkRuleDefinition, error) {
	// 读取 YAML 文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// 首先尝试解析为单文档数组格式
	var rulesArray []*model.FrameworkRuleDefinition
	if err := yaml.Unmarshal(data, &rulesArray); err == nil {
		// 成功解析为数组格式
		return rulesArray, nil
	}

	// 如果数组格式解析失败，尝试多文档格式
	yamlReader := strings.NewReader(string(data))
	decoder := yaml.NewDecoder(yamlReader)

	var rules []*model.FrameworkRuleDefinition

	// 解码 YAML 文件中的每个文档
	for {
		var rule model.FrameworkRuleDefinition
		err := decoder.Decode(&rule)
		if err != nil {
			if err == io.EOF {
				// 已到达文件末尾
				break
			}
			return nil, err
		}

		// 将解码后的规则添加到列表中
		rules = append(rules, &rule)
	}

	return rules, nil
}

// DetectFrameworks 根据加载的规则检测给定目录中的框架和组件。
// 使用文件索引进行加速。
func (e *CanvasEngine) DetectFrameworks(ctx context.Context, index *model.FileIndex, languages []string) (*model.DetectionResult, error) {
	result := &model.DetectionResult{
		Frameworks: []model.DetectedItem{},
		Components: []model.DetectedItem{},
	}

	// 创建索引匹配器
	matcher := NewIndexMatcher(index)

	// 按检测到的语言过滤规则
	filteredRules := e.filterRulesByLanguages(languages)

	// 文件内容缓存
	fileContentCache := make(map[string][]byte)

	// 辅助函数：获取文件内容（带缓存）
	getFileContent := func(path string) ([]byte, error) {
		if content, ok := fileContentCache[path]; ok {
			return content, nil
		}

		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		// 检查文件大小，如果太大（例如 > 5MB），可能需要跳过或只读部分
		stat, err := f.Stat()
		if err == nil && stat.Size() > 5*1024*1024 { // 5MB limit
			// 对于大文件，我们只读取前 1MB
			content, err := io.ReadAll(io.LimitReader(f, 1*1024*1024))
			if err != nil {
				return nil, err
			}
			fileContentCache[path] = content
			return content, nil
		}

		content, err := io.ReadAll(f)
		if err != nil {
			return nil, err
		}
		fileContentCache[path] = content
		return content, nil
	}

	// 按置信度级别顺序处理每个规则 (L1 优先, 然后 L2, 然后 L3)
	for _, rule := range filteredRules {
		// 按顺序处理检测级别 (L1, L2, L3)
		levels := []string{"L1", "L2", "L3"}
		matchFound := false
		for _, level := range levels {
			ruleSet, ok := rule.Levels[level]
			if !ok {
				continue
			}

			// 检查是否有任何路径匹配
			for _, pathPattern := range ruleSet.Paths {
				// 查找匹配的文件
				matches, err := matcher.FindFiles(pathPattern)
				if err != nil {
					continue
				}

				for _, matchPath := range matches {
					// 准备文件内容
					var content []byte
					if len(ruleSet.Contains) > 0 || ruleSet.ExtractVersionFromText != nil {
						var err error
						content, err = getFileContent(matchPath)
						if err != nil {
							continue
						}
					}

					// 检查文件是否包含所有必需的关键字
					if !containsAllKeywords(content, ruleSet.Contains) {
						continue
					}

					// 提取版本（如果可能）
					version := ""
					if ruleSet.ExtractVersionFromText != nil {
						version = extractVersionFromText(content, ruleSet.ExtractVersionFromText.Pattern)
					} else if ruleSet.ExtractVersionFromPath != nil {
						version = extractVersionFromPath(matchPath, ruleSet.ExtractVersionFromPath.Pattern)
					}

					confidence := mapConfidenceLevel(level)
					evidence := fmt.Sprintf("Found %s", matchPath)
					if len(ruleSet.Contains) > 0 {
						evidence += fmt.Sprintf(" containing %s", strings.Join(ruleSet.Contains, ", "))
					}

					item := model.DetectedItem{
						Name:       rule.Name,
						Type:       rule.Type,
						Language:   rule.Language,
						Version:    version,
						Category:   rule.Category,
						Confidence: confidence,
						Evidence:   evidence,
					}

					// 根据规则类型创建检测结果
					switch rule.Type {
					case model.RuleTypeFramework:
						result.Frameworks = append(result.Frameworks, item)

					case model.RuleTypeComponent:
						result.Components = append(result.Components, item)
					}

					// 首次匹配后停止处理此规则
					matchFound = true
					break
				}

				if matchFound {
					// 首次匹配后停止处理此级别
					break
				}
			}

			if matchFound {
				// 首次匹配级别后停止处理此规则
				break
			}
		}
	}

	return result, nil
}

// filterRulesByLanguages 过滤规则，只包含与检测到的语言匹配的规则。
func (e *CanvasEngine) filterRulesByLanguages(languages []string) []*model.FrameworkRuleDefinition {
	var filtered []*model.FrameworkRuleDefinition

	for _, rule := range e.rules {
		// 检查规则的语言是否在检测到的语言中
		for _, lang := range languages {
			if rule.Language == lang {
				filtered = append(filtered, rule)
				break
			}
		}
	}

	return filtered
}

// GetSupportedFrameworks 从规则中提取有关所有可检测框架的元数据。
func (e *CanvasEngine) GetSupportedFrameworks() []model.FrameworkMetadata {
	var frameworks []model.FrameworkMetadata

	for _, rule := range e.frameworkRules {
		// 提取级别信息
		levels := make(map[string]string)
		for level, ruleSet := range rule.Levels {
			if len(ruleSet.Paths) > 0 {
				levels[level] = ruleSet.Paths[0] // 使用第一个路径作为代表
			}
		}

		framework := model.FrameworkMetadata{
			Name:     rule.Name,
			Language: rule.Language,
			Levels:   levels,
		}

		frameworks = append(frameworks, framework)
	}

	// 按名称对框架进行排序
	sort.Slice(frameworks, func(i, j int) bool {
		return frameworks[i].Name < frameworks[j].Name
	})

	return frameworks
}

// GetSupportedComponents 从规则中提取有关所有可检测 组件 的元数据
func (e *CanvasEngine) GetSupportedComponents() []model.ComponentMetadata {
	var components []model.ComponentMetadata
	for _, rule := range e.componentRules {
		// 提取级别信息
		levels := make(map[string]string)
		for level, ruleSet := range rule.Levels {
			if len(ruleSet.Paths) > 0 {
				levels[level] = ruleSet.Paths[0] // 使用第一个路径作为代表
			}
		}
		component := model.ComponentMetadata{
			Name:     rule.Name,
			Language: rule.Language,
			Levels:   levels,
		}

		components = append(components, component)
	}
	// Sort components by name
	sort.Slice(components, func(i, j int) bool {
		return components[i].Name < components[j].Name
	})
	return components
}
