package websocket

import (
	"fmt"
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

	GetUser(id model.UserId) (User, error)
	RegisterUser(user User)
	UnregisterUserId(id model.UserId)
	UpdateUser(user User)

	UpdateUserFromPlayer(player PlayerT)
	OnUserUpdate(userId model.UserId)

	GetPlayer(playerId model.PlayerId) (PlayerT, error)

	BroadcastInfoToUser(userId model.UserId, info string)
	BroadcastErrorToUser(userId model.UserId, err error)
	BroadcastInfoToPlayers(game GameT, info string)
	BroadcastJoinableGamesToUser(userId model.UserId)
	BroadcastJoinableGames()
	BroadcastGameLayoutToPlayer(playerId model.PlayerId, game GameT)
	BroadcastGame(game GameT)
	BroadcastPlayers(game GameT)
	BroadcastBoard(game GameT)
	BroadcastPlayer(player PlayerT)
	BroadcastCookie(cookie *model.Cookie)
	BroadcastUserCookie(cookie *model.Cookie, renderUserFn func(cookie *model.Cookie) func(w io.Writer, data model.Data))

	OnJoinGame(game GameT, player PlayerT)
	OnGame(game GameT)
	OnLeaveGame(game GameT, userId model.UserId)
}

type Game[PlayerT Player] interface {
	Id() model.GameId
	Player(id model.PlayerId) (PlayerT, bool)
	Players() []PlayerT
}

type CookieServer interface {
	GetValidCookie(r *http.Request) (*model.Cookie, error)
	RenderUser(cookie *model.Cookie) func(w io.Writer, data model.Data)
}

func NewHubServer[PlayerT Player, GameT Game[PlayerT]](logger *zap.Logger, hub Hub[PlayerT], cookierServer CookieServer, newUserFromCookieFn func(cookier *model.Cookie) User, service Service[PlayerT, GameT]) HubServer[PlayerT, GameT] {
	server := &hubServer[PlayerT, GameT]{
		logger:              logger,
		hub:                 hub,
		cookierServer:       cookierServer,
		newUserFromCookieFn: newUserFromCookieFn,
		service:             service,
	}

	service.RegisterOnJoinGame(server.OnJoinGame)
	service.RegisterOnGame(server.OnGame)
	service.RegisterOnLeaveGame(server.OnLeaveGame)

	return server
}

type hubServer[PlayerT Player, GameT Game[PlayerT]] struct {
	logger              *zap.Logger
	hub                 Hub[PlayerT]
	cookierServer       CookieServer
	newUserFromCookieFn func(cookier *model.Cookie) User
	service             Service[PlayerT, GameT]
}

type Service[PlayerT Player, GameT Game[PlayerT]] interface {
	GetGame(gameId model.GameId) (GameT, error)
	GetJoinableGames() []GameT
	GetNonJoinableGames(userId model.UserId) []GameT

	RegisterOnJoinGame(func(game GameT, player PlayerT))
	RegisterOnGame(func(game GameT))
	RegisterOnLeaveGame(func(game GameT, userId model.UserId))
}

// //////////////////////////////////////////////////
// routes

func (s *hubServer[PlayerT, GameT]) RegisterAppRoutes(router *httprouter.Router, app model.App) {
	s.logger.Info(fmt.Sprintf(" (+) GET %s", app.HtmxConnectRoute()))
	router.HandlerFunc(http.MethodGet, app.HtmxConnectRoute(), s.HtmxConnect)
}

// //////////////////////////////////////////////////
// hub

func (s *hubServer[PlayerT, GameT]) Hub() Hub[PlayerT] {
	return s.hub
}

// //////////////////////////////////////////////////
// user

func (s *hubServer[PlayerT, GameT]) GetUser(id model.UserId) (User, error) {
	return s.hub.GetUser(id)
}

func (s *hubServer[PlayerT, GameT]) RegisterUser(user User) {
	s.hub.RegisterUser(user)
}

func (s *hubServer[PlayerT, GameT]) UnregisterUserId(id model.UserId) {
	s.hub.UnregisterUserId(id)
}

func (s *hubServer[PlayerT, GameT]) UpdateUser(user User) {
	s.hub.UpdateUser(user)
}

// //////////////////////////////////////////////////
// player

func (s *hubServer[PlayerT, GameT]) UpdateUserFromPlayer(player PlayerT) {
	userId := player.Id().UserId()
	user, err := s.GetUser(userId)
	if err != nil {
		s.logger.Error("user NOT found", zap.Any("id", userId))
		return
	}
	user.SetGameId(player.GameId())
	s.UpdateUser(user)
}

func (s *hubServer[PlayerT, GameT]) OnUserUpdate(userId model.UserId) {
	user, err := s.hub.GetUser(userId)
	if err != nil {
		return
	}
	s.BroadcastUser(user)
}

func (s *hubServer[PlayerT, GameT]) GetPlayer(playerId model.PlayerId) (PlayerT, error) {
	return s.hub.GetPlayer(playerId)
}

// //////////////////////////////////////////////////
// broadcast

