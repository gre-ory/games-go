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
	Status   PlayerStatus
	Cards    []Card
}

func (p *Player) GetLanguage() string {
	return p.Language
}

func (p *Player) WithCard(card Card) *Player {
	p.Cards = append(p.Cards, card)
	return p
}

func (p *Player) SelectCard(cardIndex int) (Card, error) {
	if cardIndex < 0 || cardIndex >= len(p.Cards) {
		return p.Cards[cardIndex], nil
	}
	return 0, ErrInvalidCardIndex
}

func (p *Player) PlayCard(cardIndex int) (Card, error) {
	card, err := p.SelectCard(cardIndex)
	if err != nil {
		return 0, err
	}
	p.Cards = append(p.Cards[:cardIndex], p.Cards[cardIndex+1:]...)
	return card, nil
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

// func (p *Player) CardsHtml() template.HTML {
// 	return template.HTML(fmt.Sprintf("<div class=\"cards\">%s</div>", strings.Join(list.Convert(p.Cards, Card.CardStr), "")))
// }

func (p *Player) Labels() string {
	labels := make([]string, 0)
	labels = append(labels, "player")
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
			return localizer.Loc("YouPlaying")
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
			return localizer.Loc("PlayerPlaying")
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
