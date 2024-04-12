package api

import (
	"go.uber.org/zap"
)

// ////////////////////////////////////////////////
// start game

type StartGameServer[PlayerT any] interface {
	HandleStartGame(player PlayerT) error
}

type StartGameService[PlayerT any, GameT any] interface {
	StartGame(player PlayerT) (GameT, error)
}

type OnStartGame[PlayerT any, GameT any] func(player PlayerT, game GameT)

func NewStartGameServer[PlayerT any, GameT any](logger *zap.Logger, service StartGameService[PlayerT, GameT], onStartGame OnStartGame[PlayerT, GameT]) StartGameServer[PlayerT] {
	return &startGameServer[PlayerT, GameT]{
		logger:      logger,
		service:     service,
		onStartGame: onStartGame,
	}
}

type startGameServer[PlayerT any, GameT any] struct {
	logger      *zap.Logger
	service     StartGameService[PlayerT, GameT]
	onStartGame OnStartGame[PlayerT, GameT]
}

func (s *startGameServer[PlayerT, GameT]) HandleStartGame(player PlayerT) error {
	s.logger.Info("[ws] start_game")

	game, err := s.service.StartGame(player)
	if err != nil {
		return err
	}

	if s.onStartGame != nil {
		s.onStartGame(player, game)
	}

	return nil
}
