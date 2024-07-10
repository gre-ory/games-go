package websocket

import (
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/game/share/model"
	"github.com/julienschmidt/httprouter"
)

type HubServer[PlayerT Player, GameT Game[PlayerT]] interface {
	RegisterAppRoutes(router *httprouter.Router, app model.App)
	HtmxConnect(w http.ResponseWriter, r *http.Request)

	Hub() Hub[PlayerT]

	GetPlayer(id model.PlayerId) (PlayerT, error)
	RegisterPlayer(player PlayerT)
	UnregisterPlayer(id model.PlayerId)
	UpdatePlayer(player PlayerT)
	OnPlayerUpdate(playerId model.PlayerId)

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
	BroadcastPlayer(player PlayerT)
	BroadcastPlayerCookie(cookie *model.Cookie, renderUserFn func(cookie *model.Cookie) func(w io.Writer, data any))

	OnJoinGame(game GameT, player PlayerT)
	OnLeaveGame(game GameT, player PlayerT)
	OnGame(game GameT)
}

type Game[PlayerT Player] interface {
	Id() model.GameId
	Players() []PlayerT
}

type CookieServer interface {
	GetValidCookie(r *http.Request) (*model.Cookie, error)
}

func NewHubServer[PlayerT Player, GameT Game[PlayerT]](logger *zap.Logger, hub Hub[PlayerT], cookierServer CookieServer, newPlayerFromCookieFn func(cookier *model.Cookie) PlayerT, service Service[PlayerT, GameT]) HubServer[PlayerT, GameT] {
	server := &hubServer[PlayerT, GameT]{
		logger:                logger,
		hub:                   hub,
		cookierServer:         cookierServer,
		newPlayerFromCookieFn: newPlayerFromCookieFn,
		service:               service,
	}

	service.RegisterOnJoinGame(server.OnJoinGame)
	service.RegisterOnGame(server.OnGame)
	service.RegisterOnLeaveGame(server.OnLeaveGame)

	return server
}

type hubServer[PlayerT Player, GameT Game[PlayerT]] struct {
	logger                *zap.Logger
	hub                   Hub[PlayerT]
	cookierServer         CookieServer
	newPlayerFromCookieFn func(cookier *model.Cookie) PlayerT
	service               Service[PlayerT, GameT]
}

type Service[PlayerT Player, GameT Game[PlayerT]] interface {
	GetGame(gameId model.GameId) (GameT, error)
	GetJoinableGames() []GameT
	GetNonJoinableGames(playerId model.PlayerId) []GameT

	RegisterOnJoinGame(func(game GameT, player PlayerT))
	RegisterOnGame(func(game GameT))
	RegisterOnLeaveGame(func(game GameT, player PlayerT))
}

// //////////////////////////////////////////////////
// routes

func (s *hubServer[PlayerT, GameT]) RegisterAppRoutes(router *httprouter.Router, app model.App) {
	router.HandlerFunc(http.MethodGet, app.HtmxConnectRoute(), s.HtmxConnect)
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

func (s *hubServer[PlayerT, GameT]) OnPlayerUpdate(playerId model.PlayerId) {
	player, err := s.GetPlayer(playerId)
	if err != nil {
		s.logger.Error("player NOT found", zap.Any("id", playerId))
		return
	}
	s.BroadcastPlayer(player)
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
		"players": game.Players(),
	})
}

func (s *hubServer[PlayerT, GameT]) BroadcastBoard(game GameT) {
	s.hub.BroadcastToGamePlayers("board", game.Id(), Data{
		"game": game,
	})
}

func (s *hubServer[PlayerT, GameT]) BroadcastPlayer(player PlayerT) {
	s.UpdatePlayer(player)
	s.BroadcastJoinableGames()
	if player.HasGameId() {
		game, err := s.service.GetGame(player.GameId())
		if err == nil {
			s.BroadcastGame(game)
		}
	}
}

func (s *hubServer[PlayerT, GameT]) BroadcastPlayerCookie(cookie *model.Cookie, renderCookieFn func(cookie *model.Cookie) func(w io.Writer, data any)) {
	playerId := cookie.PlayerId()
	s.Hub().BroadcastToPlayerRender(playerId, nil, renderCookieFn(cookie))
}

// //////////////////////////////////////////////////
// on game events

func (s *hubServer[PlayerT, GameT]) OnJoinGame(game GameT, player PlayerT) {
	s.BroadcastGameLayoutToPlayer(player.Id(), game)
	s.OnGame(game)
}

func (s *hubServer[PlayerT, GameT]) OnLeaveGame(game GameT, player PlayerT) {
	s.BroadcastJoinableGamesToPlayer(player.Id())
	s.OnGame(game)
}

func (s *hubServer[PlayerT, GameT]) OnGame(game GameT) {
	s.BroadcastGame(game)
	s.BroadcastJoinableGames()
}
