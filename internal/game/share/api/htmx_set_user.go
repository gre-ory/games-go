package api

import (
	"net/http"

	"github.com/gre-ory/games-go/internal/game/share/model"
	"github.com/gre-ory/games-go/internal/util"
)

// //////////////////////////////////////////////////
// set user

func (s *cookieServer) htmx_set_user(w http.ResponseWriter, r *http.Request) {

	var cookie *model.Cookie
	var err error

	switch {
	default:

		cookie, err = s.GetValidCookie(r)
		if err != nil {
			break
		}

		name := extractUserName(r)
		if name != "" {
			err = name.Validate()
			if err != nil {
				break
			}
			cookie.Name = name
		}

		avatar := extractUserAvatar(r)
		if avatar != 0 {
			err = avatar.Validate()
			if err != nil {
				break
			}
			cookie.Avatar = avatar
		}

		language := extractUserLanguage(r)
		if language != "" {
			err = language.Validate()
			if err != nil {
				break
			}
			cookie.Language = language
		}

		s.onCookie(cookie)

		s.SetCookie(w, cookie)
		s.hxServer.Render(w, "user", cookie.Data())
		return
	}

	// error response
	util.EncodeJsonErrorResponse(w, err)
}
