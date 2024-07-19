package websocket

import (
	"fmt"
	"io"
	"sync"

	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/util"

	"github.com/gre-ory/games-go/internal/game/share/model"
)

const (
	DebugLock      = false
	DebugBroadcast = false
	DebugPing      = false
	DebugMessage   = false
)

type Hub[PlayerT Player] interface {
	GetUser(id model.UserId) (User, error)
	GetUsers() []User
	GetInactiveUsers() []User
	GetNotPlayingUsers() []User
	GetPlayingUsers() []User
	FilterUsers(filterFn func(user User) bool) []User
	RegisterUser(user User)
	UnregisterUserId(id model.UserId)
	UpdateUser(user User)

	GetPlayer(id model.PlayerId) (PlayerT, error)
	GetPlayers() []PlayerT
	GetGamePlayers(gameId model.GameId) []PlayerT
	FilterPlayers(filterFn func(player PlayerT) bool) []PlayerT

	BroadcastToUser(name string, id model.UserId, data model.Data)
	BroadcastToUserFn(name string, id model.UserId, acceptFn func(user User) (bool, model.Data))
	BroadcastToUserRender(id model.UserId, data model.Data, renderFn func(w io.Writer, data model.Data))
	BroadcastToUserRenderFn(id model.UserId, acceptFn func(user User) (bool, model.Data), renderFn func(w io.Writer, data model.Data))
	BroadcastToUsers(name string, data model.Data)
	BroadcastToUsersFn(name string, acceptFn func(user User) (bool, model.Data))
	BroadcastUsersRender(data model.Data, renderFn func(w io.Writer, data model.Data))
	BroadcastUsersRenderFn(acceptFn func(user User) (bool, model.Data), renderFn func(w io.Writer, data model.Data))
	BroadcastToNotPlayingUsers(name string, data model.Data)
	BroadcastToNotPlayingUsersFn(name string, acceptFn func(user User) (bool, model.Data))
	BroadcastToPlayingUsers(name string, data model.Data)
	BroadcastToPlayingUsersFn(name string, acceptFn func(user User) (bool, model.Data))
	WrapUserData(data model.Data, user User) (bool, model.Data)

	BroadcastToPlayer(name string, id model.PlayerId, data model.Data)
	BroadcastToPlayerFn(name string, id model.PlayerId, acceptFn func(player PlayerT) (bool, model.Data))
	BroadcastToPlayerRender(id model.PlayerId, data model.Data, renderFn func(w io.Writer, data model.Data))
	BroadcastToPlayerRenderFn(id model.PlayerId, acceptFn func(player PlayerT) (bool, model.Data), renderFn func(w io.Writer, data model.Data))
	BroadcastToPlayers(name string, data model.Data)
	BroadcastToPlayersFn(name string, acceptFn func(player PlayerT) (bool, model.Data))
	BroadcastPlayersRender(data model.Data, renderFn func(w io.Writer, data model.Data))
	BroadcastPlayersRenderFn(acceptFn func(player PlayerT) (bool, model.Data), renderFn func(w io.Writer, data model.Data))
	BroadcastToGamePlayers(name string, gameId model.GameId, data model.Data)
	BroadcastToGamePlayersFn(name string, gameId model.GameId, acceptFn func(player PlayerT) (bool, model.Data))
	WrapPlayerData(data model.Data, player PlayerT) (bool, model.Data)
}

type Player interface {
	User() model.User
	Id() model.PlayerId
	GameId() model.GameId
}

