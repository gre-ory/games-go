package model

import (
	"fmt"
	"strings"
)

type CardColor int

const (
	CardColor_Blue   CardColor = 1
	CardColor_Cyan   CardColor = 2
	CardColor_Green  CardColor = 3
	CardColor_Yellow CardColor = 4
	CardColor_Red    CardColor = 5
)

var (
	CardColors = []CardColor{
		CardColor_Blue,
		CardColor_Cyan,
		CardColor_Green,
		CardColor_Yellow,
		CardColor_Red,
	}
)

func (c CardColor) String() string {
	switch c {
	case CardColor_Blue:
		return "B"
	case CardColor_Cyan:
		return "C"
	case CardColor_Green:
		return "G"
	case CardColor_Yellow:
		return "Y"
	case CardColor_Red:
		return "R"
	}
	return "?"
}

func CardColorFromString(value rune) CardColor {
	switch value {
	case 'B':
		return CardColor_Blue
	case 'C':
		return CardColor_Cyan
	case 'G':
		return CardColor_Green
	case 'Y':
		return CardColor_Yellow
	case 'R':
		return CardColor_Red
	}
	panic(ErrInvalidCardColor)
}

func (c CardColor) LabelColor() string {
	switch c {
	case CardColor_Blue:
		return "blue"
	case CardColor_Cyan:
		return "cyan"
	case CardColor_Green:
		return "green"
	case CardColor_Yellow:
		return "yellow"
	case CardColor_Red:
		return "red"
	}
	return ""
}

func (c CardColor) LabelSlice() []string {
	labels := []string{}
	labels = append(labels, c.LabelColor())
	return labels
}

func (c CardColor) Labels() string {
	return strings.Join(c.LabelSlice(), " ")
}

const (
	Card_MinValue = -2
	Card_MaxValue = 12
)

var (
	Card_NbPerValue = map[Card]int{
		-2: 5,
		-1: 10,
		0:  15,
		1:  10,
		2:  10,
		3:  10,
		4:  10,
		5:  10,
		6:  10,
		7:  10,
		8:  10,
		9:  10,
		10: 10,
		11: 10,
		12: 10,
	}
)

type Card int

func CardFromString(value rune) Card {
	card := Card(value - '0')
	if card < Card_MinValue || card > Card_MaxValue {
		panic(ErrInvalidCard)
	}
	return card
}

func (c Card) ToRune() rune {
	// note:
	// -2 = .
	// -1 = /
	// 10 = :
	// 11 = ;
	// 12 = <
	return '0' + rune(c)
}

func (c Card) Color() CardColor {
	switch {
	case c < 0:
		return CardColor_Blue
	case c == 0:
		return CardColor_Cyan
	case c < 5:
		return CardColor_Green
	case c < 9:
		return CardColor_Yellow
	default:
		return CardColor_Red
	}
}

func (c Card) Labels() string {
	labels := make([]string, 0)
	labels = append(labels, "card")
	labels = append(labels, c.Color().LabelSlice()...)
	return strings.Join(labels, " ")
}

func (c Card) String() string {
	return fmt.Sprintf("%d", c, c.Color().String())
}
