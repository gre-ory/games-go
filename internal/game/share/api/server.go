package api

import (
	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/game/share/model"
)

// ////////////////////////////////////////////////
// server

type GameServer[PlayerT any, GameT any] interface {
	HandleCreateGame(player PlayerT) error
	HandleJoinGame(player PlayerT, gameId model.GameId) error
	HandleStartGame(player PlayerT) error
	HandleLeaveGame(player PlayerT) error
}

type GameService[PlayerT any, GameT any] interface {
	CreateGame(player PlayerT) (GameT, error)
	JoinGameId(gameId model.GameId, player PlayerT) (GameT, error)
	StartPlayerGame(player PlayerT) (GameT, error)
	LeavePlayerGame(player PlayerT) (GameT, error)
}

func NewGameServer[PlayerT any, GameT any](logger *zap.Logger, service GameService[PlayerT, GameT]) GameServer[PlayerT, GameT] {
	return &gameServer[PlayerT, GameT]{
		logger:  logger,
		service: service,
	}
}

type gameServer[PlayerT any, GameT any] struct {
	logger  *zap.Logger
	service GameService[PlayerT, GameT]
}

// //////////////////////////////////////////////////
// create game

func (s *gameServer[PlayerT, GameT]) HandleCreateGame(player PlayerT) error {
	s.logger.Info("[ws] create_game")
	_, err := s.service.CreateGame(player)
	return err
}

// //////////////////////////////////////////////////
// join game

func (s *gameServer[PlayerT, GameT]) HandleJoinGame(player PlayerT, gameId model.GameId) error {
	s.logger.Info("[ws] join_game")
	_, err := s.service.JoinGameId(gameId, player)
	return err
}

// //////////////////////////////////////////////////
// leave game

func (s *gameServer[PlayerT, GameT]) HandleLeaveGame(player PlayerT) error {
	s.logger.Info("[ws] leave_game")
	_, err := s.service.LeavePlayerGame(player)
	return err
}

// //////////////////////////////////////////////////
// start game

func (s *gameServer[PlayerT, GameT]) HandleStartGame(player PlayerT) error {
	s.logger.Info("[ws] start_game")
	_, err := s.service.StartPlayerGame(player)
	return err
}
