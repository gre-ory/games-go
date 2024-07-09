package model

import (
	"fmt"
	"strings"
)

type CardColor int

const (
	CardColor_Unknown CardColor = iota
	CardColor_Blue
	CardColor_Cyan
	CardColor_Green
	CardColor_Yellow
	CardColor_Red
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

func CardColorFromValue(value int) CardColor {
	switch {
	case value < 0:
		return CardColor_Blue
	case value == 0:
		return CardColor_Cyan
	case value < 5:
		return CardColor_Green
	case value < 9:
		return CardColor_Yellow
	case value < 13:
		return CardColor_Red
	}
	return CardColor_Unknown
}

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
	Card_Unknown  = 0
	Card_Delta    = 10
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

func NewCard(value int) Card {
	color := CardColorFromValue(value)
	if value < Card_MinValue || value > Card_MaxValue {
		panic(ErrInvalidCardValue)
	}
	return Card((int(color) * 100) + value + Card_Delta)
}

func (c Card) Value() int {
	return (int(c) % 100) - Card_Delta
}

func (c Card) Color() CardColor {
	return CardColor((int(c) / 100))
}

func (c Card) Labels() string {
	labels := make([]string, 0)
	labels = append(labels, "card")
	labels = append(labels, c.Color().LabelSlice()...)
	return strings.Join(labels, " ")
}

func (c Card) String() string {
	return fmt.Sprintf("%d%s", c.Value(), c.Color().String())
}
