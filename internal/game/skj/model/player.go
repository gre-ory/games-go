package model

import (
	share_model "github.com/gre-ory/games-go/internal/game/share/model"
)

func NewPlayerFromUser(gameId share_model.GameId, user share_model.User) *Player {
	return &Player{
		Player: share_model.NewPlayerFromUser(gameId, user),
		board:  NewPlayerBoard(),
	}
}

type Player struct {
	share_model.Player
	board *PlayerBoard
}

func (p *Player) Board() *PlayerBoard {
	return p.board
}
