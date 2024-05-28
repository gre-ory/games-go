package websocket

import (
	"fmt"
	"io"
	"sync"

	"github.com/gre-ory/games-go/internal/util"
	"github.com/gre-ory/games-go/internal/util/dict"
	"github.com/gre-ory/games-go/internal/util/list"
	"go.uber.org/zap"
)

const (
	DebugLock    = false
	DebugPing    = false
	DebugMessage = false
)

type Hub[IdT comparable, GameIdT comparable, PlayerT Player[IdT, GameIdT]] interface {
	GetPlayer(id IdT) (PlayerT, error)
	RegisterPlayer(player PlayerT)
	UnregisterPlayer(id IdT)
	UpdatePlayer(player PlayerT)
	BroadcastToAll(name string, data Data)
	BroadcastToAllFn(name string, acceptFn func(player PlayerT) (bool, any))
	BroadcastToNotPlayingPlayers(name string, data Data)
	BroadcastToNotPlayingPlayersFn(name string, acceptFn func(player PlayerT) (bool, any))
	BroadcastToGamePlayers(name string, gameId GameIdT, data Data)
	BroadcastToGamePlayersFn(name string, gameId GameIdT, acceptFn func(player PlayerT) (bool, any))
	BroadcastToPlayer(name string, id IdT, data Data)
	BroadcastToPlayerFn(name string, id IdT, acceptFn func(player PlayerT) (bool, any))
	BroadcastToPlayerRender(id IdT, data Data, renderFn func(w io.Writer, data any))
	BroadcastToPlayerRenderFn(id IdT, acceptFn func(player PlayerT) (bool, any), renderFn func(w io.Writer, data any))
	BroadcastRender(data Data, renderFn func(w io.Writer, data any))
	BroadcastRenderFn(acceptFn func(player PlayerT) (bool, any), renderFn func(w io.Writer, data any))
	WrapPlayerData(data Data, player PlayerT) (bool, any)
	GetAllPlayers() []PlayerT
	GetNotPlayingPlayers() []PlayerT
	GetGamePlayers(gameId GameIdT) []PlayerT
}

func NewHub[IdT comparable, GameIdT comparable, PlayerT Player[IdT, GameIdT]](logger *zap.Logger, wrapData func(data Data, player PlayerT) (bool, any), tplRenderer util.TplRenderer) Hub[IdT, GameIdT, PlayerT] {
	h := &hub[IdT, GameIdT, PlayerT]{
		TplRenderer: tplRenderer,
		broadcast:   make(chan Template[PlayerT]),
		Register:    make(chan PlayerT),
		Unregister:  make(chan IdT),
		logger:      logger,
		players:     make(map[IdT]PlayerT),
		wrapData:    wrapData,
	}
	go h.run()
	return h
}

type hub[IdT comparable, GameIdT comparable, PlayerT Player[IdT, GameIdT]] struct {
	util.TplRenderer

	broadcast  chan Template[PlayerT]
	Register   chan PlayerT
	Unregister chan IdT

	logger   *zap.Logger
	players  map[IdT]PlayerT
	mutex    sync.RWMutex
	wrapData func(data Data, player PlayerT) (bool, any)
}

// //////////////////////////////////////////////////
// run

