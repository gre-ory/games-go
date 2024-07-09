package websocket

import (
	"fmt"
	"io"
	"sync"

	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/util"
	"github.com/gre-ory/games-go/internal/util/dict"
	"github.com/gre-ory/games-go/internal/util/list"

	"github.com/gre-ory/games-go/internal/game/share/model"
)

const (
	DebugLock    = false
	DebugPing    = false
	DebugMessage = false
)

type Hub[PlayerT Player] interface {
	GetPlayer(id model.PlayerId) (PlayerT, error)
	RegisterPlayer(player PlayerT)
	UnregisterPlayer(id model.PlayerId)
	UpdatePlayer(player PlayerT)
	BroadcastToAll(name string, data Data)
	BroadcastToAllFn(name string, acceptFn func(player PlayerT) (bool, any))
	BroadcastToNotPlayingPlayers(name string, data Data)
	BroadcastToNotPlayingPlayersFn(name string, acceptFn func(player PlayerT) (bool, any))
	BroadcastToGamePlayers(name string, gameId model.GameId, data Data)
	BroadcastToGamePlayersFn(name string, gameId model.GameId, acceptFn func(player PlayerT) (bool, any))
	BroadcastToPlayer(name string, id model.PlayerId, data Data)
	BroadcastToPlayerFn(name string, id model.PlayerId, acceptFn func(player PlayerT) (bool, any))
	BroadcastToPlayerRender(id model.PlayerId, data Data, renderFn func(w io.Writer, data any))
	BroadcastToPlayerRenderFn(id model.PlayerId, acceptFn func(player PlayerT) (bool, any), renderFn func(w io.Writer, data any))
	BroadcastRender(data Data, renderFn func(w io.Writer, data any))
	BroadcastRenderFn(acceptFn func(player PlayerT) (bool, any), renderFn func(w io.Writer, data any))
	WrapPlayerData(data Data, player PlayerT) (bool, any)
	GetAllPlayers() []PlayerT
	GetNotPlayingPlayers() []PlayerT
	GetGamePlayers(gameId model.GameId) []PlayerT
	FilterPlayers(filterFn func(player PlayerT) bool) []PlayerT
}

func NewHub[PlayerT Player](logger *zap.Logger, wrapData func(data Data, player PlayerT) (bool, any), tplRenderer util.TplRenderer) Hub[PlayerT] {
	h := &hub[PlayerT]{
		TplRenderer: tplRenderer,
		broadcast:   make(chan TplRenderer[PlayerT]),
		Register:    make(chan PlayerT),
		Unregister:  make(chan model.PlayerId),
		logger:      logger,
		players:     make(map[model.PlayerId]PlayerT),
		wrapData:    wrapData,
	}
	go h.run()
	return h
}

type hub[PlayerT Player] struct {
	util.TplRenderer

	broadcast  chan TplRenderer[PlayerT]
	Register   chan PlayerT
	Unregister chan model.PlayerId

	logger   *zap.Logger
	players  map[model.PlayerId]PlayerT
	mutex    sync.RWMutex
	wrapData func(data Data, player PlayerT) (bool, any)
}

// //////////////////////////////////////////////////
// run

func (h *hub[PlayerT]) run() {
	for {
		select {
		case player := <-h.Register:
			h.onRegisterPlayer(player)
		case playerId := <-h.Unregister:
			h.onUnregisterPlayer(playerId)
		case tpl := <-h.broadcast:
			h.onBroadcast(tpl)
		}
	}
}

// //////////////////////////////////////////////////
// get

func (h *hub[PlayerT]) GetPlayer(id model.PlayerId) (PlayerT, error) {
	if DebugLock {
		h.logger.Info("[api] GetPlayer.RLock...")
	}
	h.mutex.RLock()
	defer func() {
		h.mutex.RUnlock()
		if DebugLock {
			h.logger.Info("[api] ...GetPlayer.RUnlock")
		}
	}()

	player, found := h.players[id]
	if found {
		h.logger.Info(fmt.Sprintf("[get-player] player %v", id))
		return player, nil
	}
	var empty PlayerT
	return empty, ErrPlayerNotFound
}

// //////////////////////////////////////////////////
// register

func (h *hub[PlayerT]) RegisterPlayer(player PlayerT) {
	h.Register <- player
}

