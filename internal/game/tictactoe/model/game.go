package model

import (
	"github.com/gre-ory/games-go/internal/util"
)

type GameId string

func NewGame() *Game {
	return &Game{
		Id:      GameId(util.GenerateId()),
		Players: make(map[PlayerId]*Player),
	}
}

type Game struct {
	Id      GameId
	Players map[PlayerId]*Player
	Started bool
}

func (g *Game) WithPlayer(player *Player) *Game {
	g.Players[player.Id] = player
	return g
}
