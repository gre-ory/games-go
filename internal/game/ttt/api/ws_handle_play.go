package api

import (
	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/game/ttt/model"
)

func (s *gameServer) HandlePlay(player *model.Player, x, y int) error {
	s.logger.Info("[ws] play", zap.Int("x", x), zap.Int("y", y))

	if x == 0 {
		return model.ErrMissingPlayX
	}
	if y == 0 {
		return model.ErrMissingPlayY
	}

	game, err := s.service.PlayPlayerGame(player, x, y)
	if err != nil {
		return err
	}

	s.logger.Info("[ws] play", zap.Any("game", game))

	// s.broadcastClearToPlayers(game)
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