func (h *hub[PlayerT]) onRegisterPlayer(player PlayerT) {
	if DebugLock {
		h.logger.Info("[api] onRegisterPlayer.Lock...")
	}
	h.mutex.Lock()
	defer func() {
		h.mutex.Unlock()
		if DebugLock {
			h.logger.Info("[api] ...onRegisterPlayer.Unlock")
		}
	}()

	h.logger.Info(fmt.Sprintf("[register] (+) player %v", player.Id()))
	h.players[player.Id()] = player
}

// //////////////////////////////////////////////////
// unregister

func (h *hub[PlayerT]) UnregisterPlayer(id model.PlayerId) {
	h.Unregister <- id
}

func (h *hub[PlayerT]) onUnregisterPlayer(id model.PlayerId) {
	if DebugLock {
		h.logger.Info("[api] onUnregisterPlayer.Lock...")
	}
	h.mutex.Lock()
	defer func() {
		h.mutex.Unlock()
		if DebugLock {
			h.logger.Info("[api] ...onUnregisterPlayer.Unlock")
		}
	}()

	h.logger.Info(fmt.Sprintf("[unregister] (-) player %v", id))
	// if player, ok := h.players[id]; ok {
	if _, ok := h.players[id]; ok {
		delete(h.players, id)
		// player.Close(h.logger)
	}
}

func (h *hub[PlayerT]) UpdatePlayer(player PlayerT) {
	if DebugLock {
		h.logger.Info("[api] UpdatePlayer.Lock...")
	}
	h.mutex.Lock()
	defer func() {
		h.mutex.Unlock()
		if DebugLock {
			h.logger.Info("[api] ...UpdatePlayer.Unlock")
		}
	}()

	h.logger.Info(fmt.Sprintf("[update] (~) player %v", player.Id()))
	h.players[player.Id()] = player
}

// //////////////////////////////////////////////////
// broadcast

func (h *hub[PlayerT]) BroadcastToNotPlayingPlayers(name string, data Data) {
	h.BroadcastToNotPlayingPlayersFn(name, func(player PlayerT) (bool, any) {
		return h.WrapPlayerData(data, player)
	})
}

func (h *hub[PlayerT]) BroadcastToAllFn(name string, acceptFn func(player PlayerT) (bool, any)) {
	h.broadcast <- h.NewNamedTemplate(
		name,
		acceptFn,
	)
}

func (h *hub[PlayerT]) BroadcastToAll(name string, data Data) {
	h.BroadcastToAllFn(name, func(player PlayerT) (bool, any) {
		return h.WrapPlayerData(data, player)
	})
}

func (h *hub[PlayerT]) BroadcastToNotPlayingPlayersFn(name string, acceptFn func(player PlayerT) (bool, any)) {
	h.broadcast <- h.NewNamedTemplate(
		name,
		h.AcceptNotPlayingPlayersFn(acceptFn),
	)
}

func (h *hub[PlayerT]) AcceptNotPlayingPlayersFn(acceptFn func(player PlayerT) (bool, any)) func(player PlayerT) (bool, any) {
	return func(player PlayerT) (bool, any) {
		if player.CanJoin() {
			if acceptFn != nil {
				return acceptFn(player)
			}
			return true, nil
		}
		return false, nil
	}
}

func (h *hub[PlayerT]) BroadcastToGamePlayers(name string, gameId model.GameId, data Data) {
	h.BroadcastToGamePlayersFn(name, gameId, func(player PlayerT) (bool, any) {
		return h.WrapPlayerData(data, player)
	})
}

func (h *hub[PlayerT]) BroadcastToGamePlayersFn(name string, gameId model.GameId, acceptFn func(player PlayerT) (bool, any)) {
	h.broadcast <- h.NewNamedTemplate(
		name,
		h.AcceptGamePlayersFn(gameId, acceptFn),
	)
}

func (h *hub[PlayerT]) AcceptGamePlayersFn(gameId model.GameId, acceptFn func(player PlayerT) (bool, any)) func(player PlayerT) (bool, any) {
	return func(player PlayerT) (bool, any) {
		if player.GameId() == gameId {
			if acceptFn != nil {
				return acceptFn(player)
			}
			return true, nil
		}
		return false, nil
	}
}

func (h *hub[PlayerT]) BroadcastToPlayer(name string, id model.PlayerId, data Data) {
	h.BroadcastToPlayerFn(name, id, func(player PlayerT) (bool, any) {
		return h.WrapPlayerData(data, player)
	})
}

