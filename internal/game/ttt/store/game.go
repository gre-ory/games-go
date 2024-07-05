package store

import (
	"github.com/gre-ory/games-go/internal/game/ttt/model"
)

type GameStore interface {
	game.GameStore[*model.Game]
}

func NewGameStore() GameStore {
	return game.NewGameMemoryStore[*model.Game]()
}
