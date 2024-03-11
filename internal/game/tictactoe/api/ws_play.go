package api

import (
	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/game/tictactoe/model"
	"github.com/gre-ory/games-go/internal/util"
)

func (s *gameServer) ws_play(player *model.Player, message JsonMessage) error {
	s.logger.Info("[ws] play", zap.Any("message", message))

	playX := util.ToInt(message.PlayX)
	if playX == 0 {
		return model.ErrMissingPlayX
	}

	playY := util.ToInt(message.PlayY)
	if playY == 0 {
		return model.ErrMissingPlayY
	}

	game, err := s.service.PlayGame(player, playX, playY)
	if err != nil {
		return err
	}

	s.logger.Info("[ws] play", zap.Any("game", game))

	s.broadcastGame(game)

	return nil
}
