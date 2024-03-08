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
	funcs = template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
	}

	tpl = template.Must(
		template.ParseFS(tplFS, "tpl/*.tpl"),
	)
)