func NewHub[PlayerT Player](logger *zap.Logger, wrapUserDataFn func(data model.Data, user model.User) (bool, model.Data), getPlayerFn func(model.PlayerId) (PlayerT, error), wrapPlayerDataFn func(data model.Data, player PlayerT) (bool, model.Data), tplRenderer util.TplRenderer) Hub[PlayerT] {
	h := &hub[PlayerT]{
		TplRenderer:      tplRenderer,
		broadcastUser:    make(chan TplRenderer[User]),
		broadcastPlayer:  make(chan TplRenderer[PlayerT]),
		registerUser:     make(chan User),
		unregisterUserId: make(chan model.UserId),
		logger:           logger,
		users:            make(map[model.UserId]User),
		getPlayerFn:      getPlayerFn,
		wrapUserDataFn:   wrapUserDataFn,
		wrapPlayerDataFn: wrapPlayerDataFn,
	}
	go h.run()
	return h
}

type hub[PlayerT Player] struct {
	util.TplRenderer

	broadcastUser    chan TplRenderer[User]
	broadcastPlayer  chan TplRenderer[PlayerT]
	registerUser     chan User
	unregisterUserId chan model.UserId

	logger *zap.Logger
	users  map[model.UserId]User
	mutex  sync.RWMutex

	getPlayerFn      func(model.PlayerId) (PlayerT, error)
	wrapUserDataFn   func(data model.Data, user model.User) (bool, model.Data)
	wrapPlayerDataFn func(data model.Data, player PlayerT) (bool, model.Data)
}

// //////////////////////////////////////////////////
// run

func (h *hub[PlayerT]) run() {
	for {
		select {
		case user := <-h.registerUser:
			h.onRegisterUser(user)
		case userId := <-h.unregisterUserId:
			h.onUnregisterUser(userId)
		case tplUser := <-h.broadcastUser:
			h.onBroadcastUser(tplUser)
		case tplPlayer := <-h.broadcastPlayer:
			h.onBroadcastPlayer(tplPlayer)
		}
	}
}

// //////////////////////////////////////////////////
// user

func (h *hub[PlayerT]) GetUser(id model.UserId) (User, error) {
	unlock := h.rlock("GetUser")
	defer unlock()

	user, found := h.users[id]
	if !found {
		return nil, model.ErrUserNotFound
	}

	return user, nil
}

func (h *hub[PlayerT]) GetUsers() []User {
	return h.FilterUsers(nil)
}

func (h *hub[PlayerT]) GetInactiveUsers() []User {
	return h.FilterUsers(User.IsInactive)
}

func (h *hub[PlayerT]) GetNotPlayingUsers() []User {
	return h.FilterUsers(User.IsNotPlaying)
}

func (h *hub[PlayerT]) GetPlayingUsers() []User {
	return h.FilterUsers(User.IsPlaying)
}

func (h *hub[PlayerT]) FilterUsers(filterFn func(user User) bool) []User {
	unlock := h.rlock("FilterUsers")
	defer unlock()

	users := make([]User, 0, len(h.users))
	for _, user := range h.users {
		if filterFn == nil || filterFn(user) {
			users = append(users, user)
		}
	}
	return users
}

// //////////////////////////////////////////////////
// register user

func (h *hub[PlayerT]) RegisterUser(user User) {
	h.registerUser <- user
}

func (h *hub[PlayerT]) onRegisterUser(user User) {
	unlock := h.lock("onRegisterUser")
	defer unlock()

	h.logger.Info(fmt.Sprintf("[register] (+) user %v", user.Id()))
	h.users[user.Id()] = user
}

// //////////////////////////////////////////////////
// unregister user

func (h *hub[PlayerT]) UnregisterUserId(id model.UserId) {
	h.unregisterUserId <- id
}

func (h *hub[PlayerT]) onUnregisterUser(id model.UserId) {
	unlock := h.lock("onUnregisterUser")
	defer unlock()

	h.logger.Info(fmt.Sprintf("[unregister] (-) user %v", id))
	delete(h.users, id)
}

// //////////////////////////////////////////////////
// update user

func (h *hub[PlayerT]) UpdateUser(user User) {
	unlock := h.lock("UpdateUser")
	defer unlock()

	h.logger.Info(fmt.Sprintf("[update] (~) user %v", user.Id()))
	h.users[user.Id()] = user
}

