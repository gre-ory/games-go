package model

import "fmt"

var (
	ErrOutOfColumnBound  = fmt.Errorf("out of column bound")
	ErrOutOfRowBound     = fmt.Errorf("out of row bound")
	ErrMissingPlayX      = fmt.Errorf("missing play x")
	ErrInvalidPlayX      = fmt.Errorf("invalid play x")
	ErrMissingPlayY      = fmt.Errorf("missing play y")
	ErrInvalidPlayY      = fmt.Errorf("invalid play y")
	ErrAlreadyPlayOnCell = fmt.Errorf("already played on cell")
)
