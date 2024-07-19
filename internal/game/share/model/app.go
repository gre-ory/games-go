package model

import (
	"embed"
	"fmt"
	"strings"

	"github.com/gre-ory/games-go/internal/util/loc"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
	"golang.org/x/text/language"
)

// ////////////////////////////////////////////////
// add

type AppId string

func (id AppId) Loc() loc.AppId {
	return loc.AppId(id)
}

type App interface {
	Id() AppId

	Logger(mainLogger *zap.Logger) *zap.Logger

	NewDefaultEmbedBundle(fs embed.FS, langs ...language.Tag) *i18n.Bundle
	PlayerLocalizer(player Player) loc.Localizer
	UserLocalizer(user User) loc.Localizer
	Localizer(lang loc.Language) loc.Localizer

	Route(path string) string
	HomeRoute() string
	HtmxConnectRoute() string
}

func NewApp(id AppId) App {
	return &app{
		id: id,
	}
}

type app struct {
	id AppId
}

func (a *app) Id() AppId {
	return a.id
}

func (a *app) Logger(mainLogger *zap.Logger) *zap.Logger {
	return mainLogger.With(zap.String("app", string(a.id)))
}

func (a *app) NewDefaultEmbedBundle(fs embed.FS, langs ...language.Tag) *i18n.Bundle {
	return loc.NewDefaultEmbedBundle(a.Id().Loc(), fs, langs...)
}

func (a *app) PlayerLocalizer(player Player) loc.Localizer {
	return a.UserLocalizer(player.User())
}

func (a *app) UserLocalizer(user User) loc.Localizer {
	return a.Localizer(user.Language().Loc())
}

func (a *app) Localizer(lang loc.Language) loc.Localizer {
	return loc.NewLocalizer(a.id.Loc(), lang)
}

func (a *app) Route(path string) string {
	return fmt.Sprintf("/%s/%s", a.id, strings.TrimPrefix(path, "/"))
}

func (a *app) HomeRoute() string {
	return a.Route("")
}

func (a *app) HtmxConnectRoute() string {
	return a.Route("htmx/connect")
}
