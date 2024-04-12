package api

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/util"
	"github.com/gre-ory/games-go/internal/util/loc"

	share_model "github.com/gre-ory/games-go/internal/game/share/model"

	"github.com/gre-ory/games-go/internal/game/skyjo/model"
)

func (s *gameServer) page_home(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("[api] page_home", zap.String("path", r.URL.Path))

	var cookie *share_model.Cookie
	var err error

	switch {
	default:

		//
		// cookie
		//

		cookie = s.GetCookieOrDefault(r)

		resetCookie := util.ExtractBoolParameter(r, "reset_cookie")
		if resetCookie {
			cookie = s.NewCookie()
		}

		err = s.SetCookie(w, cookie)
		if err != nil {
			break
		}

		//
		// render
		//

		localizer := loc.NewLocalizer(model.AppId, string(cookie.Language), s.logger)
		s.Render(w, "page-home", map[string]any{
			"cookie": cookie,
			"lang":   localizer,
		})
		return
	}

	// error response
	util.EncodeJsonErrorResponse(w, err)
}
