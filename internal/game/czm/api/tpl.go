package api

import (
	"embed"
	"html/template"

	"github.com/gre-ory/games-go/internal/util"
)

var (
	//go:embed tpl/*.tpl
	tplFS embed.FS
)

var (
	// tpl = template.Must(
	// 	template.ParseFS(tplFS, "tpl/*.tpl"),
	// ).Funcs()

	tpl = template.Must(
		template.
			New("").
			Funcs(template.FuncMap{
				"dict": util.TplDict,
			}).
			ParseFS(tplFS, "tpl/*.tpl"),
	)
)
