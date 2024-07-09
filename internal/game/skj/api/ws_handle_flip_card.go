package api

import (
	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/game/skj/model"
)

func (s *gameServer) HandleFlipCard(player *model.Player, message JsonMessage) error {
	s.logger.Info("[ws] flip card", zap.Any("message", message))

	columnNumber, rowNumber, err := message.Cell()
	if err != nil {
		return err
	}

	game, err := s.service.FlipCard(player, columnNumber, rowNumber)
	if err != nil {
		return err
	}

	s.logger.Info("[ws] flip card", zap.Any("game", game))

	s.BroadcastGame(game)

	return nil
}
