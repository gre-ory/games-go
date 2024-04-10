package api

import (
	"net/http"

	"github.com/gre-ory/games-go/internal/util"
	"github.com/gre-ory/games-go/internal/util/loc"
	"go.uber.org/zap"
)

func (s *gameServer) page_home(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("[api] page_home", zap.String("path", r.URL.Path))

	//
	// cookie
	//

	cookie := s.GetCookieOrDefault(r)

	resetCookie := util.ExtractBoolParameter(r, "reset_cookie")
	if resetCookie {
		cookie = s.NewCookie()
	}

	s.SetCookie(w, cookie)

	//
	// render
	//

	localizer := loc.NewLocalizer(s.logger, string(cookie.Language))
	s.Render(w, "page-home", map[string]any{
		"cookie": cookie,
		"lang":   localizer,
	})
}
