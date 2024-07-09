package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/util"

	share_api "github.com/gre-ory/games-go/internal/game/share/api"
	share_model "github.com/gre-ory/games-go/internal/game/share/model"
	share_websocket "github.com/gre-ory/games-go/internal/game/share/websocket"

	"github.com/gre-ory/games-go/internal/game/czm/model"
	"github.com/gre-ory/games-go/internal/game/czm/service"
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
		GameServer:   share_api.NewGameServer(logger, service),
		logger:       logger,
		service:      service,
	}

	hub := share_websocket.NewHub[*model.Player](logger, server.WrapData, hxServer)
	server.HubServer = share_websocket.NewHubServer[*model.Player, *model.Game](hub, service)

	server.CookieServer.RegisterOnCookie(server.onCookie)

	return server
}

type gameServer struct {
	util.HxServer
	share_api.CookieServer
	share_websocket.HubServer[*model.Player, *model.Game]
	share_api.GameServer[*model.Player]
	logger  *zap.Logger
	service service.GameService
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
	s.logger.Info(fmt.Sprintf("[on-cookie] %s <<< ", model.AppId), zap.Any("cookie", cookie))

	playerId := cookie.PlayerId()
	player, err := s.GetPlayer(playerId)
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
	s.Hub().BroadcastToPlayerRender(playerId, nil, s.CookieServer.RenderUser(cookie))
}

// //////////////////////////////////////////////////
// on player update

func (s *gameServer) onPlayerUpdate(playerId share_model.PlayerId) {
	player, err := s.GetPlayer(playerId)
	if err != nil {
		s.logger.Error("player NOT found", zap.Any("id", playerId))
		return
	}
	s.broadcastPlayer(player)
}

func (s *gameServer) broadcastPlayer(player *model.Player) {
	s.UpdatePlayer(player)
	s.BroadcastJoinableGames()
	if player.HasGameId() {
		game, err := s.service.GetGame(player.GameId())
		if err == nil {
			s.BroadcastGame(game)
		}
	}
}

// //////////////////////////////////////////////////
// on game events

func (s *gameServer) OnCreateGame(player *model.Player, game *model.Game) {
	s.BroadcastGameLayoutToPlayer(player.Id(), game)
	s.OnGame(game)
}

func (s *gameServer) OnJoinGame(player *model.Player, game *model.Game) {
	s.BroadcastGameLayoutToPlayer(player.Id(), game)
	s.OnGame(game)
}

func (s *gameServer) OnStartGame(player *model.Player, game *model.Game) {
	s.OnGame(game)
}

func (s *gameServer) OnLeaveGame(player *model.Player, game *model.Game) {
	s.BroadcastJoinableGamesToPlayer(player.Id())
	s.OnGame(game)
}

func (s *gameServer) OnGame(game *model.Game) {
	if game != nil {
		s.BroadcastGame(game)
	}
	s.BroadcastJoinableGames()
}
