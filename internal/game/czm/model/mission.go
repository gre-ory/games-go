package model

import (
	"sort"
)

type Mission interface {
	IsCompleted(cards TopCards) bool
}

// //////////////////////////////////////////////////
// two colors next to each other

func NewTwoColorsNextToEachOtherMission(color CardColor) Mission {
	return &twoColorsNextToEachOtherMission{
		color: color,
	}
}

type twoColorsNextToEachOtherMission struct {
	color CardColor
}

func (m *twoColorsNextToEachOtherMission) IsCompleted(cards TopCards) bool {
	if cards.CountColor(m.color) != 2 {
		return false
	}
	if cards[0].IsColor(m.color) {
		return cards[1].IsColor(m.color)
	} else if cards[1].IsColor(m.color) {
		return cards[2].IsColor(m.color)
	} else if cards[2].IsColor(m.color) {
		return cards[3].IsColor(m.color)
	}
	return false
}

// //////////////////////////////////////////////////
// two colors separated by one

func NewTwoColorsSeparatedByOneMission(color CardColor) Mission {
	return &twoColorsSeparatedByOneMission{
		color: color,
	}
}

type twoColorsSeparatedByOneMission struct {
	color CardColor
}

func (m *twoColorsSeparatedByOneMission) IsCompleted(cards TopCards) bool {
	if cards.CountColor(m.color) != 2 {
		return false
	}
	if cards[0].IsColor(m.color) {
		return cards[2].IsColor(m.color)
	} else if cards[1].IsColor(m.color) {
		return cards[3].IsColor(m.color)
	}
	return false
}

// //////////////////////////////////////////////////
// two colors separated

func NewTwoColorsSeparatedMission(color CardColor) Mission {
	return &twoColorsSeparatedMission{
		color: color,
	}
}

type twoColorsSeparatedMission struct {
	color CardColor
}

func (m *twoColorsSeparatedMission) IsCompleted(cards TopCards) bool {
	if cards.CountColor(m.color) != 2 {
		return false
	}
	if cards[0].IsColor(m.color) {
		return cards[2].IsColor(m.color) || cards[3].IsColor(m.color)
	} else if cards[1].IsColor(m.color) {
		return cards[3].IsColor(m.color)
	}
	return false
}

// //////////////////////////////////////////////////
// three colors

func NewThreeColorsMission(color CardColor) Mission {
	return &threeColorsMission{
		color: color,
	}
}

type threeColorsMission struct {
	color CardColor
}

func (m *threeColorsMission) IsCompleted(cards TopCards) bool {
	return cards.CountColor(m.color) == 3
}

// //////////////////////////////////////////////////
// color double of color

func NewColorDoubleOfColorMission(color1 CardColor, color2 CardColor) Mission {
	return &colorDoubleOfColorMission{
		color1: color1,
		color2: color2,
	}
}

type colorDoubleOfColorMission struct {
	color1, color2 CardColor
}

func (m *colorDoubleOfColorMission) IsCompleted(cards TopCards) bool {
	sumColor1 := cards.ColorSum(m.color1)
	sumColor2 := cards.ColorSum(m.color2)
	return sumColor1 != 0 && sumColor2 != 0 && sumColor1 == 2*sumColor2
}

// //////////////////////////////////////////////////
// color equal color

func NewColorEqualColorMission(color1 CardColor, color2 CardColor) Mission {
	return &colorEqualColorMission{
		color1: color1,
		color2: color2,
	}
}

type colorEqualColorMission struct {
	color1, color2 CardColor
}

func (m *colorEqualColorMission) IsCompleted(cards TopCards) bool {
	sumColor1 := cards.ColorSum(m.color1)
	sumColor2 := cards.ColorSum(m.color2)
	return sumColor1 != 0 && sumColor2 != 0 && sumColor1 == sumColor2
}

// //////////////////////////////////////////////////
// color sum

func NewColorSumMission(sum int, color CardColor) Mission {
	return &colorSumMission{
		sum:   sum,
		color: color,
	}
}

type colorSumMission struct {
	sum   int
	color CardColor
}

func (m *colorSumMission) IsCompleted(cards TopCards) bool {
	sumColor := cards.ColorSum(m.color)
	return sumColor != 0 && sumColor == m.sum
}

// //////////////////////////////////////////////////
// sum

func NewSumMission(sum int) Mission {
	return &sumMission{
		sum: sum,
	}
}

type sumMission struct {
	sum int
}

func (m *sumMission) IsCompleted(cards TopCards) bool {
	sum := cards.Sum()
	return sum != 0 && sum == m.sum
}

// //////////////////////////////////////////////////
// all different

