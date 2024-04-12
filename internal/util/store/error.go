package store

import "fmt"

var (
	ErrGameNotFound   = fmt.Errorf("game not found")
	ErrPlayerNotFound = fmt.Errorf("player not found")
)
