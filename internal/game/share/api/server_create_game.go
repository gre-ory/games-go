package api

import (
	"go.uber.org/zap"
)

// ////////////////////////////////////////////////
// create game

type CreateGameServer[PlayerT any] interface {
	HandleCreateGame(player PlayerT) error
}

type CreateGameService[PlayerT any, GameT any] interface {
	CreateGame(player PlayerT) (GameT, error)
}

type OnCreateGame[PlayerT any, GameT any] func(player PlayerT, game GameT)

func NewCreateGameServer[PlayerT any, GameT any](logger *zap.Logger, service CreateGameService[PlayerT, GameT], onCreateGame OnCreateGame[PlayerT, GameT]) CreateGameServer[PlayerT] {
	return &createGameServer[PlayerT, GameT]{
		logger:       logger,
		service:      service,
		onCreateGame: onCreateGame,
	}
}

type createGameServer[PlayerT any, GameT any] struct {
	logger       *zap.Logger
	service      CreateGameService[PlayerT, GameT]
	onCreateGame OnCreateGame[PlayerT, GameT]
}

func (s *createGameServer[PlayerT, GameT]) HandleCreateGame(player PlayerT) error {
	s.logger.Info("[ws] create_game")

	game, err := s.service.CreateGame(player)
	if err != nil {
		return err
	}

	if s.onCreateGame != nil {
		s.onCreateGame(player, game)
	}

	return nil
}