func NewAllDifferentMission() Mission {
	return &allDifferentMission{}
}

type allDifferentMission struct {
	allDifferentColorMission
	allDifferentValueMission
}

func (m *allDifferentMission) IsCompleted(cards TopCards) bool {
	return m.allDifferentColorMission.IsCompleted(cards) &&
		m.allDifferentValueMission.IsCompleted(cards)
}

// //////////////////////////////////////////////////
// all different color

func NewAllDifferentColorMission() Mission {
	return &allDifferentColorMission{}
}

type allDifferentColorMission struct {
}

func (m *allDifferentColorMission) IsCompleted(cards TopCards) bool {
	colors := make(map[CardColor]struct{})
	for _, card := range cards {
		colors[card.Color()] = struct{}{}
	}
	return len(colors) == 4
}

// //////////////////////////////////////////////////
// all different value

func NewAllDifferentValueMission() Mission {
	return &allDifferentValueMission{}
}

type allDifferentValueMission struct {
}

func (m *allDifferentValueMission) IsCompleted(cards TopCards) bool {
	values := make(map[CardValue]struct{})
	for _, card := range cards {
		values[card.Value()] = struct{}{}
	}
	return len(values) == 4
}

// //////////////////////////////////////////////////
// two even separated by one

func NewTwoEvenSeparatedByOneMission() Mission {
	return &twoEvenSeparatedByOneMission{}
}

type twoEvenSeparatedByOneMission struct {
}

func (m *twoEvenSeparatedByOneMission) IsCompleted(cards TopCards) bool {
	if cards[0].IsEven() && cards[2].IsEven() {
		return cards[1].IsOdd() && cards[3].IsOdd()
	} else if cards[1].IsEven() && cards[3].IsEven() {
		return cards[0].IsOdd() && cards[2].IsOdd()
	}
	return false
}

// //////////////////////////////////////////////////
// 4 values in a row

func NewFourValuesInARowMission() Mission {
	return &fourValuesInARowMission{}
}

type fourValuesInARowMission struct {
}

func (m *fourValuesInARowMission) IsCompleted(cards TopCards) bool {
	values := make([]CardValue, 4)
	for i, card := range cards {
		values[i] = card.Value()
	}
	sort.Slice(values, func(i, j int) bool {
		return values[i] > values[j]
	})
	if values[1] == values[0]+1 {
		return values[2] == values[0]+2 && values[3] == values[0]+3
	} else if values[1] == values[0]-1 {
		return values[2] == values[0]-2 && values[3] == values[0]-3
	}
	return false
}

// //////////////////////////////////////////////////
// three ordered values

func NewThreeOrderedValuesMission() Mission {
	return &threeOrderedValuesMission{}
}

type threeOrderedValuesMission struct {
}

func (m *threeOrderedValuesMission) IsCompleted(cards TopCards) bool {
	if cards[1].Value() == cards[0].Value()+1 {
		return cards[2].Value() == cards[0].Value()+2
	} else if cards[1].Value() == cards[0].Value()-1 {
		return cards[2].Value() == cards[0].Value()-2
	} else if cards[2].Value() == cards[1].Value()+1 {
		return cards[3].Value() == cards[1].Value()+2
	} else if cards[2].Value() == cards[1].Value()-1 {
		return cards[3].Value() == cards[1].Value()-2
	}
	return false
}

// //////////////////////////////////////////////////
// all of

func NewAllSmallMission() Mission {
	return newAllOfMission(IsSmall)
}

func NewAllBigMission() Mission {
	return newAllOfMission(IsBig)
}

func NewAllEvenMission() Mission {
	return newAllOfMission(IsEven)
}

func NewAllOddMission() Mission {
	return newAllOfMission(IsOdd)
}

func NewAllTwoColorsMission(color1, color2 CardColor) Mission {
	return newAllOfMission(IsColor(color1, color2))
}

func newAllOfMission(checkFn func(card Card) bool) Mission {
	return &cardMission{
		checkFn: checkFn,
	}
}

type cardMission struct {
	checkFn func(card Card) bool
}

func (m *cardMission) IsCompleted(cards TopCards) bool {
	for _, card := range cards {
		if !m.checkFn(card) {
			return false
		}
	}
	return true
}

func IsEven(card Card) bool {
	return card.IsEven()
}

func IsOdd(card Card) bool {
	return card.IsOdd()
}

func IsSmall(card Card) bool {
	return card.IsSmall()
}

func IsBig(card Card) bool {
	return card.IsBig()
}

func IsColor(colors ...CardColor) func(card Card) bool {
	return func(card Card) bool {
		return card.IsColor(colors...)
	}
}
