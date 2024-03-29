package api

import (
	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/game/tictactoe/model"
)

// ////////////////////////////////////////////////
// create game

func (s *gameServer) ws_set_player_name(player *model.Player, message JsonMessage) error {
	s.logger.Info("[ws] set_player_name", zap.Any("message", message))

	playerId := player.Id()

	playerName := message.PlayerName
	if playerName == "" {
		s.logger.Info("[ws] missing player name", zap.Error(model.ErrMissingPlayerName))
		return model.ErrMissingPlayerName
	}

	player.Name = playerName

	s.logger.Info("[ws] before", zap.Any("player", player))

	s.hub.UpdatePlayer(player)

	s.logger.Info("[ws] after", zap.Any("player", player))

	other, err := s.hub.GetPlayer(playerId)
	if err == nil {
		s.logger.Info("[ws] after", zap.Any("other", other))
	}

	s.broadcastSelectGameToPlayer(playerId)

	return nil
}
