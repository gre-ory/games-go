package api

import (
	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/game/czm/model"
)

func (s *gameServer) HandleSelectCard(player *model.Player, cardNumber int) error {
	s.logger.Info("[ws] select_card", zap.Int("card", cardNumber))

	if cardNumber == 0 {
		return model.ErrInvalidCardNumber
	}

	game, err := s.service.SelectCard(player, cardNumber)
	if err != nil {
		return err
	}

	s.logger.Info("[ws] play", zap.Any("game", game))

	s.BroadcastGame(game)

	return nil
}
