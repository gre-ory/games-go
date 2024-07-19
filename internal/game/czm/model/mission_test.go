package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMissions(t *testing.T) {

	type TestSubCase struct {
		cards         string
		wantCompleted bool
	}

	type TestCase struct {
		mission  Mission
		subCases map[string]TestSubCase
	}

	testCases := map[string]TestCase{
		"2-reds-next-to-each-other": {
			mission: NewTwoColorsNextToEachOtherMission(CardColor_Red),
			subCases: map[string]TestSubCase{
				"none":    {cards: "B1 B2 B3 B4", wantCompleted: false},
				"one":     {cards: "B1 R2 B3 B4", wantCompleted: false},
				"two-1-2": {cards: "R1 R2 B3 B4", wantCompleted: true},
				"two-2-3": {cards: "B1 R2 R3 B4", wantCompleted: true},
				"two-3-4": {cards: "B1 B2 R3 R4", wantCompleted: true},
				"two-1-4": {cards: "R1 B2 B3 R4", wantCompleted: false},
				"two-2-4": {cards: "B1 R2 B3 R4", wantCompleted: false},
				"three":   {cards: "R1 R2 B3 R4", wantCompleted: false},
				"four":    {cards: "R1 R2 R3 R4", wantCompleted: false},
			},
		},
		"2-reds-separated-by-one": {
			mission: NewTwoColorsSeparatedByOneMission(CardColor_Red),
			subCases: map[string]TestSubCase{
				"none":    {cards: "B1 B2 B3 B4", wantCompleted: false},
				"one":     {cards: "B1 R2 B3 B4", wantCompleted: false},
				"two-1-2": {cards: "R1 R2 B3 B4", wantCompleted: false},
				"two-2-3": {cards: "B1 R2 R3 B4", wantCompleted: false},
				"two-3-4": {cards: "B1 B2 R3 R4", wantCompleted: false},
				"two-1-3": {cards: "R1 B2 R3 B4", wantCompleted: true},
				"two-2-4": {cards: "B1 R2 B3 R4", wantCompleted: true},
				"two-1-4": {cards: "R1 B2 B3 R4", wantCompleted: false},
				"three":   {cards: "R1 R2 B3 R4", wantCompleted: false},
				"four":    {cards: "R1 R2 R3 R4", wantCompleted: false},
			},
		},
		"2-reds-separated": {
			mission: NewTwoColorsSeparatedMission(CardColor_Red),
			subCases: map[string]TestSubCase{
				"none":    {cards: "B1 B2 B3 B4", wantCompleted: false},
				"one":     {cards: "B1 R2 B3 B4", wantCompleted: false},
				"two-1-2": {cards: "R1 R2 B3 B4", wantCompleted: false},
				"two-2-3": {cards: "B1 R2 R3 B4", wantCompleted: false},
				"two-3-4": {cards: "B1 B2 R3 R4", wantCompleted: false},
				"two-1-3": {cards: "R1 B2 R3 B4", wantCompleted: true},
				"two-2-4": {cards: "B1 R2 B3 R4", wantCompleted: true},
				"two-1-4": {cards: "R1 B2 B3 R4", wantCompleted: true},
				"three":   {cards: "R1 R2 B3 R4", wantCompleted: false},
				"four":    {cards: "R1 R2 R3 R4", wantCompleted: false},
			},
		},
		"3-reds": {
			mission: NewThreeColorsMission(CardColor_Red),
			subCases: map[string]TestSubCase{
				"none":        {cards: "B1 B2 B3 B4", wantCompleted: false},
				"one":         {cards: "B1 R2 B3 B4", wantCompleted: false},
				"two-1-2":     {cards: "R1 R2 B3 B4", wantCompleted: false},
				"two-2-3":     {cards: "B1 R2 R3 B4", wantCompleted: false},
				"two-3-4":     {cards: "B1 B2 R3 R4", wantCompleted: false},
				"two-1-3":     {cards: "R1 B2 R3 B4", wantCompleted: false},
				"two-2-4":     {cards: "B1 R2 B3 R4", wantCompleted: false},
				"two-1-4":     {cards: "R1 B2 B3 R4", wantCompleted: false},
				"three-1-2-3": {cards: "R1 R2 R3 B4", wantCompleted: true},
				"three-1-2-4": {cards: "R1 R2 B3 R4", wantCompleted: true},
				"three-1-3-4": {cards: "R1 B2 R3 R4", wantCompleted: true},
				"three-2-3-4": {cards: "B1 R2 R3 R4", wantCompleted: true},
				"four":        {cards: "R1 R2 R3 R4", wantCompleted: false},
			},
		},
		"2-even-separated-by-one": {
			mission: NewTwoEvenSeparatedByOneMission(),
			subCases: map[string]TestSubCase{
				"none":    {cards: "Y3 G7 R7 R1", wantCompleted: false},
				"one":     {cards: "Y3 G7 R7 R2", wantCompleted: false},
				"two-1-2": {cards: "Y2 Y4 R7 G3", wantCompleted: false},
				"two-2-3": {cards: "Y5 Y2 R6 G3", wantCompleted: false},
				"two-3-4": {cards: "Y5 Y1 R6 G6", wantCompleted: false},
				"two-1-3": {cards: "Y2 Y1 R6 G3", wantCompleted: true},
				"two-2-4": {cards: "Y5 Y2 R3 G4", wantCompleted: true},
				"two-1-4": {cards: "Y2 Y5 R3 G4", wantCompleted: false},
				"three":   {cards: "Y4 Y2 R3 G4", wantCompleted: false},
				"four":    {cards: "Y4 Y2 R6 G4", wantCompleted: false},
			},
		},
		"blue-is-double-of-yellow": {
			mission: NewColorDoubleOfColorMission(CardColor_Blue, CardColor_Yellow),
			subCases: map[string]TestSubCase{
				"none":       {cards: "G3 G7 R7 R1", wantCompleted: false},
				"equal":      {cards: "G3 B2 R7 Y2", wantCompleted: false},
				"two":        {cards: "G1 B6 R7 Y3", wantCompleted: true},
				"three":      {cards: "Y1 B6 R7 Y2", wantCompleted: true},
				"four":       {cards: "Y1 B3 B3 Y2", wantCompleted: true},
				"not-enough": {cards: "Y1 B5 R7 Y2", wantCompleted: false},
				"too-much":   {cards: "Y1 B7 R7 Y2", wantCompleted: false},
				"no-yellow":  {cards: "G3 B4 R7 R2", wantCompleted: false},
				"no-blue":    {cards: "G3 R4 R7 Y2", wantCompleted: false},
			},
		},
		"blue-equal-yellow": {
			mission: NewColorEqualColorMission(CardColor_Blue, CardColor_Yellow),
			subCases: map[string]TestSubCase{
				"none":       {cards: "G3 G7 R7 R1", wantCompleted: false},
				"two":        {cards: "G1 B4 R7 Y4", wantCompleted: true},
				"three":      {cards: "Y4 B6 R7 Y2", wantCompleted: true},
				"four":       {cards: "Y1 B3 B3 Y5", wantCompleted: true},
				"not-enough": {cards: "Y1 B4 R7 Y4", wantCompleted: false},
				"too-much":   {cards: "Y1 B1 B1 Y2", wantCompleted: false},
				"no-yellow":  {cards: "G3 B4 R7 R2", wantCompleted: false},
				"no-blue":    {cards: "G3 R4 R7 Y2", wantCompleted: false},
			},
		},
		"blue-equal-5": {
			mission: NewColorSumMission(5, CardColor_Blue),
			subCases: map[string]TestSubCase{
				"none":       {cards: "G3 G7 R7 R1", wantCompleted: false},
				"one":        {cards: "G1 B5 R7 Y4", wantCompleted: true},
				"two":        {cards: "G1 B4 R7 B1", wantCompleted: true},
				"three":      {cards: "B1 B3 B1 Y2", wantCompleted: true},
				"four":       {cards: "B1 B2 B1 B1", wantCompleted: true},
				"not-enough": {cards: "B1 Y2 B3 G4", wantCompleted: false},
				"too-much":   {cards: "Y1 B1 R5 B5", wantCompleted: false},
				"no-blue":    {cards: "G3 Y1 R7 R2", wantCompleted: false},
			},
		},
		"sum-is-9": {
			mission: NewSumMission(9),
			subCases: map[string]TestSubCase{
				"not-enough": {cards: "B1 Y2 B3 G2", wantCompleted: false},
				"ok":         {cards: "G1 G2 R5 R1", wantCompleted: true},
				"too-much":   {cards: "Y1 B1 R5 B3", wantCompleted: false},
			},
		},
		"all-different": {
			mission: NewAllDifferentMission(),
			subCases: map[string]TestSubCase{
				"same-card":            {cards: "B7 Y2 R1 B7", wantCompleted: false},
				"same-color-and-value": {cards: "B1 Y2 R2 B7", wantCompleted: false},
				"same-value":           {cards: "G1 Y2 R2 B7", wantCompleted: false},
				"same-color":           {cards: "B1 Y2 R5 B7", wantCompleted: false},
				"ok":                   {cards: "G1 Y2 R5 B7", wantCompleted: true},
			},
		},
		"all-different-color": {
			mission: NewAllDifferentColorMission(),
			subCases: map[string]TestSubCase{
				"same-card":            {cards: "B7 Y2 R1 B7", wantCompleted: false},
				"same-color-and-value": {cards: "B1 Y2 R2 B7", wantCompleted: false},
				"same-value":           {cards: "G1 Y2 R2 B7", wantCompleted: true},
				"same-color":           {cards: "B1 Y2 R5 B7", wantCompleted: false},
				"ok":                   {cards: "G1 Y2 R5 B7", wantCompleted: true},
			},
		},
		"all-different-value": {
			mission: NewAllDifferentValueMission(),
			subCases: map[string]TestSubCase{
				"same-card":            {cards: "B7 Y2 R1 B7", wantCompleted: false},
				"same-color-and-value": {cards: "B1 Y2 R2 B7", wantCompleted: false},
				"same-value":           {cards: "G1 Y2 R2 B7", wantCompleted: false},
				"same-color":           {cards: "B1 Y2 R5 B7", wantCompleted: true},
				"ok":                   {cards: "G1 Y2 R5 B7", wantCompleted: true},
			},
		},
		"all-small": {
			mission: NewAllSmallMission(),
			subCases: map[string]TestSubCase{
				"none":  {cards: "G4 Y4 Y4 B4", wantCompleted: false},
				"one":   {cards: "G1 Y7 Y4 B6", wantCompleted: false},
				"two":   {cards: "G5 Y2 Y1 B4", wantCompleted: false},
				"three": {cards: "G1 Y4 Y1 B3", wantCompleted: false},
				"four":  {cards: "G1 Y2 Y1 B3", wantCompleted: true},
			},
		},
		"all-big": {
			mission: NewAllBigMission(),
			subCases: map[string]TestSubCase{
				"none":  {cards: "G4 Y4 Y4 B4", wantCompleted: false},
				"one":   {cards: "G1 Y7 Y4 B3", wantCompleted: false},
				"two":   {cards: "G5 Y2 Y6 B4", wantCompleted: false},
				"three": {cards: "G7 Y5 Y1 B5", wantCompleted: false},
				"four":  {cards: "G7 Y6 Y7 B5", wantCompleted: true},
			},
		},
		"all-even": {
			mission: NewAllEvenMission(),
			subCases: map[string]TestSubCase{
				"none":  {cards: "G1 Y3 Y1 B7", wantCompleted: false},
				"one":   {cards: "G1 Y7 Y4 B3", wantCompleted: false},
				"two":   {cards: "G5 Y2 Y1 B4", wantCompleted: false},
				"three": {cards: "G6 Y4 Y1 B6", wantCompleted: false},
				"four":  {cards: "G4 Y2 Y2 B6", wantCompleted: true},
			},
		},
		"all-odd": {
			mission: NewAllOddMission(),
			subCases: map[string]TestSubCase{
				"none":  {cards: "G4 Y4 Y4 B4", wantCompleted: false},
				"one":   {cards: "G2 Y7 Y4 B4", wantCompleted: false},
				"two":   {cards: "G5 Y2 Y3 B4", wantCompleted: false},
				"three": {cards: "G7 Y5 Y2 B5", wantCompleted: false},
				"four":  {cards: "G7 Y7 Y7 B5", wantCompleted: true},
			},
		},
		"all-red-or-blue": {
			mission: NewAllTwoColorsMission(CardColor_Red, CardColor_Blue),
			subCases: map[string]TestSubCase{
				"none":  {cards: "G4 Y4 Y4 G4", wantCompleted: false},
				"one":   {cards: "G2 Y7 Y4 B4", wantCompleted: false},
				"two":   {cards: "G5 R2 Y3 B4", wantCompleted: false},
				"three": {cards: "R7 Y5 R2 B5", wantCompleted: false},
				"four":  {cards: "R7 B7 B7 B5", wantCompleted: true},
				"red":   {cards: "R7 R7 R7 R5", wantCompleted: true},
				"blue":  {cards: "B7 B7 B7 B5", wantCompleted: true},
			},
		},
		"4-in-a row": {
			mission: NewFourValuesInARowMission(),
			subCases: map[string]TestSubCase{
				"different":  {cards: "G7 R2 Y1 B4", wantCompleted: false},
				"same-value": {cards: "B3 G5 Y4 B3", wantCompleted: false},
				"almost":     {cards: "B5 G3 Y1 B4", wantCompleted: false},
				"unordered":  {cards: "B5 G3 Y2 B4", wantCompleted: true},
				"ordered":    {cards: "B2 G3 Y4 B5", wantCompleted: true},
			},
		},
		"3-ordered": {
			mission: NewThreeOrderedValuesMission(),
			subCases: map[string]TestSubCase{
				"same":           {cards: "B3 G3 Y3 B3", wantCompleted: false},
				"asc-unordered":  {cards: "B1 G3 Y2 B7", wantCompleted: false},
				"desc-unordered": {cards: "B7 G5 Y3 B4", wantCompleted: false},
				"asc-1-3":        {cards: "B1 G2 Y3 B7", wantCompleted: true},
				"desc-1-3":       {cards: "B5 G4 Y3 B7", wantCompleted: true},
				"asc-2-4":        {cards: "B6 G2 Y3 B4", wantCompleted: true},
				"desc-2-4":       {cards: "B7 G4 Y3 B2", wantCompleted: true},
				"asc-1-4":        {cards: "B1 G2 Y3 B4", wantCompleted: true},
				"desc-1-4":       {cards: "B5 G4 Y3 B2", wantCompleted: true},
			},
		},
	}

	for name, tc := range testCases {
		for subName, subCase := range tc.subCases {
			t.Run(name+"_"+subName, func(t *testing.T) {
				topCards := TopCardsFromString(subCase.cards)
				gotCompleted := tc.mission.IsCompleted(topCards)
				require.Equal(t, subCase.wantCompleted, gotCompleted, name+": "+subName+": "+subCase.cards)
			})
		}
	}
}
