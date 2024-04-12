package api

import (
	"net/http"

	"github.com/gre-ory/games-go/internal/game/share/model"
	"github.com/gre-ory/games-go/internal/util"
)

// //////////////////////////////////////////////////
// user language modal

func (s *cookieServer) htmx_user_language_modal(w http.ResponseWriter, r *http.Request) {

	var cookie *model.Cookie
	var err error

	switch {
	default:

		cookie, err = s.GetValidCookie(r)
		if err != nil {
			break
		}

		err = s.SetCookie(w, cookie)
		if err != nil {
			break
		}

		s.hxServer.Render(w, "user-language-modal", cookie.Data().With("available_languages", model.GetAvailableLanguages()))
		return
	}

	// error response
	util.EncodeJsonErrorResponse(w, err)
}
