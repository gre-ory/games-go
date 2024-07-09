package store

import (
	share_store "github.com/gre-ory/games-go/internal/game/share/store"

	"github.com/gre-ory/games-go/internal/game/czm/model"
)

type PlayerStore interface {
	share_store.PlayerStore[*model.Player]
}

func NewPlayerStore() PlayerStore {
	return share_store.NewPlayerMemoryStore[*model.Player]()
}