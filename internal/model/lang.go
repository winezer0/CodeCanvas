package model

// LangFeatures 定义语言的语法特征，用于辅助语言分类
// - Name: 语法特征名称
// - Tokens: 语言特有的语法标记
// - FilePatterns: 语言文件的路径模式
// - Dependencies: 语言常用的依赖包名称

type LangFeatures struct {
	Name         string   `json:"name"`
	Tokens       []string `json:"tokens"`
	FilePatterns []string `json:"file_patterns"`
	Dependencies []string `json:"dependencies"`
}

// LangRule 定义语言分类规则
// - Name: 语言名称
// - Category: 语言分类（前端/后端/桌面）
// - Features: 语言的语法特征

type LangRule struct {
	Name     string       `json:"name"`
	Category string       `json:"category"`
	Features LangFeatures `json:"features"`
}
