package api

import (
	"net/http"

	"github.com/gre-ory/games-go/internal/game/share/model"
	"github.com/gre-ory/games-go/internal/util"
)

// //////////////////////////////////////////////////
// user name modal

func (s *cookieServer) htmx_user_name_modal(w http.ResponseWriter, r *http.Request) {

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

		s.hxServer.Render(w, "user-name-modal", cookie.Data())
		return
	}

	// error response
	util.EncodeJsonErrorResponse(w, err)
}
