package api

import (
	"errors"
	"fmt"
	"net/http"

	share_model "github.com/gre-ory/games-go/internal/game/share/model"
	"github.com/gre-ory/games-go/internal/game/tictactoe/model"
	"github.com/gre-ory/games-go/internal/util/websocket"
	"go.uber.org/zap"
)

func (s *gameServer) htmx_connect(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("[api] htmx_connect ", zap.String("path", r.URL.Path))

	var cookie *share_model.Cookie
	var player *model.Player
	var game *model.Game
	var err error

	switch {
	default:

		cookie, err = s.GetValidCookie(r)
		if err != nil {
			s.logger.Info("[api] no valid cookie >>> STOP", zap.Error(err))
			break
		}
		playerId := model.PlayerId(cookie.Id)
		playerAvatar := int(cookie.Avatar)
		playerName := string(cookie.Name)
		playerLanguage := string(cookie.Language)
		s.logger.Info(fmt.Sprintf("[api] cookie %s - %s >>> connecting...", playerId, playerName), zap.Any("cookie", cookie))

		s.logger.Info("[api] s.hub.GetPlayer...")
		player, err = s.hub.GetPlayer(playerId)
		s.logger.Info("[api] ...s.hub.GetPlayer")
		if err == nil {
			s.logger.Info(fmt.Sprintf("[api] player %s already exists", playerId), zap.Any("player", player))
		} else {
			if !errors.Is(err, websocket.ErrPlayerNotFound) {
				s.logger.Info(fmt.Sprintf("[api] player %s not found >>> ERROR", playerId), zap.Error(err))
				break
			}
			s.logger.Info(fmt.Sprintf("[api] player %s not found >>> create a new one", playerId))
			// wsPlayer := websocket.NewPlayer[model.PlayerId, model.GameId](s.logger, playerId, s.onMessage, s.onPlayerUpdate, s.hub.UnregisterPlayer)
			s.logger.Info("[api] NewPlayer...")
			wsPlayer := websocket.NewPlayer[model.PlayerId, model.GameId](s.logger, playerId, s.onMessage, s.onPlayerUpdate, nil)
			player = model.NewPlayer(wsPlayer, playerAvatar, playerName, playerLanguage)
			s.logger.Info("[api] RegisterPlayer...")
			s.hub.RegisterPlayer(player)
		}
		s.logger.Info(fmt.Sprintf("[api] player %s >>> connect", playerId))
		s.logger.Info("[api] ConnectSocket...")
		player.ConnectSocket(w, r)
		s.logger.Info("[api] ...ConnectSocket")

		playerId = player.Id()

		if player.GameId() == "" {
			s.logger.Info(fmt.Sprintf("[api] player %s >>> broadcast games", playerId))
			s.broadcastJoinableGamesToPlayer(playerId)
			return
		}
		gameId := player.GameId()

		s.logger.Info("[api] GetGame...")
		game, err = s.service.GetGame(gameId)
		if err != nil {
			break
		}

		s.logger.Info(fmt.Sprintf("[api] player %s >>> broadcast game layout", playerId))
		s.broadcastGameLayoutToPlayer(playerId, game)

		s.logger.Info(fmt.Sprintf("[api] player %s >>> broadcast game", playerId))
		s.broadcastGame(game)

		return

	}

	// error response

	s.logger.Warn("[api] htmx_connect: FAILED", zap.String("path", r.URL.Path), zap.Error(err))
	s.renderError(w, err)
}