func (h *hub[PlayerT]) BroadcastToPlayerFn(name string, id model.PlayerId, acceptFn func(player PlayerT) (bool, any)) {
	h.broadcast <- h.NewNamedTemplate(
		name,
		h.AcceptPlayerFn(id, acceptFn),
	)
}

func (h *hub[PlayerT]) AcceptPlayerFn(id model.PlayerId, acceptFn func(player PlayerT) (bool, any)) func(player PlayerT) (bool, any) {
	return func(player PlayerT) (bool, any) {
		if player.Id() == id {
			if acceptFn != nil {
				return acceptFn(player)
			}
			return true, nil
		}
		return false, nil
	}
}

func (h *hub[PlayerT]) NewNamedTemplate(name string, acceptFn func(player PlayerT) (bool, any)) TplRenderer[PlayerT] {
	return h.NewTplRenderer(acceptFn, h.NewNamedRenderFn(name))
}

func (h *hub[PlayerT]) NewNamedRenderFn(name string) func(w io.Writer, data any) {
	return func(w io.Writer, data any) {
		h.Render(w, name, data)
	}
}

func (h *hub[PlayerT]) BroadcastToPlayerRender(id model.PlayerId, data Data, renderFn func(w io.Writer, data any)) {
	h.BroadcastToPlayerRenderFn(id, func(player PlayerT) (bool, any) {
		return h.WrapPlayerData(data, player)
	}, renderFn)
}

func (h *hub[PlayerT]) BroadcastToPlayerRenderFn(id model.PlayerId, acceptFn func(player PlayerT) (bool, any), renderFn func(w io.Writer, data any)) {
	h.broadcast <- h.NewTplRenderer(
		h.AcceptPlayerFn(id, acceptFn),
		renderFn,
	)
}

func (h *hub[PlayerT]) BroadcastRender(data Data, renderFn func(w io.Writer, data any)) {
	h.BroadcastRenderFn(func(player PlayerT) (bool, any) {
		return h.WrapPlayerData(data, player)
	}, renderFn)
}

func (h *hub[PlayerT]) BroadcastRenderFn(acceptFn func(player PlayerT) (bool, any), renderFn func(w io.Writer, data any)) {
	h.broadcast <- h.NewTplRenderer(
		acceptFn,
		renderFn,
	)
}

func (h *hub[PlayerT]) NewTplRenderer(acceptFn func(player PlayerT) (bool, any), renderFn func(w io.Writer, data any)) TplRenderer[PlayerT] {
	return NewTplRenderer[PlayerT](acceptFn, renderFn)
}

func (h *hub[PlayerT]) WrapPlayerData(data Data, player PlayerT) (bool, any) {
	data.With("player", player)
	if player.HasGameId() {
		data.With("game_id", player.GameId())
	}
	if h.wrapData != nil {
		return h.wrapData(data, player)
	}
	return true, data
}

func (h *hub[PlayerT]) onBroadcast(tpl TplRenderer[PlayerT]) {
	if DebugLock {
		h.logger.Info("[api] onBroadcast.RLock...")
	}
	h.mutex.RLock()
	defer func() {
		h.mutex.RUnlock()
		if DebugLock {
			h.logger.Info("[api] ...onBroadcast.RUnlock")
		}
	}()

	for _, player := range h.players {
		if bytes, ok := tpl.Render(player); ok && len(bytes) > 0 {
			player.Send(bytes)
		}
	}
}

// //////////////////////////////////////////////////
// helpers

func (h *hub[PlayerT]) GetAllPlayers() []PlayerT {
	if DebugLock {
		h.logger.Info("[api] GetAllPlayers.RLock...")
	}
	h.mutex.RLock()
	defer func() {
		h.mutex.RUnlock()
		if DebugLock {
			h.logger.Info("[api] ...GetAllPlayers.RUnlock")
		}
	}()

	return dict.ConvertToList(h.players, dict.Value)
}

func (h *hub[PlayerT]) GetNotPlayingPlayers() []PlayerT {
	return h.FilterPlayers(func(player PlayerT) bool {
		return player.CanJoin()
	})
}

func (h *hub[PlayerT]) GetGamePlayers(gameId model.GameId) []PlayerT {
	return h.FilterPlayers(func(player PlayerT) bool {
		return player.GameId() == gameId
	})
}

func (h *hub[PlayerT]) FilterPlayers(filterFn func(player PlayerT) bool) []PlayerT {
	return list.Filter(h.GetAllPlayers(), filterFn)
}
