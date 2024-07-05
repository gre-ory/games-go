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
}

func (p *Player) CanJoin() bool {
	return p.Player.CanJoin() && p.Status().IsWaitingToJoin()
}

func (p *Player) Playing() bool {
	return p.Status().IsPlaying()
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
			return localizer.Loc("YouPlaying")
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
			return localizer.Loc("PlayerPlaying")
		}
	} else {
		return localizer.Loc("PlayerDisconnected")
	}
	return ""
}