// //////////////////////////////////////////////////
// players

func (h *hub[PlayerT]) GetPlayer(id model.PlayerId) (PlayerT, error) {
	return h.getPlayerFn(id)
}

func (h *hub[PlayerT]) GetPlayers() []PlayerT {
	return h.FilterPlayers(nil)
}

func (h *hub[PlayerT]) GetGamePlayers(gameId model.GameId) []PlayerT {
	return h.FilterPlayers(func(player PlayerT) bool {
		return player.GameId() == gameId
	})
}

func (h *hub[PlayerT]) FilterPlayers(filterFn func(player PlayerT) bool) []PlayerT {
	unlock := h.rlock("FilterPlayers")
	defer unlock()

	players := make([]PlayerT, 0, len(h.users))
	for _, user := range h.users {
		if !user.HasGameId() {
			continue
		}
		playerId := user.PlayerId()
		player, err := h.getPlayerFn(playerId)
		if err != nil {
			continue
		}
		if filterFn == nil || filterFn(player) {
			players = append(players, player)
		}
	}
	return players
}

// //////////////////////////////////////////////////
// broadcast users

func (h *hub[PlayerT]) BroadcastToUser(name string, id model.UserId, data model.Data) {
	h.BroadcastToUserFn(name, id, func(user User) (bool, model.Data) {
		return h.WrapUserData(data, user)
	})
}

func (h *hub[PlayerT]) BroadcastToUserFn(name string, id model.UserId, acceptFn func(user User) (bool, model.Data)) {
	if DebugBroadcast {
		h.logger.Info(fmt.Sprintf("[broadcast] user <<< tpl %s - user %v", name, id))
	}
	h.broadcastUser <- h.NewNamedUserTemplate(
		name,
		h.AcceptUserFn(id, acceptFn),
	)
}

func (h *hub[PlayerT]) BroadcastToUserRender(id model.UserId, data model.Data, renderFn func(w io.Writer, data model.Data)) {
	h.BroadcastToUserRenderFn(id, func(user User) (bool, model.Data) {
		return h.WrapUserData(data, user)
	}, renderFn)
}

func (h *hub[PlayerT]) BroadcastToUserRenderFn(id model.UserId, acceptFn func(user User) (bool, model.Data), renderFn func(w io.Writer, data model.Data)) {
	if DebugBroadcast {
		h.logger.Info(fmt.Sprintf("[broadcast] user <<< renderFn - user %v", id))
	}
	h.broadcastUser <- h.NewTplUserRenderer(
		h.AcceptUserFn(id, acceptFn),
		renderFn,
	)
}

func (h *hub[PlayerT]) AcceptUserFn(userId model.UserId, acceptFn func(user User) (bool, model.Data)) func(user User) (bool, model.Data) {
	return func(user User) (bool, model.Data) {
		if user.IsUser(userId) {
			if acceptFn != nil {
				return acceptFn(user)
			}
			return true, nil
		}
		return false, nil
	}
}

func (h *hub[PlayerT]) BroadcastToUsers(name string, data model.Data) {
	h.BroadcastToUsersFn(name, func(user User) (bool, model.Data) {
		return h.WrapUserData(data, user)
	})
}

func (h *hub[PlayerT]) BroadcastToUsersFn(name string, acceptFn func(user User) (bool, model.Data)) {
	if DebugBroadcast {
		h.logger.Info(fmt.Sprintf("[broadcast] user <<< tpl %v", name))
	}
	h.broadcastUser <- h.NewNamedUserTemplate(
		name,
		acceptFn,
	)
}

func (h *hub[PlayerT]) BroadcastUsersRender(data model.Data, renderFn func(w io.Writer, data model.Data)) {
	h.BroadcastUsersRenderFn(func(user User) (bool, model.Data) {
		return h.WrapUserData(data, user)
	}, renderFn)
}

