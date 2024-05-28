package model

import (
	"math/rand"
	"strings"
)

type CardDeck []Card

func NewDrawCardDeck() CardDeck {
	deck := CardDeck{}
	for i := 1; i <= Card_NbPerValue; i++ {
		for _, color := range CardColors {
			for value := Card_MinValue; value <= Card_MaxValue; value++ {
				deck.Add(NewCard(value, color))
			}
		}
	}
	deck.Shuffle()
	return deck
}

func NewDiscardCardDeck() CardDeck {
	return CardDeck{}
}

func (d CardDeck) IsEmpty() bool {
	return len(d) == 0
}

func (d CardDeck) Size() int {
	return len(d)
}

func (d CardDeck) GetTopCard() Card {
	if d.IsEmpty() {
		return 0
	}
	return d[len(d)-1]
}

func (d *CardDeck) Draw() (Card, error) {
	if d.IsEmpty() {
		return 0, ErrEmptyCardDeck
	}
	card := (*d)[len(*d)-1]
	(*d) = (*d)[:len(*d)-1]
	return card, nil
}

func (d *CardDeck) Add(card Card) {
	*d = append(*d, card)
}

func (d *CardDeck) Shuffle() {
	rand.Shuffle(len(*d), func(i, j int) { (*d)[i], (*d)[j] = (*d)[j], (*d)[i] })
}

func (d CardDeck) Labels() string {
	labels := make([]string, 0)
	labels = append(labels, "card-deck")
	if d.IsEmpty() {
		labels = append(labels, "empty")
	}
	return strings.Join(labels, " ")
}
