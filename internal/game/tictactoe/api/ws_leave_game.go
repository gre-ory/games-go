package api

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/game/tictactoe/model"
)

// ////////////////////////////////////////////////
// join game

func (s *gameServer) ws_leave_game(player *model.Player, message JsonMessage) error {
	s.logger.Info("[ws] leave_game", zap.Any("message", message))

	playerId := player.GetId()
	playerName := player.Name
	if playerName == "" {
		playerName = fmt.Sprintf("Player %s", playerId)
	}

	game, err := s.service.LeaveGame(player)
	if err != nil {
		return err
	}

	// player.SetGameId(game.Id)

	s.broadcastJoinableGamesToPlayer(playerId)
	if game != nil {
		s.broadcastGame(game)
	}
	s.broadcastJoinableGames()

	return nil
}
