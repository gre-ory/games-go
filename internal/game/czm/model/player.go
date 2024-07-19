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
	Cards []Card
}

func (p *Player) WithCard(card Card) *Player {
	p.Cards = append(p.Cards, card)
	return p
}

func (p *Player) SelectCard(cardNumber int) (Card, error) {
	if cardNumber < 1 || cardNumber > len(p.Cards) {
		return p.Cards[cardNumber-1], nil
	}
	return 0, ErrInvalidCardIndex
}

func (p *Player) PlayCard(cardNumber int) (Card, error) {
	card, err := p.SelectCard(cardNumber)
	if err != nil {
		return 0, err
	}
	p.Cards = append(p.Cards[:cardNumber-1], p.Cards[cardNumber-1:]...)
	return card, nil
}
