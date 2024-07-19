package api

import (
	"fmt"
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

		if hasUserName(r) {
			name := extractUserName(r)
			if name != "" {
				err = name.Validate()
				if err != nil {
					break
				}
				s.logger.Info(fmt.Sprintf("[DEBUG] name: %s <<< %s", cookie.Name, name))
				cookie.Name = name
			} else {
				// note: empty name >>> delete name >>> default name = id
				s.logger.Info(fmt.Sprintf("[DEBUG] name: %s <<< %s (default)", cookie.Name, model.DefaultUserName(cookie.Id)))
				cookie.Name = model.DefaultUserName(cookie.Id)
			}
		} else {
			s.logger.Info(fmt.Sprintf("[DEBUG] name: %s (untouched)", cookie.Name))
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

		err = s.SetCookie(w, cookie)
		if err != nil {
			break
		}

		s.OnCookie(cookie)

		data := model.Data{
			"User": cookie,
		}
		s.hxServer.Render(w, "user", data)
		return
	}

	// error response
	util.EncodeJsonErrorResponse(w, err)
}
