package api

import (
	"net/http"

	"github.com/gre-ory/games-go/internal/util"

	"github.com/gre-ory/games-go/internal/game/share/model"
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

		data := model.Data{
			"User":               cookie,
			"AvailableLanguages": model.GetAvailableLanguages(),
		}
		s.hxServer.Render(w, "user-language-modal", data)
		return
	}

	// error response
	util.EncodeJsonErrorResponse(w, err)
}
