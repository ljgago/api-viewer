package assets

import (
	"embed"
	"io/fs"
)

//go:embed all:static
var static embed.FS

var StaticFS, _ = fs.Sub(static, "static")
