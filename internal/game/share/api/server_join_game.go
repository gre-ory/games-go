package api

import (
	"go.uber.org/zap"
)

// ////////////////////////////////////////////////
// join game

type JoinGameServer[PlayerT any] interface {
	HandleJoinGame(player PlayerT) error
}

type JoinGameService[PlayerT any, GameT any] interface {
	JoinGame(player PlayerT) (GameT, error)
}

type OnJoinGame[PlayerT any, GameT any] func(player PlayerT, game GameT)

func NewJoinGameApi[PlayerT any, GameT any](logger *zap.Logger, service JoinGameService[PlayerT, GameT], onJoinGame OnJoinGame[PlayerT, GameT]) JoinGameServer[PlayerT] {
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

func (s *joinGameServer[PlayerT, GameT]) HandleJoinGame(player PlayerT) error {
	s.logger.Info("[ws] join_game")

	game, err := s.service.JoinGame(player)
	if err != nil {
		return err
	}

	if s.onJoinGame != nil {
		s.onJoinGame(player, game)
	}

	return nil
}
