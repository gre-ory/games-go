package api

import (
	"github.com/gre-ory/games-go/internal/game/tictactoe/model"
	"github.com/gre-ory/games-go/internal/util"
	"github.com/gre-ory/games-go/internal/util/websocket"
	"go.uber.org/zap"
)

type Cookie struct {
	PlayerId model.PlayerId
}

func (c *Cookie) Validate() error {
	if c.PlayerId == "" {
		return model.ErrMissingPlayerId
	}
	return nil
}

func (c *Cookie) Data() websocket.Data {
	if c.PlayerId == "" {
		return nil
	}
	return websocket.Data{
		"player_id": c.PlayerId,
	}
}

func NewCookie() *Cookie {
	return &Cookie{
		PlayerId: model.NewPlayerId(),
	}
}

const (
	CookieKey    = "game-tictactoe"
	CookieMaxAge = 60 * 60 // seconds
)

func NewCookieHelper(logger *zap.Logger, cookieSecret string) util.CookieHelper[Cookie] {
	return util.NewCookieHelper[Cookie](logger, CookieKey, CookieMaxAge, cookieSecret)
}