func (h *hub[PlayerT]) BroadcastUsersRenderFn(acceptFn func(user User) (bool, model.Data), renderFn func(w io.Writer, data model.Data)) {
	if DebugBroadcast {
		h.logger.Info("[broadcast] user <<< renderFn")
	}
	h.broadcastUser <- h.NewTplUserRenderer(
		acceptFn,
		renderFn,
	)
}

func (h *hub[PlayerT]) BroadcastToNotPlayingUsers(name string, data model.Data) {
	h.BroadcastToNotPlayingUsersFn(name, func(user User) (bool, model.Data) {
		return h.WrapUserData(data, user)
	})
}

func (h *hub[PlayerT]) BroadcastToNotPlayingUsersFn(name string, acceptFn func(user User) (bool, model.Data)) {
	if DebugBroadcast {
		h.logger.Info(fmt.Sprintf("[broadcast] user <<< tpl %s - not playing users", name))
	}
	h.broadcastUser <- h.NewNamedUserTemplate(
		name,
		h.AcceptNotPlayingUsersFn(acceptFn),
	)
}

func (h *hub[PlayerT]) AcceptNotPlayingUsersFn(acceptFn func(user User) (bool, model.Data)) func(user User) (bool, model.Data) {
	return func(user User) (bool, model.Data) {
		if user.IsNotPlaying() {
			if acceptFn != nil {
				return acceptFn(user)
			}
			return true, nil
		}
		return false, nil
	}
}

func (h *hub[PlayerT]) BroadcastToPlayingUsers(name string, data model.Data) {
	h.BroadcastToPlayingUsersFn(name, func(user User) (bool, model.Data) {
		return h.WrapUserData(data, user)
	})
}

func (h *hub[PlayerT]) BroadcastToPlayingUsersFn(name string, acceptFn func(user User) (bool, model.Data)) {
	if DebugBroadcast {
		h.logger.Info(fmt.Sprintf("[broadcast] user <<< tpl %s - playing users", name))
	}
	h.broadcastUser <- h.NewNamedUserTemplate(
		name,
		h.AcceptPlayingUsersFn(acceptFn),
	)
}

func (h *hub[PlayerT]) AcceptPlayingUsersFn(acceptFn func(user User) (bool, model.Data)) func(user User) (bool, model.Data) {
	return func(user User) (bool, model.Data) {
		if user.IsPlaying() {
			if acceptFn != nil {
				return acceptFn(user)
			}
			return true, nil
		}
		return false, nil
	}
}

func (h *hub[PlayerT]) NewNamedUserTemplate(name string, acceptFn func(user User) (bool, model.Data)) TplRenderer[User] {
	return h.NewTplUserRenderer(acceptFn, h.NewNamedRenderFn(name))
}

func (h *hub[PlayerT]) NewTplUserRenderer(acceptFn func(user User) (bool, model.Data), renderFn func(w io.Writer, data model.Data)) TplRenderer[User] {
	return NewTplRenderer[User](acceptFn, renderFn)
}

func (h *hub[PlayerT]) WrapUserData(data model.Data, user User) (bool, model.Data) {
	data = data.With("User", user)
	if h.wrapUserDataFn != nil {
		return h.wrapUserDataFn(data, user)
	}
	return true, data
}

func (h *hub[PlayerT]) onBroadcastUser(tpl TplRenderer[User]) {
	unlock := h.rlock("onBroadcastUser")
	defer unlock()

	for _, user := range h.users {
		if bytes, ok := tpl.Render(user); ok && len(bytes) > 0 {
			if DebugBroadcast {
				h.logger.Info(fmt.Sprintf("[broadcast] user >>> render >>> user %v", user.Id()))
			}
			user.Send(bytes)
		} else if DebugBroadcast {
			h.logger.Info(fmt.Sprintf("[broadcast] user >>> SKIPPED >>> user %v", user.Id()))
		}
	}
}