func (h *hub[IdT, GameIdT, PlayerT]) run() {
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

func (h *hub[IdT, GameIdT, PlayerT]) GetPlayer(id IdT) (PlayerT, error) {
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

func (h *hub[IdT, GameIdT, PlayerT]) RegisterPlayer(player PlayerT) {
	h.Register <- player
}

func (h *hub[IdT, GameIdT, PlayerT]) onRegisterPlayer(player PlayerT) {
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

func (h *hub[IdT, GameIdT, PlayerT]) UnregisterPlayer(id IdT) {
	h.Unregister <- id
}

func (h *hub[IdT, GameIdT, PlayerT]) onUnregisterPlayer(id IdT) {
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

func (h *hub[IdT, GameIdT, PlayerT]) UpdatePlayer(player PlayerT) {
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

func (h *hub[IdT, GameIdT, PlayerT]) BroadcastToNotPlayingPlayers(name string, data Data) {
	h.BroadcastToNotPlayingPlayersFn(name, func(player PlayerT) (bool, any) {
		return h.WrapPlayerData(data, player)
	})
}

func (h *hub[IdT, GameIdT, PlayerT]) BroadcastToAllFn(name string, acceptFn func(player PlayerT) (bool, any)) {
	h.broadcast <- h.NewNamedTemplate(
		name,
		acceptFn,
	)
}

func (h *hub[IdT, GameIdT, PlayerT]) BroadcastToAll(name string, data Data) {
	h.BroadcastToAllFn(name, func(player PlayerT) (bool, any) {
		return h.WrapPlayerData(data, player)
	})
}

func (h *hub[IdT, GameIdT, PlayerT]) BroadcastToNotPlayingPlayersFn(name string, acceptFn func(player PlayerT) (bool, any)) {
	h.broadcast <- h.NewNamedTemplate(
		name,
		h.AcceptNotPlayingPlayersFn(acceptFn),
	)
}

func (h *hub[IdT, GameIdT, PlayerT]) AcceptNotPlayingPlayersFn(acceptFn func(player PlayerT) (bool, any)) func(player PlayerT) (bool, any) {
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

func (h *hub[IdT, GameIdT, PlayerT]) BroadcastToGamePlayers(name string, gameId GameIdT, data Data) {
	h.BroadcastToGamePlayersFn(name, gameId, func(player PlayerT) (bool, any) {
		return h.WrapPlayerData(data, player)
	})
}

func (h *hub[IdT, GameIdT, PlayerT]) BroadcastToGamePlayersFn(name string, gameId GameIdT, acceptFn func(player PlayerT) (bool, any)) {
	h.broadcast <- h.NewNamedTemplate(
		name,
		h.AcceptGamePlayersFn(gameId, acceptFn),
	)
}

func (h *hub[IdT, GameIdT, PlayerT]) AcceptGamePlayersFn(gameId GameIdT, acceptFn func(player PlayerT) (bool, any)) func(player PlayerT) (bool, any) {
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

func (h *hub[IdT, GameIdT, PlayerT]) BroadcastToPlayer(name string, id IdT, data Data) {
	h.BroadcastToPlayerFn(name, id, func(player PlayerT) (bool, any) {
		return h.WrapPlayerData(data, player)
	})
}

func (h *hub[IdT, GameIdT, PlayerT]) BroadcastToPlayerFn(name string, id IdT, acceptFn func(player PlayerT) (bool, any)) {
	h.broadcast <- h.NewNamedTemplate(
		name,
		h.AcceptPlayerFn(id, acceptFn),
	)
}

func (h *hub[IdT, GameIdT, PlayerT]) AcceptPlayerFn(id IdT, acceptFn func(player PlayerT) (bool, any)) func(player PlayerT) (bool, any) {
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

func (h *hub[IdT, GameIdT, PlayerT]) NewNamedTemplate(name string, acceptFn func(player PlayerT) (bool, any)) Template[PlayerT] {
	return h.NewTemplate(acceptFn, h.NewNamedRenderFn(name))
}

func (h *hub[IdT, GameIdT, PlayerT]) NewNamedRenderFn(name string) func(w io.Writer, data any) {
	return func(w io.Writer, data any) {
		h.Render(w, name, data)
	}
}

func (h *hub[IdT, GameIdT, PlayerT]) BroadcastToPlayerRender(id IdT, data Data, renderFn func(w io.Writer, data any)) {
	h.BroadcastToPlayerRenderFn(id, func(player PlayerT) (bool, any) {
		return h.WrapPlayerData(data, player)
	}, renderFn)
}

func (h *hub[IdT, GameIdT, PlayerT]) BroadcastToPlayerRenderFn(id IdT, acceptFn func(player PlayerT) (bool, any), renderFn func(w io.Writer, data any)) {
	h.broadcast <- h.NewTemplate(
		h.AcceptPlayerFn(id, acceptFn),
		renderFn,
	)
}

func (h *hub[IdT, GameIdT, PlayerT]) BroadcastRender(data Data, renderFn func(w io.Writer, data any)) {
	h.BroadcastRenderFn(func(player PlayerT) (bool, any) {
		return h.WrapPlayerData(data, player)
	}, renderFn)
}

func (h *hub[IdT, GameIdT, PlayerT]) BroadcastRenderFn(acceptFn func(player PlayerT) (bool, any), renderFn func(w io.Writer, data any)) {
	h.broadcast <- h.NewTemplate(
		acceptFn,
		renderFn,
	)
}

func (h *hub[IdT, GameIdT, PlayerT]) NewTemplate(acceptFn func(player PlayerT) (bool, any), renderFn func(w io.Writer, data any)) Template[PlayerT] {
	return NewTemplate[PlayerT](acceptFn, renderFn)
}

func (h *hub[IdT, GameIdT, PlayerT]) WrapPlayerData(data Data, player PlayerT) (bool, any) {
	data.With("player", player)
	if player.HasGameId() {
		data.With("game_id", player.GameId())
	}
	if h.wrapData != nil {
		return h.wrapData(data, player)
	}
	return true, data
}

func (h *hub[IdT, GameIdT, PlayerT]) onBroadcast(tpl Template[PlayerT]) {
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

func (h *hub[IdT, GameIdT, PlayerT]) GetAllPlayers() []PlayerT {
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

func (h *hub[IdT, GameIdT, PlayerT]) GetNotPlayingPlayers() []PlayerT {
	return list.Filter(h.GetAllPlayers(), func(player PlayerT) bool {
		return player.CanJoin()
	})
}

func (h *hub[IdT, GameIdT, PlayerT]) GetGamePlayers(gameId GameIdT) []PlayerT {
	return list.Filter(h.GetAllPlayers(), func(player PlayerT) bool {
		return player.GameId() == gameId
	})
}
