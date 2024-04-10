package websocket

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/gre-ory/games-go/internal/util"
	"github.com/gre-ory/games-go/internal/util/dict"
	"github.com/gre-ory/games-go/internal/util/list"
	"go.uber.org/zap"
)

type Hub[IdT comparable, GameIdT comparable, PlayerT Player[IdT, GameIdT]] interface {
	GetPlayer(id IdT) (PlayerT, error)
	RegisterPlayer(player PlayerT)
	UnregisterPlayer(id IdT)
	UpdatePlayer(player PlayerT)
	BroadcastToAllFn(name string, acceptFn func(player PlayerT) (bool, any))
	BroadcastToAll(name string, data Data)
	BroadcastToNotPlayingPlayersFn(name string, acceptFn func(player PlayerT) (bool, any))
	BroadcastToNotPlayingPlayers(name string, data Data)
	BroadcastToGamePlayersFn(name string, gameId GameIdT, acceptFn func(player PlayerT) (bool, any))
	BroadcastToGamePlayers(name string, gameId GameIdT, data Data)
	BroadcastToPlayerFn(name string, id IdT, acceptFn func(player PlayerT) (bool, any))
	BroadcastToPlayer(name string, id IdT, data Data)
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
	h.logger.Info("[api] GetPlayer.Lock...")
	h.mutex.RLock()
	defer func() {
		h.mutex.RUnlock()
		h.logger.Info("[api] ...GetPlayer.Unlock")
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
	h.logger.Info("[api] onRegisterPlayer.Lock...")
	h.mutex.Lock()
	defer func() {
		h.mutex.Unlock()
		h.logger.Info("[api] ...onRegisterPlayer.Unlock")
	}()

	h.logger.Info(fmt.Sprintf("[register] (+) player %s", player.Id()))
	h.players[player.Id()] = player
}

// //////////////////////////////////////////////////
// unregister

func (h *hub[IdT, GameIdT, PlayerT]) UnregisterPlayer(id IdT) {
	h.Unregister <- id
}

func (h *hub[IdT, GameIdT, PlayerT]) onUnregisterPlayer(id IdT) {
	h.logger.Info("[api] onUnregisterPlayer.Lock...")
	h.mutex.Lock()
	defer func() {
		h.mutex.Unlock()
		h.logger.Info("[api] ...onUnregisterPlayer.Unlock")
	}()

	h.logger.Info(fmt.Sprintf("[unregister] (-) player %s", id))
	if player, ok := h.players[id]; ok {
		delete(h.players, id)
		player.Close()
	}
}

func (h *hub[IdT, GameIdT, PlayerT]) UpdatePlayer(player PlayerT) {
	h.logger.Info("[api] UpdatePlayer.Lock...")
	h.mutex.Lock()
	defer func() {
		h.mutex.Unlock()
		h.logger.Info("[api] ...UpdatePlayer.Unlock")
	}()

	h.logger.Info(fmt.Sprintf("[update] (~) player %s", player.Id()))
	h.players[player.Id()] = player
}

// //////////////////////////////////////////////////
// broadcast

func (h *hub[IdT, GameIdT, PlayerT]) BroadcastToAllFn(name string, acceptFn func(player PlayerT) (bool, any)) {
	h.broadcast <- NewTemplate[PlayerT](
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
	h.broadcast <- NewTemplate[PlayerT](
		name,
		func(player PlayerT) (bool, any) {
			h.logger.Info(
				fmt.Sprintf(
					"[broadcast] %s -> player %s -> can-join: %t, active: %t, id: %t, game: %t",
					name,
					player.Id(),
					player.CanJoin(),
					player.Active(),
					player.HasId(),
					player.HasGameId(),
				),
			)
			if player.CanJoin() {
				return acceptFn(player)
			}
			return false, nil
		},
	)
}

func (h *hub[IdT, GameIdT, PlayerT]) BroadcastToNotPlayingPlayers(name string, data Data) {
	h.BroadcastToNotPlayingPlayersFn(name, func(player PlayerT) (bool, any) {
		return h.WrapPlayerData(data, player)
	})
}

func (h *hub[IdT, GameIdT, PlayerT]) BroadcastToGamePlayersFn(name string, gameId GameIdT, acceptFn func(player PlayerT) (bool, any)) {
	h.broadcast <- NewTemplate[PlayerT](
		name,
		func(player PlayerT) (bool, any) {
			if player.GameId() == gameId {
				return acceptFn(player)
			}
			return false, nil
		},
	)
}

func (h *hub[IdT, GameIdT, PlayerT]) BroadcastToGamePlayers(name string, gameId GameIdT, data Data) {
	h.BroadcastToGamePlayersFn(name, gameId, func(player PlayerT) (bool, any) {
		return h.WrapPlayerData(data, player)
	})
}

func (h *hub[IdT, GameIdT, PlayerT]) BroadcastToPlayerFn(name string, id IdT, acceptFn func(player PlayerT) (bool, any)) {
	h.broadcast <- NewTemplate[PlayerT](
		name,
		func(player PlayerT) (bool, any) {
			if player.Id() == id {
				h.logger.Info(fmt.Sprintf("[  ] accept: %s, %v == %v", name, player.Id(), id))
				return acceptFn(player)
			}
			h.logger.Info(fmt.Sprintf("[KO] accept: %s, %v == %v", name, player.Id(), id))
			return false, nil
		},
	)
}

func (h *hub[IdT, GameIdT, PlayerT]) BroadcastToPlayer(name string, id IdT, data Data) {
	h.BroadcastToPlayerFn(name, id, func(player PlayerT) (bool, any) {
		return h.WrapPlayerData(data, player)
	})
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
	h.logger.Info("[api] onBroadcast.Lock...")
	h.mutex.RLock()
	defer func() {
		h.mutex.RUnlock()
		h.logger.Info("[api] ...onBroadcast.Unlock")
	}()

	count := 0
	for _, player := range h.players {
		if ok, data := tpl.Accept(player); ok {
			buf := &bytes.Buffer{}
			h.Render(buf, tpl.GetName(), data)
			player.Send(buf.Bytes())
			count++
		}
	}
	if count > 0 {
		h.logger.Info(fmt.Sprintf("[broadcast] template: %s >>> %d player(s)", tpl.GetName(), count))
	} else {
		h.logger.Info(fmt.Sprintf("[broadcast] template: %s >>> SKIP", tpl.GetName()))
	}
}

// //////////////////////////////////////////////////
// helpers

func (h *hub[IdT, GameIdT, PlayerT]) GetAllPlayers() []PlayerT {
	h.logger.Info("[api] GetAllPlayers.Lock...")
	h.mutex.RLock()
	defer func() {
		h.mutex.RUnlock()
		h.logger.Info("[api] ...GetAllPlayers.Unlock")
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
