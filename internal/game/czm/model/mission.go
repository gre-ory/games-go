package model

import (
	"sort"
)

type Mission interface {
	IsCompleted(cards TopCards) bool
	GetTpl() (string, map[string]any)
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

func (m *twoColorsNextToEachOtherMission) GetTpl() (string, map[string]any) {
	return "mission-two-colors-next-to-each-other", map[string]any{
		"color": m.color.LabelColor(),
	}
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

func (m *twoColorsSeparatedByOneMission) GetTpl() (string, map[string]any) {
	return "mission-two-colors-separated-by-one", map[string]any{
		"color": m.color.LabelColor(),
	}
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

func (m *twoColorsSeparatedMission) GetTpl() (string, map[string]any) {
	return "mission-two-colors-separated", map[string]any{
		"color": m.color.LabelColor(),
	}
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

func (m *threeColorsMission) GetTpl() (string, map[string]any) {
	return "mission-three-colors", map[string]any{
		"color": m.color.LabelColor(),
	}
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

func (m *colorDoubleOfColorMission) GetTpl() (string, map[string]any) {
	return "mission-color-double-of-color", map[string]any{
		"color1": m.color1.LabelColor(),
		"color2": m.color2.LabelColor(),
	}
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

func (m *colorEqualColorMission) GetTpl() (string, map[string]any) {
	return "mission-color-equal-color", map[string]any{
		"color1": m.color1.LabelColor(),
		"color2": m.color2.LabelColor(),
	}
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

func (m *colorSumMission) GetTpl() (string, map[string]any) {
	return "mission-color-sum", map[string]any{
		"sum":   m.sum,
		"color": m.color.LabelColor(),
	}
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

func (m *sumMission) GetTpl() (string, map[string]any) {
	return "mission-sum", map[string]any{
		"sum": m.sum,
	}
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

func (m *allDifferentMission) GetTpl() (string, map[string]any) {
	return "mission-all-different", nil
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

func (m *allDifferentColorMission) GetTpl() (string, map[string]any) {
	return "mission-all-different-color", nil
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

func (m *allDifferentValueMission) GetTpl() (string, map[string]any) {
	return "mission-all-different-value", nil
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

func (m *twoEvenSeparatedByOneMission) GetTpl() (string, map[string]any) {
	return "mission-two-even-separated-by-one", nil
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

func (m *fourValuesInARowMission) GetTpl() (string, map[string]any) {
	return "mission-four-values-in-a-row", nil
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

func (m *threeOrderedValuesMission) GetTpl() (string, map[string]any) {
	return "mission-three-ordered-values", nil
}

// //////////////////////////////////////////////////
// all small

func NewAllSmallMission() Mission {
	return &allSmallMission{}
}

type allSmallMission struct {
}

func (m *allSmallMission) IsCompleted(cards TopCards) bool {
	return allOf(cards, func(card Card) bool {
		return card.IsSmall()
	})
}

func (m *allSmallMission) GetTpl() (string, map[string]any) {
	return "mission-all-small", nil
}

// //////////////////////////////////////////////////
// all big

func NewAllBigMission() Mission {
	return &allBigMission{}
}

type allBigMission struct {
}

func (m *allBigMission) IsCompleted(cards TopCards) bool {
	return allOf(cards, func(card Card) bool {
		return card.IsBig()
	})
}

func (m *allBigMission) GetTpl() (string, map[string]any) {
	return "mission-all-big", nil
}

// //////////////////////////////////////////////////
// all even

func NewAllEvenMission() Mission {
	return &allEvenMission{}
}

type allEvenMission struct {
}

func (m *allEvenMission) IsCompleted(cards TopCards) bool {
	return allOf(cards, func(card Card) bool {
		return card.IsEven()
	})
}

func (m *allEvenMission) GetTpl() (string, map[string]any) {
	return "mission-all-even", nil
}

// //////////////////////////////////////////////////
// all odd

func NewAllOddMission() Mission {
	return &allOddMission{}
}

type allOddMission struct {
}

func (m *allOddMission) IsCompleted(cards TopCards) bool {
	return allOf(cards, func(card Card) bool {
		return card.IsOdd()
	})
}

func (m *allOddMission) GetTpl() (string, map[string]any) {
	return "mission-all-odd", nil
}

// //////////////////////////////////////////////////
// all two colors

func NewAllTwoColorsMission(color1, color2 CardColor) Mission {
	return &allTwoColorsMission{
		color1: color1,
		color2: color2,
	}
}

type allTwoColorsMission struct {
	color1 CardColor
	color2 CardColor
}

func (m *allTwoColorsMission) IsCompleted(cards TopCards) bool {
	return allOf(cards, func(card Card) bool {
		return card.IsColor(m.color1, m.color2)
	})
}

func (m *allTwoColorsMission) GetTpl() (string, map[string]any) {
	return "mission-all-two-colors", map[string]any{
		"color1": m.color1.LabelColor(),
		"color2": m.color2.LabelColor(),
	}
}

// //////////////////////////////////////////////////
// all of

func allOf(cards TopCards, checkFn func(card Card) bool) bool {
	for _, card := range cards {
		if !checkFn(card) {
			return false
		}
	}
	return true
}
