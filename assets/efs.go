package assets

import "embed"

//go:embed "templates" "docs"
var EmbededFiles embed.FS
