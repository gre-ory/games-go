package model

import "github.com/gre-ory/games-go/internal/util"

type PlayerId string

func NewPlayer(name string) *Player {
	return &Player{
		Id:   PlayerId(util.GenerateId()),
		Name: name,
	}
}

type Player struct {
	Id   PlayerId
	Name string
}
