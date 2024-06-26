package api

import (
	"embed"
	"html/template"
)

var (
	//go:embed tpl/*.tpl
	tplFS embed.FS
)

var (
	tpl = template.Must(
		template.ParseFS(tplFS, "tpl/*.tpl"),
	)
)
