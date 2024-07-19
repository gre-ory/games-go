package model

import (
	share_model "github.com/gre-ory/games-go/internal/game/share/model"
)

const (
	Game_MinPlayer = 2
	Game_MaxPlayer = 4
)

func NewGame() *Game {
	game := &Game{
		Game:         share_model.NewGame[*Player](Game_MinPlayer, Game_MaxPlayer),
		DrawCardDeck: NewDrawCardDeck(),
		DiscardCardDecks: [NbCardDeck]CardDeck{
			NewDiscardCardDeck(),
			NewDiscardCardDeck(),
			NewDiscardCardDeck(),
			NewDiscardCardDeck(),
		},
		SelectedCardNumber: 0,
		DrawMissionDeck:    NewDrawMissionDeck(),
		Missions:           [NbMission]Mission{},
		DiscardMissionDeck: NewDiscardMissionDeck(),
	}
	return game
}

type GameStatus int

const (
	Joinable GameStatus = iota
	NotJoinable
	Started
	Stopped
)

const (
	NbCardDeck = 4
	NbMission  = 4
)

type TopMissions [NbMission]Mission

type Game struct {
	share_model.Game[*Player]
	DrawCardDeck          CardDeck
	DiscardCardDecks      [NbCardDeck]CardDeck
	SelectedCardNumber    int
	DrawMissionDeck       MissionDeck
	Missions              [NbMission]Mission
	DiscardMissionDeck    MissionDeck
	ValidatedMissionIndex int
	Medal                 Medal
}

func (g *Game) IsCardSelectable(playerId share_model.PlayerId, cardNumber int) bool {
	if g.IsPlayingPlayer(playerId) {
		return cardNumber != g.SelectedCardNumber
	}
	return false
}

func (g *Game) IsCardSelected(playerId share_model.PlayerId, cardNumber int) bool {
	if g.IsPlayingPlayer(playerId) {
		return cardNumber == g.SelectedCardNumber
	}
	return false
}

func (g *Game) GetTopCards() TopCards {
	topCards := TopCards{}
	for index, discardDeck := range g.DiscardCardDecks {
		topCards[index] = discardDeck.GetTopCard()
	}
	return topCards
}

func (g *Game) Play(player *Player, x, y int) error {
	// if row, ok := g.Rows[y]; ok {
	// 	return row.Play(player, x)
	// }
	return ErrOutOfRowBound
}

func (g *Game) HasValidatedMission() bool {
	return g.ValidatedMissionIndex >= 0
}

func (g *Game) GetValidatedMissionIndex() int {
	return g.ValidatedMissionIndex
}

func (g *Game) HasBronzeMedal() bool {
	return len(g.DiscardMissionDeck) > 15
}

func (g *Game) HasSilverMedal() bool {
	return len(g.DiscardMissionDeck) > 19
}

func (g *Game) HasGoldMedal() bool {
	return len(g.DiscardMissionDeck) > 23
}

func (g *Game) HasWinner() (bool, share_model.PlayerId) {
	// for x := 1; x <= 3; x++ {
	// 	same, symbol := g.HasSameSymbol(g.Rows[1].Cells[x], g.Rows[2].Cells[x], g.Rows[3].Cells[x])
	// 	if same {
	// 		return true, g.GetPlayerIdFromRune(symbol)
	// 	}
	// }
	// for y := 1; y <= 3; y++ {
	// 	same, symbol := g.HasSameSymbol(g.Rows[y].Cells[1], g.Rows[y].Cells[2], g.Rows[y].Cells[3])
	// 	if same {
	// 		return true, g.GetPlayerIdFromRune(symbol)
	// 	}
	// }
	// same, symbol := g.HasSameSymbol(g.Rows[1].Cells[1], g.Rows[2].Cells[2], g.Rows[3].Cells[3])
	// if same {
	// 	return true, g.GetPlayerIdFromRune(symbol)
	// }
	// same, symbol = g.HasSameSymbol(g.Rows[1].Cells[3], g.Rows[2].Cells[2], g.Rows[3].Cells[1])
	// if same {
	// 	return true, g.GetPlayerIdFromRune(symbol)
	// }
	return false, ""
}
