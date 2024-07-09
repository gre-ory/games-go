package websocket

import (
	"github.com/gre-ory/games-go/internal/game/share/model"
)

type HubServer[PlayerT Player, GameT Game[PlayerT]] interface {
	Hub() Hub[PlayerT]
	GetPlayer(id model.PlayerId) (PlayerT, error)
	RegisterPlayer(player PlayerT)
	UnregisterPlayer(id model.PlayerId)
	UpdatePlayer(player PlayerT)
	BroadcastInfoToPlayers(game GameT, info string)
	BroadcastInfoToPlayer(playerId model.PlayerId, info string)
	BroadcastErrorToPlayer(playerId model.PlayerId, err error)
	BroadcastSelectGameToPlayer(playerId model.PlayerId)
	BroadcastGameLayoutToPlayer(playerId model.PlayerId, game GameT)
	BroadcastJoinableGamesToPlayer(playerId model.PlayerId)
	BroadcastJoinableGames()
	BroadcastGame(game GameT)
	BroadcastPlayers(game GameT)
	BroadcastBoard(game GameT)
}

type Game[PlayerT Player] interface {
	Id() model.GameId
	GetPlayers() []PlayerT
}

func NewHubServer[PlayerT Player, GameT Game[PlayerT]](hub Hub[PlayerT], service Service[PlayerT, GameT]) HubServer[PlayerT, GameT] {
	return &hubServer[PlayerT, GameT]{
		hub:     hub,
		service: service,
	}
}

type hubServer[PlayerT Player, GameT Game[PlayerT]] struct {
	hub     Hub[PlayerT]
	service Service[PlayerT, GameT]
}

type Service[PlayerT Player, GameT Game[PlayerT]] interface {
	GetJoinableGames() []GameT
	GetNonJoinableGames(playerId model.PlayerId) []GameT
}

// //////////////////////////////////////////////////
// hub

func (s *hubServer[PlayerT, GameT]) Hub() Hub[PlayerT] {
	return s.hub
}

func (s *hubServer[PlayerT, GameT]) GetPlayer(id model.PlayerId) (PlayerT, error) {
	return s.hub.GetPlayer(id)
}

func (s *hubServer[PlayerT, GameT]) RegisterPlayer(player PlayerT) {
	s.hub.RegisterPlayer(player)
}

func (s *hubServer[PlayerT, GameT]) UnregisterPlayer(id model.PlayerId) {
	s.hub.UnregisterPlayer(id)
}

func (s *hubServer[PlayerT, GameT]) UpdatePlayer(player PlayerT) {
	s.hub.UpdatePlayer(player)
}

// //////////////////////////////////////////////////
// broadcast

func (s *hubServer[PlayerT, GameT]) BroadcastInfoToPlayers(game GameT, info string) {
	s.hub.BroadcastToGamePlayers("info", game.Id(), Data{
		"info": info,
	})
}

func (s *hubServer[PlayerT, GameT]) BroadcastInfoToPlayer(playerId model.PlayerId, info string) {
	s.hub.BroadcastToPlayer("info", playerId, Data{
		"info": info,
	})
}

func (s *hubServer[PlayerT, GameT]) BroadcastErrorToPlayer(playerId model.PlayerId, err error) {
	s.hub.BroadcastToPlayer("error", playerId, Data{
		"error": err.Error(),
	})
}

func (s *hubServer[PlayerT, GameT]) BroadcastSelectGameToPlayer(playerId model.PlayerId) {
	data := s.getJoinableGamesData(playerId)
	s.hub.BroadcastToPlayer("select-game", playerId, data)
}

func (s *hubServer[PlayerT, GameT]) BroadcastGameLayoutToPlayer(playerId model.PlayerId, game GameT) {
	s.hub.BroadcastToPlayer("game-layout", playerId, Data{
		"game": game,
	})
}

func (s *hubServer[PlayerT, GameT]) BroadcastJoinableGamesToPlayer(playerId model.PlayerId) {
	data := s.getJoinableGamesData(playerId)
	s.hub.BroadcastToPlayer("select-game", playerId, data)
}

func (s *hubServer[PlayerT, GameT]) BroadcastJoinableGames() {
	s.hub.BroadcastToNotPlayingPlayersFn("select-game", func(player PlayerT) (bool, any) {
		data := s.getJoinableGamesData(player.Id())
		return s.hub.WrapPlayerData(data, player)
	})
}

func (s *hubServer[PlayerT, GameT]) getJoinableGamesData(playerId model.PlayerId) Data {
	waitingPlayers := s.getWaitingPlayers(playerId)
	data := make(Data)
	data["new_games"] = s.service.GetJoinableGames()
	data["other_games"] = s.service.GetNonJoinableGames(playerId)
	data["has_waiting_players"] = len(waitingPlayers) > 0
	data["waiting_players"] = waitingPlayers
	return data
}

func (s *hubServer[PlayerT, GameT]) getWaitingPlayers(playerId model.PlayerId) []PlayerT {
	players := s.hub.GetNotPlayingPlayers()
	waitingPlayers := make([]PlayerT, 0, len(players))
	for _, player := range players {
		if player.Id() == playerId {
			continue
		}
		waitingPlayers = append(waitingPlayers, player)
	}
	return waitingPlayers
}

func (s *hubServer[PlayerT, GameT]) BroadcastGame(game GameT) {
	s.BroadcastPlayers(game)
	s.BroadcastBoard(game)
}

func (s *hubServer[PlayerT, GameT]) BroadcastPlayers(game GameT) {
	s.hub.BroadcastToGamePlayers("players", game.Id(), Data{
		"players": game.GetPlayers(),
	})
}

func (s *hubServer[PlayerT, GameT]) BroadcastBoard(game GameT) {
	s.hub.BroadcastToGamePlayers("board", game.Id(), Data{
		"game": game,
	})
}
