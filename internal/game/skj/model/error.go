package model

import "fmt"

var (
	ErrInvalidCardColor    = fmt.Errorf("invalid card color")
	ErrInvalidCard         = fmt.Errorf("invalid card")
	ErrEmptyCardDeck       = fmt.Errorf("empty card deck")
	ErrCardAlreadyFlipped  = fmt.Errorf("card already flipped")
	ErrInvalidNumberOfRow  = fmt.Errorf("invalid number of row")
	ErrInvalidRow          = fmt.Errorf("invalid row")
	ErrInvalidNumberOfCard = fmt.Errorf("invalid number of card")
	ErrInvalidColumn       = fmt.Errorf("invalid column")
)
