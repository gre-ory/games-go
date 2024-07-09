package api

import (
	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/game/skj/model"
)

func (s *gameServer) HandleDiscardCard(player *model.Player) error {
	s.logger.Info("[ws] discard card", zap.Any("player", player))

	game, err := s.service.DiscardCard(player)
	if err != nil {
		return err
	}

	s.logger.Info("[ws] discard card", zap.Any("game", game))

	s.BroadcastGame(game)

	return nil
}
