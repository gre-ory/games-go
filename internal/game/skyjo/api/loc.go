package api

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/gre-ory/games-go/internal/util/loc"

	"github.com/gre-ory/games-go/internal/game/skyjo/model"
)

//go:embed loc/*.toml
var LocFS embed.FS

var bundle = loc.NewDefaultEmbedBundle(model.AppId, LocFS, language.English, language.French)
