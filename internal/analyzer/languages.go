package analyzer

import "strings"

type LanguageDefinition struct {
	Name         string
	LineComments []string
	MultiLine    [][]string // [[Start, End], [Start, End]]
	Extensions   []string
	Filenames    []string
}

var languages = []LanguageDefinition{
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

var (
	extToLanguage  = make(map[string]*LanguageDefinition)
	fileToLanguage = make(map[string]*LanguageDefinition)
)

func init() {
	for i := range languages {
		lang := &languages[i]
		for _, ext := range lang.Extensions {
			extToLanguage[ext] = lang
		}
		for _, name := range lang.Filenames {
			fileToLanguage[name] = lang
		}
	}
}

// GetLanguageByExtension returns the language definition for a given file extension
func GetLanguageByExtension(ext string) *LanguageDefinition {
	return extToLanguage[strings.ToLower(ext)]
}

// GetLanguageByFilename returns the language definition for a given filename
func GetLanguageByFilename(name string) *LanguageDefinition {
	return fileToLanguage[name]
}
