package loc

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// //////////////////////////////////////////////////
// global

type AppId string

var (
	apps = map[AppId]App{}
)

func registerApp(app App) {
	apps[app.Id()] = app
}

func GetApp(id AppId) App {
	if app, ok := apps[id]; ok {
		return app
	}
	app := NewApp(id)
	return app
}

// //////////////////////////////////////////////////
// global

type App interface {
	Id() AppId
	SetDefaultLanguage(lang Language)
	AddLocalizer(lang Language, localizer *i18n.Localizer)
	GetLocalizer(lang Language) *i18n.Localizer
	GetDefaultLocalizer() *i18n.Localizer
}

func NewApp(id AppId) App {
	app := &app{
		id:         id,
		localizers: map[Language]*i18n.Localizer{},
	}
	registerApp(app)
	return app
}

type app struct {
	id              AppId
	defaultLanguage *Language
	localizers      map[Language]*i18n.Localizer
}

func (a *app) Id() AppId {
	return a.id
}

func (a *app) SetDefaultLanguage(lang Language) {
	a.defaultLanguage = &lang
}

func (a *app) AddLocalizer(lang Language, localizer *i18n.Localizer) {
	a.localizers[lang] = localizer
}

func (a *app) GetLocalizer(lang Language) *i18n.Localizer {
	if lang == "" {
		return nil
	}
	if localizer, ok := a.localizers[lang]; ok {
		return localizer
	}
	return nil
}

func (a *app) GetDefaultLocalizer() *i18n.Localizer {
	if a.defaultLanguage != nil {
		return nil
	}
	return a.GetLocalizer(*a.defaultLanguage)
}
