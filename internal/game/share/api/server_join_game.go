package api

import (
	"go.uber.org/zap"
)

// ////////////////////////////////////////////////
// join game

type JoinGameServer[PlayerT any, GameIdT comparable] interface {
	HandleJoinGame(player PlayerT, gameId GameIdT) error
}

type JoinGameService[GameIdT comparable, PlayerT any, GameT any] interface {
	JoinGame(gameId GameIdT, player PlayerT) (GameT, error)
}

type OnJoinGame[PlayerT any, GameT any] func(player PlayerT, game GameT)

func NewJoinGameServer[PlayerT any, GameIdT comparable, GameT any](
	logger *zap.Logger,
	service JoinGameService[GameIdT, PlayerT, GameT],
	onJoinGame OnJoinGame[PlayerT, GameT],
) JoinGameServer[PlayerT, GameIdT] {
	return &joinGameServer[PlayerT, GameIdT, GameT]{
		logger:     logger,
		service:    service,
		onJoinGame: onJoinGame,
	}
}

type joinGameServer[PlayerT any, GameIdT comparable, GameT any] struct {
	logger     *zap.Logger
	service    JoinGameService[GameIdT, PlayerT, GameT]
	onJoinGame OnJoinGame[PlayerT, GameT]
}

func (s *joinGameServer[PlayerT, GameIdT, GameT]) HandleJoinGame(player PlayerT, gameId GameIdT) error {
	s.logger.Info("[ws] join_game")

	game, err := s.service.JoinGame(gameId, player)
	if err != nil {
		return err
	}

	if s.onJoinGame != nil {
		s.onJoinGame(player, game)
	}

	return nil
}
