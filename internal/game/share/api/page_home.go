package api

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/util"

	"github.com/gre-ory/games-go/internal/game/share/model"
)

func PageHome(logger *zap.Logger, app model.App, cookieServer CookieServer, hxServer util.HxServer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("[api] page_home", zap.String("path", r.URL.Path))

		var cookie *model.Cookie
		var err error

		switch {
		default:

			//
			// cookie
			//

			cookie = cookieServer.GetCookieOrDefault(r)

			resetCookie := util.ExtractBoolParameter(r, "reset_cookie")
			if resetCookie {
				cookie = cookieServer.NewCookie()
			}

			err = cookieServer.SetCookie(w, cookie)
			if err != nil {
				break
			}

			//
			// render
			//

			renderer := NewRenderer()
			hxServer.Render(w, "page-home", map[string]any{
				"AppId":      app.Id(),
				"Cookie":     cookie,
				"Lang":       app.Localizer(cookie.Language.Loc()),
				"ConnectUrl": app.HtmxConnectRoute(),
				"Share":      renderer,
			})
			return
		}

		// error response
		util.EncodeJsonErrorResponse(w, err)
	}
}
