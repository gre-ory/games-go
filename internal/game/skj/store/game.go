package store

import (
	share_store "github.com/gre-ory/games-go/internal/game/share/store"

	"github.com/gre-ory/games-go/internal/game/skj/model"
)

type GameStore interface {
	share_store.GameStore[*model.Game]
}

func NewGameStore() GameStore {
	return share_store.NewGameMemoryStore[*model.Game]()
}
