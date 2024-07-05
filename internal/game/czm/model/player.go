package model

import (
	"html/template"
	"strings"

	"github.com/gre-ory/games-go/internal/util/loc"
	"github.com/gre-ory/games-go/internal/util/websocket"
)

func NewPlayer(player websocket.Player) *Player {
	return &Player{
		Player: player,
	}
}

type Player struct {
	*websocket.Player
	Cards []Card
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
	return p.Status().IsWaitingToJoin() && p.Player.CanJoin()
}

func (p *Player) Playing() bool {
	return p.Status().IsPlaying()
}

// func (p *Player) CardsHtml() template.HTML {
// 	return template.HTML(fmt.Sprintf("<div class=\"cards\">%s</div>", strings.Join(list.Convert(p.Cards, Card.CardStr), "")))
// }

func (p *Player) Labels() string {
	labels := make([]string, 0)
	labels = append(labels, "player")
	if p.IsActive() {
		labels = append(labels, p.Status().LabelSlice()...)
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
