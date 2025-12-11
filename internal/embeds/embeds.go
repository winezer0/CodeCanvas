package embeds

import (
	"github.com/winezer0/codecanvas/internal/embedfs"
	"github.com/winezer0/codecanvas/internal/model"
	"io"
	"io/fs"
	"strings"

	"gopkg.in/yaml.v3"
)

// EmbeddedCanvas returns the default set of framework and component detection rules.
// Rules are embedded in the binary using the embed package and loaded from YAML files.
func EmbeddedCanvas() []*model.FrameworkRuleDefinition {
	var allRules []*model.FrameworkRuleDefinition

	// Read all files from the embedded filesystem
	files, err := fs.Glob(embedfs.CanvasEmbedFS, "*.yml")
	if err != nil {
		// Should not happen in a valid build
		return []*model.FrameworkRuleDefinition{}
	}

	for _, filename := range files {
		fileContent, err := embedfs.CanvasEmbedFS.ReadFile(filename)
		if err != nil {
			continue
		}

		// Try to parse as single-document array format first
		var rulesArray []*model.FrameworkRuleDefinition
		if err := yaml.Unmarshal(fileContent, &rulesArray); err == nil {
			// Check if we actually got something valid (array of structs)
			// yaml.Unmarshal might succeed with empty array or zero values
			if len(rulesArray) > 0 && rulesArray[0].Name != "" {
				allRules = append(allRules, rulesArray...)
				continue
			}
		}

		// Fall back to multi-document format
		yamlReader := strings.NewReader(string(fileContent))
		decoder := yaml.NewDecoder(yamlReader)

		for {
			var rule model.FrameworkRuleDefinition
			if err := decoder.Decode(&rule); err != nil {
				if err == io.EOF {
					break
				}
				// Skip malformed documents but continue with other files
				break
			}
			// Only append valid rules
			if rule.Name != "" {
				allRules = append(allRules, &rule)
			}
		}
	}

	return allRules
}
