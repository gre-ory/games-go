package loc

import (
	"embed"
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/gre-ory/games-go/internal/util/list"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func NewDefaultEmbedBundle(appId AppId, fs embed.FS, langs ...language.Tag) *i18n.Bundle {
	if len(langs) == 0 {
		panic("[loc] missing languages!")
	}

	return NewEmbedBundle(appId, langs[0], fs, list.Convert(langs, defaultLocPath)...)
}

func defaultLocPath(lang language.Tag) string {
	return fmt.Sprintf("loc/%s.toml", lang.String())
}

func NewEmbedBundle(appId AppId, defaultLang language.Tag, fs embed.FS, paths ...string) *i18n.Bundle {
	if len(paths) == 0 {
		panic("[loc] missing language files!")
	}

	app := GetApp(appId)

	bundle := i18n.NewBundle(defaultLang)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	for _, path := range paths {
		file, err := bundle.LoadMessageFileFS(fs, path)
		if err != nil {
			fmt.Printf("[loc] %s: Error while loading >>> %s\n", path, err.Error())
			continue
		}
		lang := file.Tag.String()
		fmt.Printf("[loc] %s: (+) lang: %s, path: %s \n", appId, lang, path)

		app.AddLocalizer(Language(lang), i18n.NewLocalizer(bundle, lang))
		if lang == defaultLang.String() {
			app.SetDefaultLanguage(Language(defaultLang.String()))
		}
	}

	fmt.Printf("[loc] languages: %s \n", joinLanguages(bundle.LanguageTags()))

	return bundle
}

func joinLanguages(langs []language.Tag) string {
	return strings.Join(list.Convert(langs, language.Tag.String), ",")
}
