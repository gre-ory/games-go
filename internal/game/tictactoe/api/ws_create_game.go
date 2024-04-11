package api

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/game/tictactoe/model"
)

// ////////////////////////////////////////////////
// create game

func (s *gameServer) ws_create_game(player *model.Player, message JsonMessage) error {
	s.logger.Info("[ws] create_game", zap.Any("message", message))

	playerId := player.GetId()
	playerName := player.Name
	if playerName == "" {
		playerName = fmt.Sprintf("Player %s", playerId)
	}

	game, err := s.service.NewGame()
	if err != nil {
		return err
	}

	game, err = s.service.JoinGame(game.Id, player)
	if err != nil {
		return err
	}

	// player.SetGameId(game.Id)

	s.broadcastGameLayoutToPlayer(playerId, game)
	s.broadcastGame(game)
	s.broadcastJoinableGames()

	return nil
}
