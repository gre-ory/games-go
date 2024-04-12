package api

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/gre-ory/games-go/internal/game/tictactoe/model"
	"github.com/gre-ory/games-go/internal/util/loc"
)

//go:embed loc/*.toml
var LocFS embed.FS

var bundle = loc.NewDefaultEmbedBundle(model.AppId, LocFS, language.English, language.French)
