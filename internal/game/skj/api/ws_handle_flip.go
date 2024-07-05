package api

import (
	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/util"

	"github.com/gre-ory/games-go/internal/game/skj/model"
)

func (s *gameServer) HandleFlip(player *model.Player, message JsonMessage) error {
	s.logger.Info("[ws] flip", zap.Any("message", message))

	columnNumber := util.ToInt(message.ColumnNumber)
	if columnNumber == 0 {
		return model.ErrInvalidColumn
	}

	rowNumber := util.ToInt(message.RowNumber)
	if rowNumber == 0 {
		return model.ErrInvalidRow
	}

	game, err := s.service.FlipCard(player, columnNumber-1, rowNumber-1)
	if err != nil {
		return err
	}

	s.logger.Info("[ws] flip", zap.Any("game", game))

	// s.broadcastClearToPlayers(game)
	s.broadcastGame(game)

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
