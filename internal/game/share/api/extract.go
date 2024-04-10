package api

import (
	"net/http"

	"github.com/gre-ory/games-go/internal/util"

	"github.com/gre-ory/games-go/internal/game/share/model"
)

const (
	UserNameParameter     = "user_name"
	UserAvatarParameter   = "user_avatar"
	UserLanguageParameter = "user_language"
)

func hasUserName(r *http.Request) bool {
	return util.HasParameter(r, UserNameParameter)
}

func extractUserName(r *http.Request) model.UserName {
	return model.UserName(util.ExtractParameter(r, UserNameParameter))
}

func extractUserAvatar(r *http.Request) model.UserAvatar {
	return model.UserAvatar(util.ExtractIntParameter(r, UserAvatarParameter))
}

func extractUserLanguage(r *http.Request) model.UserLanguage {
	return model.ToLanguage(util.ExtractParameter(r, UserLanguageParameter))
}
