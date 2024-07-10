package model

import (
	"strings"

	"github.com/gre-ory/games-go/internal/util/list"
	"github.com/gre-ory/games-go/internal/util/loc"
)

type UserLanguage string

var (
	UserLanguage_En UserLanguage = "en"
	UserLanguage_Fr UserLanguage = "fr"

	SupportedLanguages = []UserLanguage{
		UserLanguage_En,
		UserLanguage_Fr,
	}
)

func (l UserLanguage) Loc() loc.Language {
	return loc.Language(l)
}

func ToLanguage(value string) UserLanguage {
	value = strings.ToLower(value)
	for _, lang := range SupportedLanguages {
		if string(lang) == value {
			return lang
		}
	}
	return ""
}

func (l UserLanguage) Validate() error {
	if !list.Contains(SupportedLanguages, l) {
		return ErrUnsupportedLanguage
	}
	return nil
}

func GetAvailableLanguages() [][]UserLanguage {
	result := make([][]UserLanguage, 0, 1)
	result = append(result, SupportedLanguages)
	return result
}
