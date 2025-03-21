package model

import (
	"html/template"
	"strings"

	"github.com/gre-ory/games-go/internal/util/loc"
)

type PlayerStatus int

const (
	PlayerStatus_WaitingToJoin PlayerStatus = iota
	PlayerStatus_WaitingToStart
	PlayerStatus_WaitingToPlay
	PlayerStatus_Playing
	PlayerStatus_Played
)

func (s PlayerStatus) IsWaitingToJoin() bool {
	return s == PlayerStatus_WaitingToJoin
}

func (s PlayerStatus) IsWaitingToStart() bool {
	return s == PlayerStatus_WaitingToStart
}

func (s PlayerStatus) IsWaitingToPlay() bool {
	return s == PlayerStatus_WaitingToPlay
}

func (s PlayerStatus) IsPlaying() bool {
	return s == PlayerStatus_Playing
}

func (s PlayerStatus) HasPlayed() bool {
	return s == PlayerStatus_Played
}

func (s PlayerStatus) IsValid() bool {
	switch s {
	case PlayerStatus_WaitingToJoin,
		PlayerStatus_WaitingToStart,
		PlayerStatus_WaitingToPlay,
		PlayerStatus_Playing,
		PlayerStatus_Played:
		return true
	default:
		return false
	}
}

func (s PlayerStatus) String() string {
	switch s {
	case PlayerStatus_WaitingToJoin:
		return "waiting-to-join"
	case PlayerStatus_WaitingToStart:
		return "waiting-to-start"
	case PlayerStatus_WaitingToPlay:
		return "waiting-to-play"
	case PlayerStatus_Playing:
		return "playing"
	case PlayerStatus_Played:
		return "played"
	default:
		return ""
	}
}

func (s PlayerStatus) YourMessage(localizer loc.Localizer, args ...any) template.HTML {
	switch s {
	case PlayerStatus_WaitingToJoin:
		return localizer.Loc("YouWaitingToJoin", args...)
	case PlayerStatus_WaitingToStart:
		return localizer.Loc("YouWaitingToStart", args...)
	case PlayerStatus_WaitingToPlay:
		return localizer.Loc("YouWaitingToPlay", args...)
	case PlayerStatus_Playing:
		return localizer.Loc("YouPlaying", args...)
	case PlayerStatus_Played:
		return localizer.Loc("YouPlayed", args...)
	}
	return ""
}

func (s PlayerStatus) Message(localizer loc.Localizer, args ...any) template.HTML {
	switch s {
	case PlayerStatus_WaitingToJoin:
		return localizer.Loc("PlayerWaitingToJoin", args...)
	case PlayerStatus_WaitingToStart:
		return localizer.Loc("PlayerWaitingToStart", args...)
	case PlayerStatus_WaitingToPlay:
		return localizer.Loc("PlayerWaitingToPlay", args...)
	case PlayerStatus_Playing:
		return localizer.Loc("PlayerPlaying", args...)
	case PlayerStatus_Played:
		return localizer.Loc("PlayerPlayed", args...)
	}
	return ""
}

func (s PlayerStatus) LabelSlice() []string {
	var labels []string
	labels = append(labels, "player-status")
	if s.IsValid() {
		labels = append(labels, s.String())
	}
	return labels
}

func (s PlayerStatus) Labels() string {
	return strings.Join(s.LabelSlice(), " ")
}

func (s PlayerStatus) Icon() string {
	switch s {
	case PlayerStatus_WaitingToJoin,
		PlayerStatus_WaitingToStart,
		PlayerStatus_WaitingToPlay,
		PlayerStatus_Played:
		return "icon-pause"
	case PlayerStatus_Playing:
		return "icon-play"
	}
	return ""
}
