package model

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/gre-ory/games-go/internal/util"
	"github.com/gre-ory/games-go/internal/util/loc"
	"github.com/gre-ory/games-go/internal/util/websocket"
)

type PlayerId string

func NewPlayerId() PlayerId {
	return PlayerId(util.GeneratePlayerId())
}

func NewPlayer(player websocket.Player[PlayerId, GameId], avatar int, name string, language string) *Player {
	return &Player{
		Player:   player,
		Avatar:   avatar,
		Name:     name,
		Language: language,
		Status:   WaitingToJoin,
	}
}

type PlayerStatus int

const (
	WaitingToJoin PlayerStatus = iota
	WaitingToJoinOrStart
	WaitingToStart
	WaitingToPlay
	Playing
)

type Player struct {
	websocket.Player[PlayerId, GameId]
	Avatar   int
	Name     string
	Language string
	Symbol   rune
	Status   PlayerStatus
}

func (p *Player) GetLanguage() string {
	return p.Language
}

func (p *Player) WithSymbol(symbol rune) *Player {
	p.Symbol = symbol
	return p
}

func (p *Player) CanJoin() bool {
	return p.Player.CanJoin() && (p.Status == WaitingToJoin)
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
	if p.IsActive() {
		switch p.Status {
		case WaitingToJoin:
			labels = append(labels, "waiting-to-join")
		case WaitingToStart:
			labels = append(labels, "waiting-to-start")
		case WaitingToPlay:
			labels = append(labels, "waiting-to-play")
		case Playing:
			labels = append(labels, "playing")
		}
	} else {
		labels = append(labels, "disconnected")
	}
	return strings.Join(labels, " ")
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

func (p *Player) YourMessage(localizer loc.Localizer) template.HTML {
	if p.IsActive() {
		switch p.Status {
		case WaitingToJoin:
			return localizer.Loc("YouWaitingToJoin")
		case WaitingToStart:
			return localizer.Loc("YouWaitingToStart")
		case WaitingToPlay:
			return localizer.Loc("YouWaitingToPlay")
		case Playing:
			return localizer.Loc("YouPlaying", p.IconHtml())
		}
	} else {
		return localizer.Loc("YouDisconnected")
	}
	return ""
}

func (p *Player) Message(localizer loc.Localizer) template.HTML {
	if p.IsActive() {
		switch p.Status {
		case WaitingToJoin:
			return localizer.Loc("PlayerWaitingToJoin")
		case WaitingToStart:
			return localizer.Loc("PlayerWaitingToStart")
		case WaitingToPlay:
			return localizer.Loc("PlayerWaitingToPlay")
		case Playing:
			return localizer.Loc("PlayerPlaying", p.IconHtml())
		}
	} else {
		return localizer.Loc("PlayerDisconnected")
	}
	return ""
}

func (p *Player) StatusIcon() string {
	switch p.Status {
	case WaitingToJoin,
		WaitingToStart,
		WaitingToPlay:
		return "icon-pause"
	case Playing:
		return "icon-play"
	}
	return ""
}
