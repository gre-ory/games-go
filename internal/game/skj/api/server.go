package api

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/util"

	share_api "github.com/gre-ory/games-go/internal/game/share/api"
	share_model "github.com/gre-ory/games-go/internal/game/share/model"
	"github.com/gre-ory/games-go/internal/game/share/websocket"
	share_websocket "github.com/gre-ory/games-go/internal/game/share/websocket"

	"github.com/gre-ory/games-go/internal/game/skj/model"
	"github.com/gre-ory/games-go/internal/game/skj/service"
)

// //////////////////////////////////////////////////
// game server

type GameServer interface {
	util.Server
}

func NewGameServer(logger *zap.Logger, cookieServer share_api.CookieServer, service service.GameService) GameServer {
	hxServer := util.NewHxServer(logger, tpl)

	server := &gameServer{
		HxServer:     hxServer,
		CookieServer: cookieServer,
		logger:       logger,
		service:      service,
	}

	server.CreateGameServer = share_api.NewCreateGameServer(logger, service, server.OnCreateGame)
	server.JoinGameServer = share_api.NewJoinGameServer(logger, service, server.OnJoinGame)
	server.StartGameServer = share_api.NewStartGameServer(logger, service, server.OnStartGame)
	server.LeaveGameServer = share_api.NewLeaveGameServer(logger, service, server.OnLeaveGame)

	server.hub = share_websocket.NewHub[*model.Player](logger, server.WrapData, hxServer)

	cookieServer.RegisterOnCookie(server.onCookie)

	return server
}

type gameServer struct {
	util.HxServer
	share_api.CookieServer
	share_api.CreateGameServer[*model.Player]
	share_api.JoinGameServer[*model.Player]
	share_api.StartGameServer[*model.Player]
	share_api.LeaveGameServer[*model.Player]
	logger  *zap.Logger
	service service.GameService
	hub     share_websocket.Hub[*model.Player]
}

// //////////////////////////////////////////////////
// register

func (s *gameServer) RegisterRoutes(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, s.path(""), s.page_home)
	router.HandlerFunc(http.MethodGet, s.path("htmx/connect"), s.htmx_connect)
}

func (s *gameServer) path(path string) string {
	return fmt.Sprintf("/%s/%s", model.AppId, strings.TrimPrefix(path, "/"))
}

// //////////////////////////////////////////////////
// wrap data

func (s *gameServer) WrapData(data share_websocket.Data, player *model.Player) (bool, any) {
	data.With("share", share_api.NewRenderer())
	return s.service.WrapData(data, player)
}

// //////////////////////////////////////////////////
// on cookie

func (s *gameServer) onCookie(cookie *share_model.Cookie) {
	s.logger.Info("[on-cookie] ttt <<< ", zap.Any("cookie", cookie))

	playerId := cookie.PlayerId()
	player, err := s.hub.GetPlayer(playerId)
	if err != nil {
		s.logger.Error("player NOT found", zap.Any("cookie", cookie))
		return
	}
	player.SetCookie(cookie)

	s.broadcastPlayer(player)
	s.broadcastUser(cookie)
}

func (s *gameServer) broadcastUser(cookie *share_model.Cookie) {
	playerId := cookie.PlayerId()
	s.hub.BroadcastToPlayerRender(playerId, nil, s.CookieServer.RenderUser(cookie))
}

// //////////////////////////////////////////////////
// on player update

func (s *gameServer) onPlayerUpdate(playerId share_model.PlayerId) {
	player, err := s.hub.GetPlayer(playerId)
	if err != nil {
		s.logger.Error("player NOT found", zap.Any("id", playerId))
		return
	}
	s.broadcastPlayer(player)
}

func (s *gameServer) broadcastPlayer(player *model.Player) {
	s.hub.UpdatePlayer(player)
	s.broadcastJoinableGames()
	if player.HasGameId() {
		game, err := s.service.GetGame(player.GameId())
		if err == nil {
			s.broadcastPlayers(game)
			s.broadcastBoard(game)
		}
	}
}

// //////////////////////////////////////////////////
// on game events

func (s *gameServer) OnCreateGame(player *model.Player, game *model.Game) {
	s.broadcastGameLayoutToPlayer(player.Id(), game)
	s.OnGame(game)
}

