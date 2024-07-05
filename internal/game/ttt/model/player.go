package model

import (
	"html/template"

	"github.com/gre-ory/games-go/internal/util/loc"

	share_model "github.com/gre-ory/games-go/internal/game/share/model"
	share_websocket "github.com/gre-ory/games-go/internal/game/share/websocket"
)

func NewPlayer(player share_websocket.Player) *Player {
	return &Player{
		Player: player,
	}
}

type Player struct {
	share_websocket.Player
	Symbol rune
}

func (p *Player) SetSymbol(symbol rune) {
	p.Symbol = symbol
}

func (p *Player) CanJoin() bool {
	return p.Player.CanJoin() && p.Status().IsWaitingToJoin()
}

func (p *Player) Playing() bool {
	return p.Status().IsPlaying()
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

func (p *Player) YourMessage(localizer loc.Localizer) template.HTML {
	if p.IsActive() {
		switch p.Status() {
		case share_model.PlayerStatus_WaitingToJoin:
			return localizer.Loc("YouWaitingToJoin")
		case share_model.PlayerStatus_WaitingToStart:
			return localizer.Loc("YouWaitingToStart")
		case share_model.PlayerStatus_WaitingToPlay:
			return localizer.Loc("YouWaitingToPlay")
		case share_model.PlayerStatus_Playing:
			return localizer.Loc("YouPlaying", p.IconHtml())
		}
	} else {
		return localizer.Loc("YouDisconnected")
	}
	return ""
}

func (p *Player) Message(localizer loc.Localizer) template.HTML {
	if p.IsActive() {
		switch p.Status() {
		case share_model.PlayerStatus_WaitingToJoin:
			return localizer.Loc("PlayerWaitingToJoin")
		case share_model.PlayerStatus_WaitingToStart:
			return localizer.Loc("PlayerWaitingToStart")
		case share_model.PlayerStatus_WaitingToPlay:
			return localizer.Loc("PlayerWaitingToPlay")
		case share_model.PlayerStatus_Playing:
			return localizer.Loc("PlayerPlaying", p.IconHtml())
		}
	} else {
		return localizer.Loc("PlayerDisconnected")
	}
	return ""
}

func (p *Player) StatusIcon() string {
	switch p.Status() {
	case share_model.PlayerStatus_WaitingToJoin,
		share_model.PlayerStatus_WaitingToStart,
		share_model.PlayerStatus_WaitingToPlay:
		return "icon-pause"
	case share_model.PlayerStatus_Playing:
		return "icon-play"
	}
	return ""
}
