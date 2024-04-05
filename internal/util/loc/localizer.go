package loc

import (
	"fmt"
	"html/template"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
	"golang.org/x/text/language"
)

// //////////////////////////////////////////////////
// global

var (
	defaultLanguage  *language.Tag
	loadedLocalizers = map[string]*i18n.Localizer{}
)

// //////////////////////////////////////////////////
// localizer

type Localizer interface {
	LocalizeText(id string) template.HTML
	Loc(id string, args ...any) template.HTML
	Localize(id string, data any) template.HTML
}

// //////////////////////////////////////////////////
// constructor

func NewLocalizer(logger *zap.Logger, lang string) Localizer {
	localizers := make([]*i18n.Localizer, 2)
	if lang != "" {
		if localizer, ok := loadedLocalizers[lang]; ok {
			localizers = append(localizers, localizer)
		}
	}
	if defaultLanguage != nil && (lang == "" || defaultLanguage.String() != lang) {
		if defaultLocalizer, ok := loadedLocalizers[defaultLanguage.String()]; ok {
			localizers = append(localizers, defaultLocalizer)
		}
	}
	return &localizer{
		logger:     logger,
		localizers: localizers,
	}
}

// //////////////////////////////////////////////////
// localizer

type localizer struct {
	logger     *zap.Logger
	localizers []*i18n.Localizer
}

func (l *localizer) LocalizeText(id string) template.HTML {
	return l.Localize(id, nil)
}

func (l *localizer) Loc(id string, args ...any) template.HTML {
	data := map[string]any{}
	for i, arg := range args {
		data[fmt.Sprintf("arg%d", i+1)] = arg
	}
	// l.logger.Info(fmt.Sprintf("[loc] %s: %#v \n", id, data))
	return l.Localize(id, data)
}

func (l *localizer) Localize(id string, data any) template.HTML {
	if len(l.localizers) == 0 {
		return l.error("missing localizers", id)
	}
	cfg := &i18n.LocalizeConfig{
		MessageID:    id,
		TemplateData: data,
	}
	var firstErr error
	for _, localizer := range l.localizers {
		if localizer == nil {
			continue
		}
		localized, err := localizer.Localize(cfg)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		if firstErr != nil {
			return l.warning(localized)
		}
		return template.HTML(localized)
	}
	if firstErr != nil {
		return l.error(firstErr.Error())
	}
	return l.error("missing %q", id)
}

func (l *localizer) warning(msg string, args ...any) template.HTML {
	return template.HTML(fmt.Sprintf("<span style=\"color: #ff9b00; background: rgba(255, 143, 0, 0.2);\">"+msg+"</span>", args...))
}

func (l *localizer) error(msg string, args ...any) template.HTML {
	return template.HTML(fmt.Sprintf("<span style=\"color: red;background: rgba(255,0,0,0.2);\">"+msg+"</span>", args...))
}
