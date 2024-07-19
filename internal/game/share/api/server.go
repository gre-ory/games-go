package api

import (
	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/util/dict"

	"github.com/gre-ory/games-go/internal/game/share/model"
)

// ////////////////////////////////////////////////
// server

type GameServer[PlayerT model.Player, GameT model.Game[PlayerT]] interface {
	HandleCreateGame(user model.User) error
	HandleJoinGame(gameId model.GameId, user model.User) error
	HandleStartGame(player PlayerT) error
	HandleLeaveGame(player PlayerT) error
}

type GameService[PlayerT model.Player, GameT model.Game[PlayerT]] interface {
	CreateGame(user model.User) (GameT, error)
	JoinGameId(gameId model.GameId, user model.User) (GameT, error)
	StartPlayerGame(player PlayerT) (GameT, error)
	LeavePlayerGame(player PlayerT) (GameT, error)
}

func NewGameServer[PlayerT model.Player, GameT model.Game[PlayerT]](logger *zap.Logger, service GameService[PlayerT, GameT]) GameServer[PlayerT, GameT] {
	return &gameServer[PlayerT, GameT]{
		logger:  logger,
		service: service,
	}
}

type gameServer[PlayerT model.Player, GameT model.Game[PlayerT]] struct {
	logger  *zap.Logger
	service GameService[PlayerT, GameT]
	players map[model.PlayerId]PlayerT
}

// //////////////////////////////////////////////////
// players

func (s *gameServer[PlayerT, GameT]) GetPlayer(playerId model.PlayerId) (PlayerT, bool) {
	if playerId == "" {
		var empty PlayerT
		return empty, false
	}
	player, ok := s.players[playerId]
	return player, ok
}

func (s *gameServer[PlayerT, GameT]) GetPlayers() []PlayerT {
	return dict.Values(s.players)
}

func (s *gameServer[PlayerT, GameT]) RegisterPlayer(player PlayerT) {
	s.players[player.Id()] = player
}

func (s *gameServer[PlayerT, GameT]) UnregisterPlayerId(playerId model.PlayerId) {
	delete(s.players, playerId)
}

// //////////////////////////////////////////////////
// create game

func (s *gameServer[PlayerT, GameT]) HandleCreateGame(user model.User) error {
	s.logger.Info("[ws] create_game")
	_, err := s.service.CreateGame(user)
	return err
}

// //////////////////////////////////////////////////
// join game

func (s *gameServer[PlayerT, GameT]) HandleJoinGame(gameId model.GameId, user model.User) error {
	s.logger.Info("[ws] join_game")
	if gameId == "" {
		return model.ErrMissingGameId
	}
	_, err := s.service.JoinGameId(gameId, user)
	return err
}

// //////////////////////////////////////////////////
// start game

func (s *gameServer[PlayerT, GameT]) HandleStartGame(player PlayerT) error {
	s.logger.Info("[ws] start_game")
	_, err := s.service.StartPlayerGame(player)
	return err
}

// //////////////////////////////////////////////////
// leave game

func (s *gameServer[PlayerT, GameT]) HandleLeaveGame(player PlayerT) error {
	s.logger.Info("[ws] leave_game")
	_, err := s.service.LeavePlayerGame(player)
	return err
}
