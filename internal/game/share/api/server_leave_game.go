package api

import (
	"go.uber.org/zap"
)

// ////////////////////////////////////////////////
// leave game

type LeaveGameServer[PlayerT any] interface {
	HandleLeaveGame(player PlayerT) error
}

type LeaveGameService[PlayerT any, GameT any] interface {
	LeaveGame(player PlayerT) (GameT, error)
}

type OnLeaveGame[PlayerT any, GameT any] func(player PlayerT, game GameT)

func NewLeaveGameServer[PlayerT any, GameT any](logger *zap.Logger, service LeaveGameService[PlayerT, GameT], onLeaveGame OnLeaveGame[PlayerT, GameT]) LeaveGameServer[PlayerT] {
	return &leaveGameServer[PlayerT, GameT]{
		logger:      logger,
		service:     service,
		onLeaveGame: onLeaveGame,
	}
}

type leaveGameServer[PlayerT any, GameT any] struct {
	logger      *zap.Logger
	service     LeaveGameService[PlayerT, GameT]
	onLeaveGame OnLeaveGame[PlayerT, GameT]
}

func (s *leaveGameServer[PlayerT, GameT]) HandleLeaveGame(player PlayerT) error {
	s.logger.Info("[ws] leave_game")

	game, err := s.service.LeaveGame(player)
	if err != nil {
		return err
	}

	if s.onLeaveGame != nil {
		s.onLeaveGame(player, game)
	}

	return nil
}
