package api

import (
	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/game/share/model"
)

// ////////////////////////////////////////////////
// server

type GameServer[PlayerT any] interface {
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

func NewGameServer[PlayerT any, GameT any](logger *zap.Logger, service GameService[PlayerT, GameT]) GameServer[PlayerT] {
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

	game, err := s.service.CreateGame(player)
	if err != nil {
		return err
	}

	s.OnCreateGame(player, game)

	return nil
}

func (s *gameServer[PlayerT, GameT]) OnCreateGame(player PlayerT, game GameT) {
}

// //////////////////////////////////////////////////
// join game

func (s *gameServer[PlayerT, GameT]) HandleJoinGame(player PlayerT, gameId model.GameId) error {
	s.logger.Info("[ws] join_game")

	game, err := s.service.JoinGameId(gameId, player)
	if err != nil {
		return err
	}

	s.OnJoinGame(player, game)

	return nil
}

func (s *gameServer[PlayerT, GameT]) OnJoinGame(player PlayerT, game GameT) {
}

// //////////////////////////////////////////////////
// leave game

func (s *gameServer[PlayerT, GameT]) HandleLeaveGame(player PlayerT) error {
	s.logger.Info("[ws] leave_game")

	game, err := s.service.LeavePlayerGame(player)
	if err != nil {
		return err
	}

	s.OnLeaveGame(player, game)

	return nil
}

func (s *gameServer[PlayerT, GameT]) OnLeaveGame(player PlayerT, game GameT) {
}

// //////////////////////////////////////////////////
// start game

func (s *gameServer[PlayerT, GameT]) HandleStartGame(player PlayerT) error {
	s.logger.Info("[ws] start_game")

	game, err := s.service.StartPlayerGame(player)
	if err != nil {
		return err
	}

	s.OnStartGame(player, game)

	return nil
}

func (s *gameServer[PlayerT, GameT]) OnStartGame(player PlayerT, game GameT) {
}