// //////////////////////////////////////////////////
// broadcast players

func (h *hub[PlayerT]) BroadcastToPlayer(name string, id model.PlayerId, data model.Data) {
	h.BroadcastToPlayerFn(name, id, func(player PlayerT) (bool, model.Data) {
		return h.WrapPlayerData(data, player)
	})
}

func (h *hub[PlayerT]) BroadcastToPlayerFn(name string, id model.PlayerId, acceptFn func(player PlayerT) (bool, model.Data)) {
	if DebugBroadcast {
		h.logger.Info(fmt.Sprintf("[broadcast] player <<< tpl %v - player %v", name, id))
	}
	h.broadcastPlayer <- h.NewNamedPlayerTemplate(
		name,
		h.AcceptPlayerFn(id, acceptFn),
	)
}

func (h *hub[PlayerT]) BroadcastToPlayerRender(id model.PlayerId, data model.Data, renderFn func(w io.Writer, data model.Data)) {
	h.BroadcastToPlayerRenderFn(id, func(player PlayerT) (bool, model.Data) {
		return h.WrapPlayerData(data, player)
	}, renderFn)
}

func (h *hub[PlayerT]) BroadcastToPlayerRenderFn(id model.PlayerId, acceptFn func(player PlayerT) (bool, model.Data), renderFn func(w io.Writer, data model.Data)) {
	if DebugBroadcast {
		h.logger.Info(fmt.Sprintf("[broadcast] player <<< renderFn - player %v", id))
	}
	h.broadcastPlayer <- h.NewTplPlayerRenderer(
		h.AcceptPlayerFn(id, acceptFn),
		renderFn,
	)
}

func (h *hub[PlayerT]) AcceptPlayerFn(playerId model.PlayerId, acceptFn func(player PlayerT) (bool, model.Data)) func(player PlayerT) (bool, model.Data) {
	return func(player PlayerT) (bool, model.Data) {
		if player.Id() == playerId {
			if acceptFn != nil {
				return acceptFn(player)
			}
			return true, nil
		}
		return false, nil
	}
}

func (h *hub[PlayerT]) BroadcastToPlayers(name string, data model.Data) {
	h.BroadcastToPlayersFn(name, func(player PlayerT) (bool, model.Data) {
		return h.WrapPlayerData(data, player)
	})
}

func (h *hub[PlayerT]) BroadcastToPlayersFn(name string, acceptFn func(player PlayerT) (bool, model.Data)) {
	if DebugBroadcast {
		h.logger.Info(fmt.Sprintf("[broadcast] player <<< tpl %v - players", name))
	}
	h.broadcastPlayer <- h.NewNamedPlayerTemplate(
		name,
		acceptFn,
	)
}

func (h *hub[PlayerT]) BroadcastPlayersRender(data model.Data, renderFn func(w io.Writer, data model.Data)) {
	h.BroadcastPlayersRenderFn(func(player PlayerT) (bool, model.Data) {
		return h.WrapPlayerData(data, player)
	}, renderFn)
}

func (h *hub[PlayerT]) BroadcastPlayersRenderFn(acceptFn func(player PlayerT) (bool, model.Data), renderFn func(w io.Writer, data model.Data)) {
	if DebugBroadcast {
		h.logger.Info("[broadcast] player <<< renderFn - players")
	}
	h.broadcastPlayer <- h.NewTplPlayerRenderer(
		acceptFn,
		renderFn,
	)
}

func (h *hub[PlayerT]) BroadcastToGamePlayers(name string, gameId model.GameId, data model.Data) {
	h.BroadcastToGamePlayersFn(name, gameId, func(player PlayerT) (bool, model.Data) {
		return h.WrapPlayerData(data, player)
	})
}

