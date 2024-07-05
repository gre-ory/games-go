package model

import (
	"fmt"
	"strings"
)

type PlayerResult int

const (
	PlayerResult_Unknown = 0
	rank_shift           = 10
	result_unknown       = 0
	result_win           = 1
	result_tie           = 2
	result_loose         = 3
)

func NewWinResult() PlayerResult {
	return PlayerResult(result_win)
}

func NewTieResult() PlayerResult {
	return PlayerResult(result_tie)
}

func NewLooseResult() PlayerResult {
	return PlayerResult(result_loose)
}

func NewWinRankResult(rank int) PlayerResult {
	return PlayerResult(result_win + rank*rank_shift)
}

func NewTieRankResult(rank int) PlayerResult {
	return PlayerResult(result_tie + rank*rank_shift)
}

func NewLooseRankResult(rank int) PlayerResult {
	return PlayerResult(result_loose + rank*rank_shift)
}

func (s PlayerResult) result() int {
	return int(s) % rank_shift
}

func (s PlayerResult) IsUnknown() bool {
	return s.result() == result_unknown
}

func (s PlayerResult) IsWin() bool {
	return s.result() == result_win
}

func (s PlayerResult) IsTie() bool {
	return s.result() == result_tie
}

func (s PlayerResult) IsLoose() bool {
	return s.result() == result_loose
}

func (s PlayerResult) HasRank() bool {
	return s.Rank() > 0
}

func (s PlayerResult) Rank() int {
	return int(s) - s.result()/rank_shift
}

func (s PlayerResult) IsFirst() bool {
	return s.Rank() == 1
}

func (s PlayerResult) IsSecond() bool {
	return s.Rank() == 2
}

func (s PlayerResult) IsThird() bool {
	return s.Rank() == 3
}

func (s PlayerResult) HasMedal() bool {
	return s.HasRank() && s.Rank() < 4
}

func (s PlayerResult) Medal() string {
	if s.IsFirst() {
		return "gold"
	} else if s.IsSecond() {
		return "silver"
	} else if s.IsThird() {
		return "bronze"
	} else {
		return ""
	}
}

func (s PlayerResult) IsEmpty() bool {
	return s == 0
}

func (s PlayerResult) IsValid() bool {
	return s > 0
}

func (s PlayerResult) RankString() string {
	if s.HasRank() {
		return fmt.Sprintf("rank-%02d", s)
	} else {
		return ""
	}
}

func (s PlayerResult) RankLabels() []string {
	var labels []string
	if s.HasRank() {
		labels = append(labels, "player-rank")
		labels = append(labels, s.RankString())
	}
	return labels
}

func (s PlayerResult) MedalLabels() []string {
	var labels []string
	if s.HasMedal() {
		labels = append(labels, "player-medal")
		labels = append(labels, s.Medal())
	}
	return labels
}

func (s PlayerResult) String() string {
	switch {
	case s.IsUnknown():
		return "unknown"
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
	labels = append(labels, s.RankLabels()...)
	labels = append(labels, s.MedalLabels()...)
	return labels
}

func (s PlayerResult) Labels() string {
	return strings.Join(s.LabelSlice(), " ")
}
