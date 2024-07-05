package api

import (
	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/game/share/model"
)

// ////////////////////////////////////////////////
// join game

type JoinGameServer[PlayerT any] interface {
	HandleJoinGame(player PlayerT, gameId model.GameId) error
}

type JoinGameService[PlayerT any, GameT any] interface {
	JoinGameId(gameId model.GameId, player PlayerT) (GameT, error)
}

type OnJoinGame[PlayerT any, GameT any] func(player PlayerT, game GameT)

func NewJoinGameServer[PlayerT any, GameT any](
	logger *zap.Logger,
	service JoinGameService[PlayerT, GameT],
	onJoinGame OnJoinGame[PlayerT, GameT],
) JoinGameServer[PlayerT] {
	return &joinGameServer[PlayerT, GameT]{
		logger:     logger,
		service:    service,
		onJoinGame: onJoinGame,
	}
}

type joinGameServer[PlayerT any, GameT any] struct {
	logger     *zap.Logger
	service    JoinGameService[PlayerT, GameT]
	onJoinGame OnJoinGame[PlayerT, GameT]
}

func (s *joinGameServer[PlayerT, GameT]) HandleJoinGame(player PlayerT, gameId model.GameId) error {
	s.logger.Info("[ws] join_game")

	game, err := s.service.JoinGameId(gameId, player)
	if err != nil {
		return err
	}

	if s.onJoinGame != nil {
		s.onJoinGame(player, game)
	}

	return nil
}