func (h *hub[PlayerT]) BroadcastToGamePlayersFn(name string, gameId model.GameId, acceptFn func(player PlayerT) (bool, model.Data)) {
	if DebugBroadcast {
		h.logger.Info(fmt.Sprintf("[broadcast] player <<< tpl %v - game %v players", name, gameId))
	}
	h.broadcastPlayer <- h.NewNamedPlayerTemplate(
		name,
		h.AcceptGamePlayersFn(gameId, acceptFn),
	)
}

func (h *hub[PlayerT]) AcceptGamePlayersFn(gameId model.GameId, acceptFn func(player PlayerT) (bool, model.Data)) func(player PlayerT) (bool, model.Data) {
	return func(player PlayerT) (bool, model.Data) {
		if player.GameId() == gameId {
			if acceptFn != nil {
				return acceptFn(player)
			}
			return true, nil
		}
		return false, nil
	}
}

func (h *hub[PlayerT]) NewNamedPlayerTemplate(name string, acceptFn func(player PlayerT) (bool, model.Data)) TplRenderer[PlayerT] {
	return h.NewTplPlayerRenderer(acceptFn, h.NewNamedRenderFn(name))
}

func (h *hub[PlayerT]) NewNamedRenderFn(name string) func(w io.Writer, data model.Data) {
	return func(w io.Writer, data model.Data) {
		h.Render(w, name, data)
	}
}

func (h *hub[PlayerT]) NewTplPlayerRenderer(acceptFn func(player PlayerT) (bool, model.Data), renderFn func(w io.Writer, data model.Data)) TplRenderer[PlayerT] {
	return NewTplRenderer[PlayerT](acceptFn, renderFn)
}

func (h *hub[PlayerT]) WrapPlayerData(data model.Data, player PlayerT) (bool, model.Data) {
	data = data.With("Player", player)
	if h.wrapPlayerDataFn != nil {
		return h.wrapPlayerDataFn(data, player)
	}
	return true, data
}

func (h *hub[PlayerT]) onBroadcastPlayer(tpl TplRenderer[PlayerT]) {
	unlock := h.rlock("onBroadcastPlayer")
	defer unlock()

	for _, user := range h.users {
		if !user.HasGameId() {
			if DebugBroadcast {
				h.logger.Info(fmt.Sprintf("[broadcast] player >>> SKIPPED >>> user %v", user.Id()))
			}
			continue
		}
		playerId := user.PlayerId()
		player, err := h.getPlayerFn(playerId)
		if err != nil {
			continue
		}
		if bytes, ok := tpl.Render(player); ok && len(bytes) > 0 {
			if DebugBroadcast {
				h.logger.Info(fmt.Sprintf("[broadcast] player >>> render >>> player %v", playerId))
			}
			user.Send(bytes)
		} else if DebugBroadcast {
			h.logger.Info(fmt.Sprintf("[broadcast] player >>> SKIPPED >>> player %v", playerId))
		}

	}
}

// //////////////////////////////////////////////////
// lock

func (h *hub[PlayerT]) rlock(requester string) func() {
	logger := h.logger.WithOptions(zap.AddCallerSkip(1))
	if DebugLock {
		logger.Info(fmt.Sprintf(" >>> R-LOCK >>> hub >>> %s ", requester))
	}
	h.mutex.RLock()
	return func() {
		h.mutex.RUnlock()
		if DebugLock {
			logger.Info(fmt.Sprintf(" <<< R-LOCK <<< hub <<< %s ", requester))
		}
	}
}

func (h *hub[PlayerT]) lock(requester string) func() {
	logger := h.logger.WithOptions(zap.AddCallerSkip(1))
	if DebugLock {
		logger.Info(fmt.Sprintf(" >>> W-LOCK >>> hub >>> %s ", requester))
	}
	h.mutex.Lock()
	return func() {
		h.mutex.Unlock()
		if DebugLock {
			logger.Info(fmt.Sprintf(" <<< W-LOCK <<< hub <<< %s ", requester))
		}
	}
}
