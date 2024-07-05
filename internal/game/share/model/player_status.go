package model

import "strings"

type PlayerStatus int

const (
	PlayerStatus_WaitingToJoin PlayerStatus = iota
	PlayerStatus_WaitingToStart
	PlayerStatus_WaitingToPlay
	PlayerStatus_Playing
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

func (s PlayerStatus) IsValid() bool {
	switch s {
	case PlayerStatus_WaitingToJoin,
		PlayerStatus_WaitingToStart,
		PlayerStatus_WaitingToPlay,
		PlayerStatus_Playing:
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
	default:
		return ""
	}
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
