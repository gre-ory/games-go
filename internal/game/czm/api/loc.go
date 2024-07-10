package api

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/gre-ory/games-go/internal/game/czm/model"
)

//go:embed loc/*.toml
var LocFS embed.FS

var bundle = model.App.NewDefaultEmbedBundle(LocFS, language.English, language.French)
