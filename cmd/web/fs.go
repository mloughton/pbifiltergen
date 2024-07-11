package embedfs

import "embed"

//go:embed assets
var StaticFiles embed.FS
