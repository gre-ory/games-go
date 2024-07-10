package api

import (
	"fmt"
	"net/http"

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

	hub := share_websocket.NewHub(logger, server.WrapData, hxServer)
	server.HubServer = share_websocket.NewHubServer(logger, hub, cookieServer, server.newPlayerFromCookie, service)

	server.CookieServer.RegisterOnCookie(server.onCookie)

	return server
}

type gameServer struct {
	util.HxServer
	share_api.CookieServer
	share_websocket.HubServer[*model.Player, *model.Game]
	share_api.GameServer[*model.Player, *model.Game]
	logger  *zap.Logger
	service service.GameService
}

// //////////////////////////////////////////////////
// register

func (s *gameServer) RegisterRoutes(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, model.App.HomeRoute(), s.page_home())
	s.HubServer.RegisterAppRoutes(router, model.App)
}

func (s *gameServer) page_home() func(http.ResponseWriter, *http.Request) {
	return share_api.PageHome(s.logger, model.App, s, s)
}

// //////////////////////////////////////////////////
// wrap data

func (s *gameServer) WrapData(data share_websocket.Data, player *model.Player) (bool, any) {
	data.With("share", share_api.NewRenderer())
	return s.service.WrapData(data, player)
}

// //////////////////////////////////////////////////
// cookie

func (s *gameServer) newPlayerFromCookie(cookie *share_model.Cookie) *model.Player {
	ws_player := share_websocket.NewPlayerFromCookie(s.logger, cookie, s.onMessage, s.OnPlayerUpdate, nil)
	return model.NewPlayer(ws_player)
}

func (s *gameServer) onCookie(cookie *share_model.Cookie) {

	playerId := cookie.PlayerId()
	player, err := s.Hub().GetPlayer(playerId)
	if err != nil {
		s.logger.Info(fmt.Sprintf("[on-cookie] %s :: player %s NOT found >>> SKIPPED", model.App.Id(), playerId), zap.Any("cookie", cookie))
		return
	}
	s.logger.Info(fmt.Sprintf("[on-cookie] %s :: update player %s + broadcast", model.App.Id(), playerId), zap.Any("cookie", cookie))

	player.SetCookie(cookie)

	s.BroadcastPlayer(player)
	s.BroadcastPlayerCookie(cookie, s.CookieServer.RenderUser)
}
