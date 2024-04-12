package loc

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// //////////////////////////////////////////////////
// global

var (
	apps = map[string]App{}
)

func registerApp(app App) {
	apps[app.GetId()] = app
}

func GetApp(id string) App {
	if app, ok := apps[id]; ok {
		return app
	}
	app := NewApp(id)
	return app
}

// //////////////////////////////////////////////////
// global

type App interface {
	GetId() string
	SetDefaultLanguage(lang string)
	AddLocalizer(lang string, localizer *i18n.Localizer)
	GetLocalizer(lang string) *i18n.Localizer
	GetDefaultLocalizer() *i18n.Localizer
}

func NewApp(id string) App {
	app := &app{
		ID:         id,
		localizers: map[string]*i18n.Localizer{},
	}
	registerApp(app)
	return app
}

type app struct {
	ID              string
	defaultLanguage *string
	localizers      map[string]*i18n.Localizer
}

func (a *app) GetId() string {
	return a.ID
}

func (a *app) SetDefaultLanguage(lang string) {
	a.defaultLanguage = &lang
}

func (a *app) AddLocalizer(lang string, localizer *i18n.Localizer) {
	a.localizers[lang] = localizer
}

func (a *app) GetLocalizer(lang string) *i18n.Localizer {
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
