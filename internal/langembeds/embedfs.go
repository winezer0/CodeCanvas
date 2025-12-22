// Package lcembeds provides embedded rules for language classification.
package langembeds

import (
	"embed"
)

//go:embed *.yml
var LanguageEmbedFS embed.FS
