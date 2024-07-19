package model

import (
	share_model "github.com/gre-ory/games-go/internal/game/share/model"
)

const (
	MinNbPlayer = 2
	MaxNbPlayer = 4
)

func NewGame(nbRow, nbColumn int) *Game {
	return &Game{
		Game:        share_model.NewGame[*Player](MinNbPlayer, MaxNbPlayer),
		NbRow:       nbRow,
		NbColumn:    nbColumn,
		DrawDeck:    NewDrawCardDeck(),
		DiscardDeck: NewDiscardCardDeck(),
		boards:      make(map[share_model.PlayerId]*PlayerBoard),
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
	boards       map[share_model.PlayerId]*PlayerBoard
}

func (g *Game) CanJoin() bool {
	return g.NbPlayer() < MaxNbPlayer
}

func (g *Game) CanStart() bool {
	return g.NbPlayer() >= MinNbPlayer
}

func (g *Game) AddBoard(playerId share_model.PlayerId, board *PlayerBoard) {
	g.boards[playerId] = board
}

func (g *Game) GetBoard(playerId share_model.PlayerId) (*PlayerBoard, bool) {
	board, found := g.boards[playerId]
	return board, found
}
