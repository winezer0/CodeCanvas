// Package analyzer 提供代码分析相关的功能
package analyzer

// LangDefine 定义编程语言的基本特征和属性
// 用于代码分析、语法高亮、注释处理等功能

type LangDefine struct {
	Name         string     // 语言名称（如 "Go", "JavaScript"）
	LineComments []string   // 行注释标记（如 []string{"//"}）
	MultiLine    [][]string // 多行注释标记 [[Start, End], ...]（如 [][]string{{"/*", "*/"}}）
	Extensions   []string   // 文件扩展名（如 []string{".go"}）
	Filenames    []string   // 特定文件名（如 []string{"Dockerfile"}）
}

// LangDefines 预定义的编程语言列表
// 包含常见编程语言的基本特征和属性定义
var LangDefines = []LangDefine{
	{
		Name:         "Go",
		LineComments: []string{"//"},
		MultiLine:    [][]string{{"/*", "*/"}},
		Extensions:   []string{".go"},
	},
	{
		Name:         "Java",
		LineComments: []string{"//"},
		MultiLine:    [][]string{{"/*", "*/"}},
		Extensions:   []string{".java"},
	},
	{
		Name:         "C",
		LineComments: []string{"//"},
		MultiLine:    [][]string{{"/*", "*/"}},
		Extensions:   []string{".c", ".h"},
	},
	{
		Name:         "C++",
		LineComments: []string{"//"},
		MultiLine:    [][]string{{"/*", "*/"}},
		Extensions:   []string{".cpp", ".hpp", ".cc", ".cxx", ".hxx"},
	},
	{
		Name:         "C#",
		LineComments: []string{"//"},
		MultiLine:    [][]string{{"/*", "*/"}},
		Extensions:   []string{".cs"},
	},
	{
		Name:         "Python",
		LineComments: []string{"#"},
		MultiLine:    [][]string{{"\"\"\"", "\"\"\""}, {"'''", "'''"}},
		Extensions:   []string{".py"},
	},
	{
		Name:         "JavaScript",
		LineComments: []string{"//"},
		MultiLine:    [][]string{{"/*", "*/"}},
		Extensions:   []string{".js", ".mjs", ".cjs"},
	},
	{
		Name:         "TypeScript",
		LineComments: []string{"//"},
		MultiLine:    [][]string{{"/*", "*/"}},
		Extensions:   []string{".ts"},
	},
	{
		Name:         "JSX",
		LineComments: []string{"//"},
		MultiLine:    [][]string{{"/*", "*/"}},
		Extensions:   []string{".jsx"},
	},
	{
		Name:         "TSX",
		LineComments: []string{"//"},
		MultiLine:    [][]string{{"/*", "*/"}},
		Extensions:   []string{".tsx"},
	},
	{
		Name:         "Vue",
		LineComments: []string{"//"},
		MultiLine:    [][]string{{"<!--", "-->"}, {"/*", "*/"}},
		Extensions:   []string{".vue"},
	},
	{
		Name:         "HTML",
		LineComments: []string{},
		MultiLine:    [][]string{{"<!--", "-->"}},
		Extensions:   []string{".html", ".htm"},
	},
	{
		Name:         "CSS",
		LineComments: []string{},
		MultiLine:    [][]string{{"/*", "*/"}},
		Extensions:   []string{".css"},
	},
	{
		Name:         "SCSS",
		LineComments: []string{"//"},
		MultiLine:    [][]string{{"/*", "*/"}},
		Extensions:   []string{".scss"},
	},
	{
		Name:         "Less",
		LineComments: []string{"//"},
		MultiLine:    [][]string{{"/*", "*/"}},
		Extensions:   []string{".less"},
	},
	{
		Name:         "JSON",
		LineComments: []string{},
		MultiLine:    [][]string{},
		Extensions:   []string{".json"},
	},
	{
		Name:         "YAML",
		LineComments: []string{"#"},
		MultiLine:    [][]string{},
		Extensions:   []string{".yaml", ".yml"},
	},
	{
		Name:         "XML",
		LineComments: []string{},
		MultiLine:    [][]string{{"<!--", "-->"}},
		Extensions:   []string{".xml"},
	},
	{
		Name:         "PHP",
		LineComments: []string{"//", "#"},
		MultiLine:    [][]string{{"/*", "*/"}},
		Extensions:   []string{".php"},
	},
	{
		Name:         "Ruby",
		LineComments: []string{"#"},
		MultiLine:    [][]string{{"=begin", "=end"}},
		Extensions:   []string{".rb"},
	},
	{
		Name:         "Rust",
		LineComments: []string{"//"},
		MultiLine:    [][]string{{"/*", "*/"}},
		Extensions:   []string{".rs"},
	},
	{
		Name:         "Kotlin",
		LineComments: []string{"//"},
		MultiLine:    [][]string{{"/*", "*/"}},
		Extensions:   []string{".kt", ".kts"},
	},
	{
		Name:         "Swift",
		LineComments: []string{"//"},
		MultiLine:    [][]string{{"/*", "*/"}},
		Extensions:   []string{".swift"},
	},
	{
		Name:         "Shell",
		LineComments: []string{"#"},
		MultiLine:    [][]string{},
		Extensions:   []string{".sh", ".bash", ".zsh"},
	},
	{
		Name:         "PowerShell",
		LineComments: []string{"#"},
		MultiLine:    [][]string{{"<#", "#>"}},
		Extensions:   []string{".ps1"},
	},
	{
		Name:         "SQL",
		LineComments: []string{"--"},
		MultiLine:    [][]string{{"/*", "*/"}},
		Extensions:   []string{".sql"},
	},
	{
		Name:         "Lua",
		LineComments: []string{"--"},
		MultiLine:    [][]string{{"--[[", "]]"}},
		Extensions:   []string{".lua"},
	},
	{
		Name:         "Markdown",
		LineComments: []string{},
		MultiLine:    [][]string{{"<!--", "-->"}},
		Extensions:   []string{".md", ".markdown"},
	},
	{
		Name:         "Dockerfile",
		LineComments: []string{"#"},
		MultiLine:    [][]string{},
		Filenames:    []string{"Dockerfile"},
	},
	{
		Name:         "Makefile",
		LineComments: []string{"#"},
		MultiLine:    [][]string{},
		Filenames:    []string{"Makefile"},
	},
}
