package model

import (
	"fmt"
	"strings"
)

type CardColor int

const (
	CardColor_Red    CardColor = 1
	CardColor_Green  CardColor = 2
	CardColor_Blue   CardColor = 3
	CardColor_Yellow CardColor = 4
)

var (
	CardColors = []CardColor{
		CardColor_Red,
		CardColor_Green,
		CardColor_Blue,
		CardColor_Yellow,
	}
)

func (c CardColor) String() string {
	switch c {
	case CardColor_Red:
		return "R"
	case CardColor_Green:
		return "G"
	case CardColor_Blue:
		return "B"
	case CardColor_Yellow:
		return "Y"
	}
	return "?"
}

func CardColorFromString(value rune) CardColor {
	switch value {
	case 'R':
		return CardColor_Red
	case 'G':
		return CardColor_Green
	case 'B':
		return CardColor_Blue
	case 'Y':
		return CardColor_Yellow
	}
	panic(fmt.Sprintf("invalid color: %v", value))
}

func (c CardColor) LabelSlice() []string {
	labels := []string{}
	switch c {
	case CardColor_Red:
		labels = append(labels, "red")
	case CardColor_Green:
		labels = append(labels, "green")
	case CardColor_Blue:
		labels = append(labels, "blue")
	case CardColor_Yellow:
		labels = append(labels, "yellow")
	}
	return labels
}

func (c CardColor) Labels() string {
	return strings.Join(c.LabelSlice(), " ")
}

type CardValue int

func CardValueFromString(value rune) CardValue {
	cardValue := CardValue(value - '0')
	if cardValue < Card_MinValue || cardValue > Card_MaxValue {
		panic(fmt.Sprintf("invalid value: %v", value))
	}
	return cardValue
}

const (
	Card_MinValue   CardValue = 1
	Card_MaxValue   CardValue = 7
	Card_NbPerValue           = 4
)

type Card int

func NewCard(value CardValue, color CardColor) Card {
	return Card(int(color)*10 + int(value))
}

func (c Card) Color() CardColor {
	return CardColor(c / 10)
}

func (c Card) IsColor(colors ...CardColor) bool {
	for _, color := range colors {
		if c.Color() == color {
			return true
		}
	}
	return false
}

func (c Card) IsNotColor(colors ...CardColor) bool {
	for _, color := range colors {
		if c.Color() == color {
			return false
		}
	}
	return true
}

func (c Card) Value() CardValue {
	return CardValue(c % 10)
}

func (c Card) IsEven() bool {
	return c.Value()%2 == 0
}

func (c Card) IsOdd() bool {
	return c.Value()%2 == 1
}

func (c Card) IsSmall() bool {
	return c.Value() < 4
}

func (c Card) IsBig() bool {
	return c.Value() > 4
}

func (c Card) Labels() string {
	labels := make([]string, 0)
	labels = append(labels, "card")
	labels = append(labels, c.Color().LabelSlice()...)
	return strings.Join(labels, " ")
}

// func (c Card) CardStr() string {
// 	return fmt.Sprintf(
// 		"<div class=\"%s\"><div class=\"symbol-top\">%d</div><div class=\"value\">%d</div><div class=\"symbol-bottom\">%d</div></div>",
// 		c.Color().Labels(),
// 		c.Value(),
// 		c.Value(),
// 		c.Value(),
// 	)
// }

// func (c Card) CardHtml() template.HTML {
// 	return template.HTML(c.CardStr())
// }

func (c Card) String() string {
	return fmt.Sprintf("%s%d", c.Color().String(), c.Value())
}

func CardFromString(value string) Card {
	return NewCard(CardValueFromString(rune(value[1])), CardColorFromString(rune(value[0])))
}

// //////////////////////////////////////////////////
// top cards

type TopCards [NbCardDeck]Card

func (t TopCards) CountColor(colors ...CardColor) int {
	count := 0
	for _, card := range t {
		for _, color := range colors {
			if card.Color() == color {
				count++
			}
		}
	}
	return count
}

func (t TopCards) Sum() int {
	sum := 0
	for _, card := range t {
		sum += int(card.Value())
	}
	return sum
}

func (t TopCards) ColorSum(colors ...CardColor) int {
	sum := 0
	for _, card := range t {
		for _, color := range colors {
			if card.Color() == color {
				sum += int(card.Value())
			}
		}
	}
	return sum
}

func (t TopCards) String() string {
	values := make([]string, 0)
	for _, card := range t {
		values = append(values, card.String())
	}
	return strings.Join(values, " ")
}

func TopCardsFromString(value string) TopCards {
	values := strings.Split(value, " ")
	if len(values) != NbCardDeck {
		panic(fmt.Sprintf("invalid top cards: %s", value))
	}
	topCards := TopCards{}
	for index, value := range values {
		topCards[index] = CardFromString(value)
	}
	return topCards
}
