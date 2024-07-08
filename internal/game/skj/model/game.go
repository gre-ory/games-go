package model

import (
	"html/template"

	"github.com/gre-ory/games-go/internal/util/loc"

	share_model "github.com/gre-ory/games-go/internal/game/share/model"
	share_websocket "github.com/gre-ory/games-go/internal/game/share/websocket"
)

func NewGame(nbRow, nbColumn int) *Game {
	return &Game{
		Game:        share_model.NewGame[*Player](),
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

const (
	MinNbPlayer = 2
	MaxNbPlayer = 4
)

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

func (g *Game) UpdateStatus() {
	if !g.CanJoin() {
		g.SetStatus(share_model.GameStatus_NotJoinableAndStartable)
		for _, player := range g.GetPlayers() {
			player.SetStatus(share_model.PlayerStatus_WaitingToStart)
		}
	} else {
		if g.CanStart() {
			g.SetStatus(share_model.GameStatus_JoinableAndStartable)
			for _, player := range g.GetPlayers() {
				player.SetStatus(share_model.PlayerStatus_WaitingToStart)
			}
		} else {
			g.SetStatus(share_model.GameStatus_JoinableNotStartable)
			for _, player := range g.GetPlayers() {
				player.SetStatus(share_model.PlayerStatus_WaitingToJoin)
			}
		}
	}
}

func (g *Game) WrapData(data share_websocket.Data, player *Player) (bool, any) {
	data = data.With("game", g)

	playerId := player.Id()
	if playerId == "" {
		return true, data
	}
	player, found := g.GetPlayer(playerId)
	if !found {
		return false, nil
	}
	return true, data.With("player", player)
}

func (g *Game) IsStopped() bool {
	return g.Status().IsStopped()
}

func (g *Game) PlayerLabels(playerId share_model.PlayerId) string {
	player, found := g.GetPlayer(playerId)
	if !found {
		return "error"
	}
	return player.Labels()
}

func (g *Game) YourPlayerMessage(localizer loc.Localizer, playerId share_model.PlayerId) template.HTML {
	player, found := g.GetPlayer(playerId)
	if !found {
		return localizer.Loc("Error", share_model.ErrPlayerNotFound.Error())
	}
	if player.HasResult() {
		result := player.Result()
		switch {
		case result.IsWin():
			return localizer.Loc("YouWin")
		case result.IsTie():
			return localizer.Loc("YouTie")
		case result.IsLoose():
			return localizer.Loc("YouLoose")
		}
	}
	return player.YourMessage(localizer)
}

func (g *Game) PlayerMessage(localizer loc.Localizer, playerId share_model.PlayerId) template.HTML {
	player, found := g.GetPlayer(playerId)
	if !found {
		return localizer.Loc("Error", share_model.ErrPlayerNotFound.Error())
	}
	if player.HasResult() {
		result := player.Result()
		switch {
		case result.IsWin():
			return localizer.Loc("PlayerWin")
		case result.IsTie():
			return localizer.Loc("PlayerTie")
		case result.IsLoose():
			return localizer.Loc("PlayerLoose")
		}
	}
	return player.Message(localizer)
}

func (g *Game) PlayerStatusIcon(playerId share_model.PlayerId) string {
	player, found := g.GetPlayer(playerId)
	if !found {
		return ""
	}
	if player.HasResult() {
		return player.Result().Icon()
	}
	return ""
}
