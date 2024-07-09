package api

import (
	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/game/skj/model"
)

func (s *gameServer) HandleDrawCard(player *model.Player) error {
	s.logger.Info("[ws] draw card", zap.Any("player", player))

	game, err := s.service.DrawCard(player)
	if err != nil {
		return err
	}

	s.logger.Info("[ws] draw card", zap.Any("game", game))

	s.BroadcastGame(game)

	return nil
}
