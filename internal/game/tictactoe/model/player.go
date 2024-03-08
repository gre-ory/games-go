package model

import (
	"strings"

	"github.com/gre-ory/games-go/internal/util"
	"github.com/gre-ory/games-go/internal/util/websocket"
)

type PlayerId string

func NewPlayerId() PlayerId {
	return PlayerId(util.GeneratePlayerId())
}

func NewPlayer(player websocket.Player[PlayerId, GameId], name string) *Player {
	return &Player{
		Player: player,
		Name:   name,
		Status: WaitingToStart,
	}
}

type PlayerStatus int

const (
	WaitingToStart PlayerStatus = iota
	WaitingToPlay
	Playing
	Win
	Tie
	Loose
)

type Player struct {
	websocket.Player[PlayerId, GameId]
	Name   string
	Symbol rune
	Status PlayerStatus
}

func (p *Player) WithSymbol(symbol rune) *Player {
	p.Symbol = symbol
	return p
}

func (p *Player) Playing() bool {
	return p.Status == Playing
}

func (p *Player) Labels() string {
	labels := make([]string, 0)
	labels = append(labels, "player")
	switch p.Symbol {
	case PLAYER_ONE_SYMBOL:
		labels = append(labels, "symbol-1")
	case PLAYER_TWO_SYMBOL:
		labels = append(labels, "symbol-2")
	}
	switch p.Status {
	case WaitingToStart:
		labels = append(labels, "waiting-to-start")
	case WaitingToPlay:
		labels = append(labels, "waiting-to-play")
	case Playing:
		labels = append(labels, "playing")
	case Win:
		labels = append(labels, "win")
	case Tie:
		labels = append(labels, "tie")
	case Loose:
		labels = append(labels, "loose")
	}
	return strings.Join(labels, " ")
}
