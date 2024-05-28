package model

import (
	"math/rand"
	"strings"
)

type MissionDeck []Mission

func NewDrawMissionDeck() MissionDeck {
	deck := MissionDeck{}

	for _, color := range CardColors {
		deck.Add(NewTwoColorsNextToEachOtherMission(color))
		deck.Add(NewTwoColorsSeparatedByOneMission(color))
		deck.Add(NewTwoColorsSeparatedMission(color))
		deck.Add(NewThreeColorsMission(color))
	}

	deck.Add(NewTwoEvenSeparatedByOneMission())

	deck.Add(NewColorDoubleOfColorMission(CardColor_Green, CardColor_Blue))
	deck.Add(NewColorDoubleOfColorMission(CardColor_Blue, CardColor_Red))
	deck.Add(NewColorDoubleOfColorMission(CardColor_Red, CardColor_Yellow))
	deck.Add(NewColorDoubleOfColorMission(CardColor_Yellow, CardColor_Green))

	deck.Add(NewColorEqualColorMission(CardColor_Green, CardColor_Blue))
	deck.Add(NewColorEqualColorMission(CardColor_Blue, CardColor_Red))
	deck.Add(NewColorEqualColorMission(CardColor_Red, CardColor_Yellow))
	deck.Add(NewColorEqualColorMission(CardColor_Yellow, CardColor_Green))

	deck.Add(NewColorSumMission(2, CardColor_Yellow))
	deck.Add(NewColorSumMission(3, CardColor_Blue))
	deck.Add(NewColorSumMission(4, CardColor_Red))
	deck.Add(NewColorSumMission(6, CardColor_Green))
	deck.Add(NewColorSumMission(7, CardColor_Green))
	deck.Add(NewColorSumMission(9, CardColor_Blue))
	deck.Add(NewColorSumMission(10, CardColor_Red))
	deck.Add(NewColorSumMission(11, CardColor_Yellow))

	deck.Add(NewSumMission(10))
	deck.Add(NewSumMission(15))
	deck.Add(NewSumMission(18))
	deck.Add(NewSumMission(20))

	deck.Add(NewAllDifferentMission())
	deck.Add(NewAllDifferentColorMission())
	deck.Add(NewAllDifferentValueMission())

	deck.Add(NewAllSmallMission())
	deck.Add(NewAllBigMission())

	deck.Add(NewAllEvenMission())
	deck.Add(NewAllOddMission())

	deck.Add(NewAllTwoColorsMission(CardColor_Blue, CardColor_Yellow))
	deck.Add(NewAllTwoColorsMission(CardColor_Yellow, CardColor_Green))
	deck.Add(NewAllTwoColorsMission(CardColor_Green, CardColor_Red))
	deck.Add(NewAllTwoColorsMission(CardColor_Red, CardColor_Blue))

	deck.Add(NewFourValuesInARowMission())
	deck.Add(NewThreeOrderedValuesMission())

	deck.Shuffle()

	return deck
}

func NewDiscardMissionDeck() MissionDeck {
	return MissionDeck{}
}

func (d MissionDeck) IsEmpty() bool {
	return len(d) == 0
}

func (d *MissionDeck) Add(mission Mission) {
	(*d) = append(*d, mission)
}

func (d MissionDeck) Size() int {
	return len(d)
}

func (d MissionDeck) GetTopMission() Mission {
	if d.IsEmpty() {
		return nil
	}
	return d[len(d)-1]
}

func (d *MissionDeck) Draw() (Mission, error) {
	if d.IsEmpty() {
		return nil, ErrEmptyMissionDeck
	}
	mission := (*d)[len(*d)-1]
	(*d) = (*d)[:len(*d)-1]
	return mission, nil
}

func (d MissionDeck) Discard(mission Mission) {
	d = append(d, mission)
}

func (d *MissionDeck) Shuffle() {
	rand.Shuffle(len(*d), func(i, j int) { (*d)[i], (*d)[j] = (*d)[j], (*d)[i] })
}

func (d MissionDeck) Labels() string {
	labels := make([]string, 0)
	labels = append(labels, "mission-deck")
	if d.IsEmpty() {
		labels = append(labels, "empty")
	}
	return strings.Join(labels, " ")
}
