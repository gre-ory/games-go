package api

import (
	"go.uber.org/zap"

	share_api "github.com/gre-ory/games-go/internal/game/share/api"
)

const (
	CookieKey    = "game-tictactoe"
	CookieMaxAge = 60 * 60 // seconds
)

func NewCookieServer(logger *zap.Logger, cookieSecret string, onCookie share_api.CookieCallback) share_api.CookieServer {
	return share_api.NewCookieServer(logger, CookieKey, CookieMaxAge, cookieSecret, onCookie)
}
