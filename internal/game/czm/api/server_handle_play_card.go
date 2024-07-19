package api

import (
	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/game/czm/model"
)

func (s *gameServer) HandlePlayCard(player *model.Player, discardNumber int) error {
	s.logger.Info("[ws] play_card", zap.Int("discardNumber", discardNumber))

	if discardNumber == 0 {
		return model.ErrInvalidDiscardNumber
	}

	game, err := s.service.PlayCard(player, discardNumber)
	if err != nil {
		return err
	}

	s.logger.Info("[ws] play", zap.Any("game", game))

	s.BroadcastGame(game)

	// if game.Stopped() {
	// 	if yes, winnerId := game.HasWinner(); yes {
	// 		winner, err := game.GetPlayer(winnerId)
	// 		if err == nil {
	// 			s.broadcastInfoToPlayers(
	// 				game,
	// 				fmt.Sprintf("%s wins!", winner.Name),
	// 			)
	// 		}
	// 	} else if game.IsTie() {
	// 		s.broadcastInfoToPlayers(
	// 			game,
	// 			"It is a tie!",
	// 		)
	// 	}
	// }

	return nil
}