func (s *hubServer[PlayerT, GameT]) BroadcastCookie(cookie *model.Cookie) {

	s.logger.Info("[on-cookie] broadcast cookie...", zap.Any("cookie", cookie))
	s.BroadcastUserCookie(cookie, s.cookierServer.RenderUser)

	userId := cookie.Id
	s.logger.Info(fmt.Sprintf("[on-cookie] broadcast user %s...", userId), zap.Any("cookie", cookie))
	user, err := s.Hub().GetUser(userId)
	if err != nil {
		s.logger.Info(fmt.Sprintf("[on-cookie] user %s NOT found >>> SKIPPED", userId), zap.Any("cookie", cookie), zap.Error(err))
		return
	}
	s.logger.Info(fmt.Sprintf("[on-cookie] update user %s + broadcast", userId), zap.Any("cookie", cookie))
	user.SetCookie(cookie)

	s.BroadcastUser(user)

	if !user.HasGameId() {
		return
	}

	playerId := user.PlayerId()
	s.logger.Info(fmt.Sprintf("[on-cookie] broadcast player %s...", playerId), zap.Any("cookie", cookie))
	player, err := s.Hub().GetPlayer(playerId)
	if err != nil {
		s.logger.Info(fmt.Sprintf("[on-cookie] player %s NOT found >>> SKIPPED", playerId), zap.Any("cookie", cookie), zap.Error(err))
		return
	}
	s.logger.Info(fmt.Sprintf("[on-cookie] update player %s + broadcast", playerId), zap.Any("cookie", cookie))
	player.User().SetCookie(cookie)

	s.BroadcastPlayer(player)
}

// //////////////////////////////////////////////////
// broadcast

func (s *hubServer[PlayerT, GameT]) BroadcastInfoToUser(userId model.UserId, info string) {
	s.hub.BroadcastToUser("info", userId, model.Data{
		"Info": info,
	})
}

func (s *hubServer[PlayerT, GameT]) BroadcastErrorToUser(userId model.UserId, err error) {
	s.hub.BroadcastToUser("error", userId, model.Data{
		"Error": err.Error(),
	})
}

func (s *hubServer[PlayerT, GameT]) BroadcastInfoToPlayers(game GameT, info string) {
	s.hub.BroadcastToGamePlayers("info", game.Id(), model.Data{
		"Info": info,
	})
}

func (s *hubServer[PlayerT, GameT]) BroadcastJoinableGamesToUser(userId model.UserId) {
	data := s.getJoinableGamesData(userId)
	s.hub.BroadcastToUser("select-game", userId, data)
}

func (s *hubServer[PlayerT, GameT]) BroadcastJoinableGames() {
	s.hub.BroadcastToNotPlayingUsersFn("select-game", func(user User) (bool, model.Data) {
		data := s.getJoinableGamesData(user.Id())
		return s.hub.WrapUserData(data, user)
	})
}

func (s *hubServer[PlayerT, GameT]) getJoinableGamesData(userId model.UserId) model.Data {
	waitingUsers := s.getWaitingUsers(userId)
	return model.Data{
		"NewGames":        s.service.GetJoinableGames(),
		"OtherGames":      s.service.GetNonJoinableGames(userId),
		"HasWaitingUsers": len(waitingUsers) > 0,
		"WaitingUsers":    waitingUsers,
	}
}

func (s *hubServer[PlayerT, GameT]) getWaitingUsers(userId model.UserId) []User {
	users := s.hub.GetNotPlayingUsers()
	waitingUsers := make([]User, 0, len(users))
	for _, user := range users {
		if user.IsNotUser(userId) {
			waitingUsers = append(waitingUsers, user)
		}
	}
	return waitingUsers
}

func (s *hubServer[PlayerT, GameT]) BroadcastGameLayoutToPlayer(playerId model.PlayerId, game GameT) {
	s.hub.BroadcastToPlayer("game-layout", playerId, model.Data{
		"Game": game,
	})
}

func (s *hubServer[PlayerT, GameT]) BroadcastGame(game GameT) {
	s.BroadcastPlayers(game)
	s.BroadcastBoard(game)
}

func (s *hubServer[PlayerT, GameT]) BroadcastPlayers(game GameT) {
	s.hub.BroadcastToGamePlayers("players", game.Id(), model.Data{
		"Players": game.Players(),
	})
}

func (s *hubServer[PlayerT, GameT]) BroadcastBoard(game GameT) {
	s.hub.BroadcastToGamePlayers("board", game.Id(), model.Data{
		"Game": game,
	})
}

func (s *hubServer[PlayerT, GameT]) BroadcastPlayer(player PlayerT) {
	s.UpdateUserFromPlayer(player)
	s.BroadcastJoinableGames()
	game, err := s.service.GetGame(player.GameId())
	if err == nil {
		s.BroadcastGame(game)
	}
}

func (s *hubServer[PlayerT, GameT]) BroadcastUserCookie(cookie *model.Cookie, renderCookieFn func(cookie *model.Cookie) func(w io.Writer, data model.Data)) {
	s.Hub().BroadcastToUserRender(cookie.Id, nil, renderCookieFn(cookie))
}

func (s *hubServer[PlayerT, GameT]) BroadcastUser(user User) {
	s.BroadcastJoinableGames()
}

// //////////////////////////////////////////////////
// on game events

func (s *hubServer[PlayerT, GameT]) OnJoinGame(game GameT, player PlayerT) {

	userId := player.User().Id()
	user, err := s.GetUser(userId)
	if err != nil {
		return
	}
	user.SetGameId(game.Id())

	s.BroadcastGameLayoutToPlayer(player.Id(), game)
	s.OnGame(game)
}

func (s *hubServer[PlayerT, GameT]) OnLeaveGame(game GameT, userId model.UserId) {

	user, err := s.GetUser(userId)
	if err != nil {
		return
	}
	user.UnsetGameId()

	s.OnGame(game)
	s.BroadcastJoinableGamesToUser(userId)
}

func (s *hubServer[PlayerT, GameT]) OnGame(game GameT) {
	s.BroadcastGame(game)
	s.BroadcastJoinableGames()
}
