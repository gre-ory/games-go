package model

import (
	share_model "github.com/gre-ory/games-go/internal/game/share/model"
)

const (
	NbRow       = 3
	NbColumn    = 4
	MinNbPlayer = 2
	MaxNbPlayer = 4
)

func NewGame() *Game {
	return &Game{
		Game:        share_model.NewGame[*Player](MinNbPlayer, MaxNbPlayer),
		NbRow:       NbRow,
		NbColumn:    NbColumn,
		DrawDeck:    NewDrawCardDeck(),
		DiscardDeck: NewDiscardCardDeck(),
	}
}

type Game struct {
	share_model.Game[*Player]
	NbRow        int
	NbColumn     int
	DrawDeck     CardDeck
	DiscardDeck  CardDeck
	SelectedCard *Card
	ShouldFlip   bool
	LastTurn     bool
}

func (g *Game) CanDrawCard(player *Player) bool {
	return g.SelectedCard == nil &&
		!g.ShouldFlip &&
		!g.DrawDeck.IsEmpty() &&
		player.IsPlaying()
}

func (g *Game) CanDrawDiscardCard(player *Player) bool {
	return g.SelectedCard == nil &&
		!g.ShouldFlip &&
		!g.DiscardDeck.IsEmpty() &&
		player.IsPlaying()
}

func (g *Game) CanDiscardCard(player *Player) bool {
	return g.SelectedCard != nil &&
		!g.ShouldFlip &&
		player.IsPlaying()
}

func (g *Game) CanPutCard(player *Player, cell *PlayerCell) bool {
	return g.SelectedCard != nil &&
		!g.ShouldFlip &&
		player.IsPlaying()
}

func (g *Game) CanFlipCard(player *Player, cell *PlayerCell) bool {
	return g.SelectedCard == nil &&
		g.ShouldFlip &&
		player.IsPlaying() &&
		cell.CanFlip()
}

func (g *Game) HasSelectedCard() bool {
	return g.SelectedCard != nil
}

func (g *Game) Board(playerId share_model.PlayerId) *PlayerBoard {
	player, found := g.Player(playerId)
	if !found {
		panic(share_model.ErrPlayerNotFound)
	}
	return player.Board()
}

func (g *Game) GetBoard(playerId share_model.PlayerId) (*PlayerBoard, bool) {
	player, found := g.Player(playerId)
	if !found {
		return nil, false
	}
	return player.Board(), true
}

func (g *Game) FlipAll() {
	for _, player := range g.Players() {
		player.Board().FlipAll()
	}
}

func (g *Game) WinnerId() share_model.PlayerId {
	var winnerId share_model.PlayerId
	var bestTotal int
	for _, player := range g.Players() {
		if winnerId == "" {
			winnerId = player.Id()
			bestTotal = player.Board().Total()
		} else {
			if total := player.Board().Total(); total < bestTotal {
				winnerId = player.Id()
				bestTotal = total
			}
		}
	}
	return winnerId
}
