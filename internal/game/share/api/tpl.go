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
	ShareTpl = template.Must(
		template.ParseFS(tplFS, "tpl/*.tpl"),
	)
)
