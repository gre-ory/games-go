package api

import (
	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/game/skj/model"
)

func (s *gameServer) HandlePutCard(player *model.Player, columnNumber, rowNumber int) error {
	s.logger.Info("[ws] put card", zap.Int("column", columnNumber), zap.Int("row", rowNumber))

	if columnNumber == 0 {
		return model.ErrInvalidColumn
	}
	if rowNumber == 0 {
		return model.ErrInvalidRow
	}

	game, err := s.service.PutCard(player, columnNumber, rowNumber)
	if err != nil {
		return err
	}

	s.logger.Info("[ws] put card", zap.Any("game", game))

	s.BroadcastGame(game)

	return nil
}
