// Package embeds provides embedded rules for framework and component detection.
package embedfs

import (
	"embed"
)

//go:embed *.yml
var CanvasEmbedFS embed.FS
