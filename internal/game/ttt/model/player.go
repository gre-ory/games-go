package model

import (
	"html/template"

	share_model "github.com/gre-ory/games-go/internal/game/share/model"
)

func NewPlayerFromUser(gameId share_model.GameId, user share_model.User) *Player {
	return &Player{
		Player: share_model.NewPlayerFromUser(gameId, user),
	}
}

type Player struct {
	share_model.Player
	Symbol rune
}

func (p *Player) SetSymbol(symbol rune) {
	p.Symbol = symbol
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

func (p *Player) LabelSlice() []string {
	labels := make([]string, 0)
	labels = append(labels, p.Player.LabelSlice()...)
	switch p.Symbol {
	case PLAYER_ONE_SYMBOL:
		labels = append(labels, "symbol-1")
	case PLAYER_TWO_SYMBOL:
		labels = append(labels, "symbol-2")
	}
	return labels
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
