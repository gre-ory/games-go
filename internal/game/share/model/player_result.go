package model

import (
	"strings"
)

type PlayerResult int

const (
	PlayerResult_Unknown PlayerResult = iota
	PlayerResult_Win
	PlayerResult_Tie
	PlayerResult_Loose
)

func (s PlayerResult) IsUnknown() bool {
	return s == PlayerResult_Unknown
}

func (s PlayerResult) IsWin() bool {
	return s == PlayerResult_Win
}

func (s PlayerResult) IsTie() bool {
	return s == PlayerResult_Tie
}

func (s PlayerResult) IsLoose() bool {
	return s == PlayerResult_Loose
}

func (s PlayerResult) String() string {
	switch {
	case s.IsWin():
		return "win"
	case s.IsTie():
		return "tie"
	case s.IsLoose():
		return "loose"
	default:
		return ""
	}
}

func (s PlayerResult) LabelSlice() []string {
	var labels []string
	labels = append(labels, "player-result")
	labels = append(labels, s.String())
	return labels
}

func (s PlayerResult) Labels() string {
	return strings.Join(s.LabelSlice(), " ")
}

func (s PlayerResult) Icon() string {
	switch {
	case s.IsWin():
		return "icon-win"
	case s.IsTie():
		return "icon-tie"
	case s.IsLoose():
		return "icon-loose"
	}
	return ""
}
