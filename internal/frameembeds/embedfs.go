// Package embeds provides embedded rules for framework and component detection.
package frameembeds

import (
	"embed"
)

//go:embed *.yml
var FrameEmbedFS embed.FS
