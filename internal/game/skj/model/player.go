package model

import (
	share_model "github.com/gre-ory/games-go/internal/game/share/model"
)

func NewPlayerFromUser(gameId share_model.GameId, user share_model.User) *Player {
	return &Player{
		Player: share_model.NewPlayerFromUser(gameId, user),
	}
}

type Player struct {
	share_model.Player
}
