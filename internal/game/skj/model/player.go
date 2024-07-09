package model

import (
	share_websocket "github.com/gre-ory/games-go/internal/game/share/websocket"
)

func NewPlayer(player share_websocket.Player) *Player {
	return &Player{
		Player: player,
	}
}

type Player struct {
	share_websocket.Player
}

func (p *Player) Playing() bool {
	return p.Status().IsPlaying()
}
