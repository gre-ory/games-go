package api

import (
	"net/http"

	"github.com/gre-ory/games-go/internal/game/share/model"
	"github.com/gre-ory/games-go/internal/util"
)

// //////////////////////////////////////////////////
// user avatar modal

func (s *cookieServer) htmx_user_avatar_modal(w http.ResponseWriter, r *http.Request) {

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
			"User":             cookie,
			"AvailableAvatars": model.GetAvailableAvatars(),
		}
		s.hxServer.Render(w, "user-avatar-modal", data)
		return
	}

	// error response
	util.EncodeJsonErrorResponse(w, err)
}
