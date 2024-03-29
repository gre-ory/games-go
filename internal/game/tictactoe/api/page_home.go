package api

import (
	"net/http"

	"github.com/gre-ory/games-go/internal/util"
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

	s.Render(w, "page-home", map[string]any{
		// "title":  "Tic Tac Toe",
		"title":  "TTT",
		"cookie": cookie,
	})
}
