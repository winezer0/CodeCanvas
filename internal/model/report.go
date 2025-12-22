package model

import "time"

// 代码画板包定义了 CodeCanvas 的核心数据结构
// ，CodeCanvas 是一款轻量级的代码性能分析和框架检测引擎
// 本版本仅专注于技术栈识别——不包含安全元数据。
// 常量定义
const (
	// 规则类型
	RuleTypeFramework = "framework"
	RuleTypeComponent = "component"

	// 代码所处的应用类别
	CategoryFrontend = "frontend"
	CategoryBackend  = "backend"
	CategoryDesktop  = "desktop"
	CategoryOther    = "other"
)

var AllCategory = []string{CategoryFrontend, CategoryBackend, CategoryDesktop, CategoryOther}

// CanvasReport 最终分析报告
type CanvasReport struct {
	CodeProfile CodeProfile   `json:"code_profile"`
	Detection   DetectionInfo `json:"detection"`
	Timestamp   time.Time     `json:"timestamp"`
	Version     string        `json:"version"`
}
type CodeProfile struct {
	Path              string     `json:"path"`
	TotalFiles        int        `json:"total_files"`
	TotalLines        int        `json:"total_lines"`
	ErrorFiles        int        `json:"error_files"`        // Number of files that failed to process
	LanguageInfos     []LangInfo `json:"language_infos"`     // 所有检测到的语言的完整列表
	FrontendLanguages []string   `json:"frontend_languages"` // 例如: ["TypeScript", "JavaScript"]
	BackendLanguages  []string   `json:"backend_languages"`  // 例如: ["Java", "Go"]
	DesktopLanguages  []string   `json:"desktop_languages"`  // 例如: ["C#", "C++"]
	OtherLanguages    []string   `json:"other_languages"`    // 例如: ["JSON", "YAML"]
	Languages         []string   `json:"languages"`
	ExpandLanguages   []string   `json:"expand_languages"`
}

// DetectionInfo 框架与组件识别结果 包含已检测到的框架和组件的列表。
type DetectionInfo struct {
	Frameworks []DetectedItem `json:"frameworks"`
	Components []DetectedItem `json:"components"`
}

// DetectedItem  框架与组件识别结果代表了一项已检测到的技术项目（框架或组件）。
type DetectedItem struct {
	Name     string `json:"name"`     // 例如: "gin", "log4j-core", "wails"
	Type     string `json:"type"`     // "framework" 或 "component"
	Language string `json:"language"` // 例如: "Go", "Java", "JavaScript"
	Version  string `json:"version"`  // 版本字符串，可能为空
	Category string `json:"category"` // "frontend" | "backend" | "desktop"
	Evidence string `json:"evidence"` // 人类可读的检测原因
}

// LangInfo  某一编程语言或标记语言的详细统计数据。
type LangInfo struct {
	Name         string `json:"name"` // 例如: "Java", "YAML"
	Files        int    `json:"files"`
	CodeLines    int    `json:"code_lines"`
	CommentLines int    `json:"comment_lines"`
	BlankLines   int    `json:"blank_lines"`
}

// AnalysisResult 包含了 CodeCanvas 分析的完整结果
type AnalysisResult struct {
	// 语言信息列表
	LanguageInfos []LangInfo `json:"language_infos"`
	// 语言列表
	Languages []string `json:"languages"`
	// 前端语言列表
	FrontendLanguages []string `json:"frontend_languages"`
	// 后端语言列表
	BackendLanguages []string `json:"backend_languages"`
	// 桌面语言列表
	DesktopLanguages []string `json:"desktop_languages"`
	// 其他语言列表
	OtherLanguages []string `json:"other_languages"`
	// 主要后端语言列表 (Top 3)
	MainBackendLanguages []string `json:"main_backend_languages"`
	// 主要前端语言列表 (Top 3)
	MainFrontendLanguages []string `json:"main_frontend_languages"`
	// 框架信息列表，名称到版本的映射
	Frameworks map[string]string `json:"frameworks"`
	// 组件信息列表，名称到版本的映射
	Components map[string]string `json:"components"`
}

// FrameworkMetadata 支持列表元数据 描述了 CodeCanvas 能够识别的一种框架。
type FrameworkMetadata struct {
	Name     string            `json:"name"`
	Language string            `json:"language"`
	Levels   map[string]string `json:"levels"` // 例如: {"L1": "pom.xml", "L2": "Application.java"}
}

// ComponentMetadata 描述了 CodeCanvas 能够识别的一种组件。
type ComponentMetadata struct {
	Name     string            `json:"name"`
	Language string            `json:"language"`
	Levels   map[string]string `json:"levels"` // 例如: {"L1": "pom.xml", "L2": "Application.java"}
}

// LangSummary 保存单一语言的统计结果
type LangSummary struct {
	Name    string
	Code    int64
	Comment int64
	Blank   int64
	Count   int64
}
