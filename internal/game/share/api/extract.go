package api

import (
	"net/http"

	"github.com/gre-ory/games-go/internal/util"

	"github.com/gre-ory/games-go/internal/game/share/model"
)

func extractUserName(r *http.Request) model.UserName {
	return model.UserName(util.ExtractParameter(r, "user_name"))
}

func extractUserAvatar(r *http.Request) model.UserAvatar {
	return model.UserAvatar(util.ExtractIntParameter(r, "user_avatar"))
}

func extractUserLanguage(r *http.Request) model.UserLanguage {
	return model.UserLanguage(util.ExtractParameter(r, "user_language"))
}
