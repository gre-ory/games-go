package api

import (
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/util"
	"github.com/gre-ory/games-go/internal/util/websocket"

	user_api "github.com/gre-ory/games-go/internal/game/user/api"

	"github.com/gre-ory/games-go/internal/game/tictactoe/model"
	"github.com/gre-ory/games-go/internal/game/tictactoe/service"
)

// //////////////////////////////////////////////////
// game server

type GameServer interface {
	util.Server
}

func NewGameServer(logger *zap.Logger, sessionServer user_api.SessionServer, service service.GameService, cookieSecret string) GameServer {
	hxServer := util.NewHxServer(logger, tpl)
	cookieHelper := NewCookieHelper(logger, cookieSecret)
	return &gameServer{
		HxServer:      hxServer,
		CookieHelper:  cookieHelper,
		logger:        logger,
		sessionServer: sessionServer,
		service:       service,
		hub:           websocket.NewHub[model.PlayerId, model.GameId, *model.Player](logger, service.WrapData, hxServer),
	}
}

type gameServer struct {
	util.HxServer
	util.CookieHelper[Cookie]
	logger        *zap.Logger
	sessionServer user_api.SessionServer
	service       service.GameService
	hub           websocket.Hub[model.PlayerId, model.GameId, *model.Player]
}

// //////////////////////////////////////////////////
// register

func (s *gameServer) RegisterRoutes(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/game/tictactoe", s.page_home)

	router.HandlerFunc(http.MethodGet, "/htmx/tictactoe/connect", s.htmx_connect)
	// router.HandlerFunc(http.MethodPut, "/htmx/tictactoe/create", s.htmx_create_game)
	// router.HandlerFunc(http.MethodPut, "/htmx/tictactoe/join/:game_id", s.htmx_join_game)
	// router.HandlerFunc(http.MethodPut, "/htmx/tictactoe/play/:game_id/:play_x/:play_y", s.htmx_play_game)
	router.HandlerFunc(http.MethodDelete, "/htmx/tictactoe/leave/:game_id", s.htmx_leave_game)
	router.HandlerFunc(http.MethodDelete, "/htmx/tictactoe/delete/:game_id", s.htmx_delete_game)
}

// //////////////////////////////////////////////////
// broadcast

func (s *gameServer) broadcastErrorToPlayer(playerId model.PlayerId, err error) {
	s.hub.BroadcastToPlayer("error", playerId, websocket.Data{
		"error": err.Error(),
	})
}

func (s *gameServer) broadcastSelectNameToPlayer(playerId model.PlayerId) {
	s.hub.BroadcastToPlayer("select-name", playerId, nil)
}

func (s *gameServer) broadcastSelectGameToPlayer(playerId model.PlayerId) {
	data := make(websocket.Data)
	data["games"] = s.service.GetJoinableGames()
	s.hub.BroadcastToPlayer("select-game", playerId, data)
}

func (s *gameServer) broadcastGameLayoutToPlayer(playerId model.PlayerId, game *model.Game) {
	s.hub.BroadcastToPlayer("game-layout", playerId, websocket.Data{
		"game": game,
	})
}

func (s *gameServer) broadcastJoinableGamesToPlayer(playerId model.PlayerId) {
	data := make(websocket.Data)
	data["games"] = s.service.GetJoinableGames()
	s.hub.BroadcastToPlayer("select-game", playerId, data)
}

func (s *gameServer) broadcastJoinableGames() {
	games := s.service.GetJoinableGames()
	s.hub.BroadcastToNotPlayingPlayers("select-game", websocket.Data{
		"games": games,
	})
}

func (s *gameServer) broadcastGame(game *model.Game) {
	s.broadcastPlayers(game)
	s.broadcastBoard(game)
}

func (s *gameServer) broadcastPlayers(game *model.Game) {
	s.hub.BroadcastToGamePlayers("players", game.Id, websocket.Data{
		"players": game.Players,
	})
}

func (s *gameServer) broadcastBoard(game *model.Game) {
	s.hub.BroadcastToGamePlayers("board", game.Id, websocket.Data{
		"game": game,
	})
}

// //////////////////////////////////////////////////
// render

func (s *gameServer) renderSelectName(w io.Writer, cookie *Cookie) {
	s.Render(w, "select-name", cookie.Data())
}

func (s *gameServer) renderSelectGame(w io.Writer, cookie *Cookie) {
	data := cookie.Data()
	data["games"] = s.service.GetJoinableGames()
	s.Render(w, "select-game", data)
}

func (s *gameServer) renderGameLayout(w io.Writer, game *model.Game, cookie *Cookie) {
	data := cookie.Data()
	data["game"] = game
	s.Render(w, "game-layout", data)
}

func (s *gameServer) renderError(w io.Writer, err error) {
	s.Render(w, "error", websocket.Data{
		"error": err.Error(),
	})
}