func (s *gameServer) OnJoinGame(player *model.Player, game *model.Game) {
	s.broadcastGameLayoutToPlayer(player.Id(), game)
	s.OnGame(game)
}

func (s *gameServer) OnStartGame(player *model.Player, game *model.Game) {
	s.OnGame(game)
}

func (s *gameServer) OnLeaveGame(player *model.Player, game *model.Game) {
	s.broadcastJoinableGamesToPlayer(player.Id())
	s.OnGame(game)
}

func (s *gameServer) OnGame(game *model.Game) {
	if game != nil {
		s.broadcastGame(game)
	}
	s.broadcastJoinableGames()
}

// //////////////////////////////////////////////////
// broadcast

func (s *gameServer) broadcastInfoToPlayers(game *model.Game, info string) {
	s.hub.BroadcastToGamePlayers("info", game.Id(), share_websocket.Data{
		"info": info,
	})
}

func (s *gameServer) broadcastInfoToPlayer(playerId share_model.PlayerId, info string) {
	s.hub.BroadcastToPlayer("info", playerId, share_websocket.Data{
		"info": info,
	})
}

func (s *gameServer) broadcastErrorToPlayer(playerId share_model.PlayerId, err error) {
	s.hub.BroadcastToPlayer("error", playerId, share_websocket.Data{
		"error": err.Error(),
	})
}

func (s *gameServer) broadcastSelectGameToPlayer(playerId share_model.PlayerId) {
	data := s.getJoinableGamesData(playerId)
	s.hub.BroadcastToPlayer("select-game", playerId, data)
}

func (s *gameServer) broadcastGameLayoutToPlayer(playerId share_model.PlayerId, game *model.Game) {
	s.hub.BroadcastToPlayer("game-layout", playerId, share_websocket.Data{
		"game": game,
	})
}

func (s *gameServer) broadcastJoinableGamesToPlayer(playerId share_model.PlayerId) {
	data := s.getJoinableGamesData(playerId)
	s.hub.BroadcastToPlayer("select-game", playerId, data)
}

func (s *gameServer) broadcastJoinableGames() {
	s.hub.BroadcastToNotPlayingPlayersFn("select-game", func(player *model.Player) (bool, any) {
		data := s.getJoinableGamesData(player.Id())
		return s.hub.WrapPlayerData(data, player)
	})
}

func (s *gameServer) getJoinableGamesData(playerId share_model.PlayerId) websocket.Data {
	waitingPlayers := s.getWaitingPlayers(playerId)
	data := make(websocket.Data)
	data["new_games"] = s.service.GetJoinableGames()
	data["other_games"] = s.service.GetNonJoinableGames(playerId)
	data["has_waiting_players"] = len(waitingPlayers) > 0
	data["waiting_players"] = waitingPlayers
	return data
}

func (s *gameServer) getWaitingPlayers(playerId share_model.PlayerId) []*model.Player {
	players := s.hub.GetNotPlayingPlayers()
	waitingPlayers := make([]*model.Player, 0, len(players))
	for _, player := range players {
		if player == nil {
			continue
		}
		if player.Id() == playerId {
			continue
		}
		waitingPlayers = append(waitingPlayers, player)
	}
	return waitingPlayers
}

func (s *gameServer) broadcastGame(game *model.Game) {
	s.broadcastPlayers(game)
	s.broadcastBoard(game)
}

func (s *gameServer) broadcastPlayers(game *model.Game) {
	s.hub.BroadcastToGamePlayers("players", game.Id(), websocket.Data{
		"players": game.GetPlayers(),
	})
}

func (s *gameServer) broadcastBoard(game *model.Game) {
	s.hub.BroadcastToGamePlayers("board", game.Id(), websocket.Data{
		"game": game,
	})
}

// //////////////////////////////////////////////////
// render

func (s *gameServer) renderInfo(w io.Writer, info string) {
	s.Render(w, "info", websocket.Data{
		"info": info,
	})
}

func (s *gameServer) renderError(w io.Writer, err error) {
	s.Render(w, "error", websocket.Data{
		"error": err.Error(),
	})
}
