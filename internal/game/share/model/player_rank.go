package model

import (
	"fmt"
	"strings"
)

type PlayerRank int

func (s PlayerRank) IsUnknown() bool {
	return s == 0
}

func (s PlayerRank) IsFirst() bool {
	return s == 1
}

func (s PlayerRank) IsSecond() bool {
	return s == 2
}

func (s PlayerRank) IsThird() bool {
	return s == 3
}

func (s PlayerRank) IsValid() bool {
	return s > 0
}

func (s PlayerRank) String() string {
	if s.IsValid() {
		return fmt.Sprintf("rank-%02d", s)
	} else {
		return ""
	}
}

func (s PlayerRank) LabelSlice() []string {
	var labels []string
	labels = append(labels, "player-rank")
	labels = append(labels, s.String())
	return labels
}

func (s PlayerRank) Labels() string {
	return strings.Join(s.LabelSlice(), " ")
}

func (s PlayerRank) HasMedal() bool {
	return 0 < s && s < 4
}

func (s PlayerRank) MedalString() string {
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

func (s PlayerRank) MedalLabelSlice() []string {
	var labels []string
	if s.HasMedal() {
		labels = append(labels, "player-medal")
		labels = append(labels, s.MedalString())
	}
	return labels
}
