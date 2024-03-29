package model

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/gre-ory/games-go/internal/util"
	"github.com/gre-ory/games-go/internal/util/websocket"
)

type PlayerId string

func NewPlayerId() PlayerId {
	return PlayerId(util.GeneratePlayerId())
}

func NewPlayer(player websocket.Player[PlayerId, GameId], avatar int, name string) *Player {
	return &Player{
		Player: player,
		Avatar: avatar,
		Name:   name,
		Status: WaitingToJoin,
	}
}

type PlayerStatus int

const (
	WaitingToJoin PlayerStatus = iota
	WaitingToStart
	WaitingToPlay
	Playing
	Win
	Tie
	Loose
)

type Player struct {
	websocket.Player[PlayerId, GameId]
	Avatar int
	Name   string
	Symbol rune
	Status PlayerStatus
}

func (p *Player) WithSymbol(symbol rune) *Player {
	p.Symbol = symbol
	return p
}

func (p *Player) CanJoin() bool {
	return p.Player.CanJoin() && p.Status == WaitingToJoin && p.Name != ""
}

func (p *Player) Playing() bool {
	return p.Status == Playing
}

func (p *Player) ExtraSmallAvatarHtml() template.HTML {
	if p.Avatar != 0 {
		return template.HTML(fmt.Sprintf("<div class=\"avatar-%d xs\"></div>", p.Avatar))
	}
	return ""
}

func (p *Player) SmallAvatarHtml() template.HTML {
	if p.Avatar != 0 {
		return template.HTML(fmt.Sprintf("<div class=\"avatar-%d s\"></div>", p.Avatar))
	}
	return ""
}

func (p *Player) AvatarHtml() template.HTML {
	if p.Avatar != 0 {
		return template.HTML(fmt.Sprintf("<div class=\"avatar-%d m\"></div>", p.Avatar))
	}
	return ""
}

func (p *Player) IconHtml() template.HTML {
	switch p.Symbol {
	case PLAYER_ONE_SYMBOL:
		return "<div class=\"icon-cross\"></div>"
	case PLAYER_TWO_SYMBOL:
		return "<div class=\"icon-circle\"></div>"
	}
	return ""
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
	if p.Active() {
		switch p.Status {
		case WaitingToJoin:
			labels = append(labels, "waiting-to-start")
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
	} else {
		labels = append(labels, "disconnected")
	}
	return strings.Join(labels, " ")
}

func (p *Player) YourMessage() template.HTML {
	if p.Active() {
		switch p.Status {
		case WaitingToJoin:
			return "Wait other player!"
		case WaitingToStart:
			return "Start?"
		case WaitingToPlay:
			return "Wait!"
		case Playing:
			return "Play " + p.IconHtml() + "!"
		case Win:
			return "You wins!"
		case Tie:
			return "Tie!"
		case Loose:
			return "You looses!"
		}
	} else {
		return "Disconnected..."
	}
	return ""
}

func (p *Player) Message() template.HTML {
	if p.Active() {
		switch p.Status {
		case WaitingToJoin:
			return "Waiting..."
		case WaitingToStart:
			return "Start?"
		case WaitingToPlay:
			return "Waiting..."
		case Playing:
			return "Playing " + p.IconHtml() + "..."
		case Win:
			return "Wins!"
		case Tie:
			return "Tie!"
		case Loose:
			return "Looses!"
		}
	} else {
		return "Disconnected..."
	}
	return ""
}

func (p *Player) SymbolIcon() string {
	switch p.Symbol {
	case PLAYER_ONE_SYMBOL:
		return "icon-cross"
	case PLAYER_TWO_SYMBOL:
		return "icon-circle"
	}
	return ""
}

func (p *Player) StatusIcon() string {
	switch p.Status {
	case WaitingToJoin:
		return "icon-pause"
	case WaitingToPlay:
		return "icon-pause"
	case Playing:
		return "icon-play"
	case Win:
		return "icon-win"
	case Tie:
		return "icon-tie"
	case Loose:
		return "icon-loose"
	}
	return ""
}
