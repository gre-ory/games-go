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

	"github.com/gre-ory/games-go/internal/game/ttt/model"
	"github.com/gre-ory/games-go/internal/game/ttt/service"
)

// //////////////////////////////////////////////////
// game server

type GameServer interface {
	util.Server
	share_websocket.HubServer[*model.Player, *model.Game]
}

func NewGameServer(logger *zap.Logger, cookieServer share_api.CookieServer, service service.GameService) GameServer {
	logger = model.App.Logger(logger)
	hxServer := util.NewHxServer(logger, tpl)

	server := &gameServer{
		HxServer:     hxServer,
		CookieServer: cookieServer,
		GameServer:   share_api.NewGameServer(logger, service),
		logger:       logger,
		service:      service,
	}

	hub := share_websocket.NewHub(logger, server.WrapUserData, service.GetPlayer, server.WrapPlayerData, hxServer)
	server.HubServer = share_websocket.NewHubServer(logger, hub, cookieServer, server.newUserFromCookie, service)

	server.CookieServer.RegisterOnCookie(server.BroadcastCookie)

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
	s.logger.Info(fmt.Sprintf(" (+) GET %s", model.App.HomeRoute()))
	router.HandlerFunc(http.MethodGet, model.App.HomeRoute(), s.page_home())
	s.HubServer.RegisterAppRoutes(router, model.App)
}

func (s *gameServer) page_home() func(http.ResponseWriter, *http.Request) {
	return share_api.PageHome(s.logger, model.App, s, s)
}

// //////////////////////////////////////////////////
// wrap data

func (s *gameServer) WrapUserData(data share_model.Data, user share_model.User) (bool, share_model.Data) {
	data = data.With("Share", share_api.NewRenderer())
	if user == nil {
		return true, data
	}
	data = data.With("Lang", model.App.UserLocalizer(user))
	return true, data
}

func (s *gameServer) WrapPlayerData(data share_model.Data, player *model.Player) (bool, share_model.Data) {
	ok, data := s.WrapUserData(data, player.User())
	if !ok {
		return false, nil
	}
	if game, err := s.service.GetGame(player.GameId()); err != nil {
		return false, nil
	} else {
		data = data.With("Game", game)
	}
	return true, data
}

// //////////////////////////////////////////////////
// cookie

func (s *gameServer) newUserFromCookie(cookie *share_model.Cookie) share_websocket.User {
	return share_websocket.NewUser(s.logger, cookie, s.onMessage, s.OnUserUpdate, nil)
}
