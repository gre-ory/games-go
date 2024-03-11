package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gre-ory/games-go/internal/game/tictactoe/model"
	"github.com/gre-ory/games-go/internal/util/websocket"
	"go.uber.org/zap"
)

func (s *gameServer) htmx_connect(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("[api] htmx_connect ", zap.String("path", r.URL.Path))

	var cookie *Cookie
	var player *model.Player
	var game *model.Game
	var err error

	switch {
	default:

		cookie, err = s.GetCookie(r)
		if err != nil {
			break
		}
		err = cookie.Validate()
		if err != nil {
			break
		}
		playerId := cookie.PlayerId

		player, err = s.hub.GetPlayer(playerId)
		if err == nil {
			s.logger.Info(fmt.Sprintf("[api] player %s already exists >>> re-use + connect", playerId), zap.Any("player", player))
		} else {
			if !errors.Is(err, websocket.ErrPlayerNotFound) {
				s.logger.Info(fmt.Sprintf("[api] player %s not found >>> ERROR", playerId), zap.Error(err))
				break
			}
			s.logger.Info(fmt.Sprintf("[api] player %s not found >>> create new one + connect", playerId))
			wsPlayer := websocket.NewPlayer[model.PlayerId, model.GameId](s.logger, playerId, s.onMessage, s.hub.UnregisterPlayer)
			player = model.NewPlayer(wsPlayer, "")
			s.hub.RegisterPlayer(player)
		}
		player.ConnectSocket(w, r)
		player.Activate()

		playerId = player.Id()

		if player.Name == "" {
			s.broadcastSelectNameToPlayer(playerId)
			return
		}

		if player.GameId() == "" {
			s.broadcastJoinableGamesToPlayer(playerId)
			return
		}
		gameId := player.GameId()

		game, err = s.service.GetGame(gameId)
		if err != nil {
			break
		}

		s.broadcastGameLayoutToPlayer(playerId, game)
		s.broadcastGame(game)

		return

	}

	// error response

	s.logger.Warn("[api] htmx_connect: FAILED", zap.String("path", r.URL.Path), zap.Error(err))
	s.renderError(w, err)
}
