package api

import (
	"net/http"

	"github.com/gre-ory/games-go/internal/game/share/model"
	"github.com/gre-ory/games-go/internal/util"
)

// //////////////////////////////////////////////////
// get user

func (s *cookieServer) htmx_get_user(w http.ResponseWriter, r *http.Request) {

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
			"User": cookie,
		}
		s.hxServer.Render(w, "user", data)
		return
	}

	// error response
	util.EncodeJsonErrorResponse(w, err)
}
