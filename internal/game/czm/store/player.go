package store

import (
	util_store "github.com/gre-ory/games-go/internal/util/store"

	"github.com/gre-ory/games-go/internal/game/czm/model"
)

type PlayerStore interface {
	Set(player *model.Player) error
	Get(id model.PlayerId) (*model.Player, error)
	Delete(id model.PlayerId) error
}

func NewPlayerStore() PlayerStore {
	return util_store.NewPlayerMemoryStore[model.PlayerId, *model.Player]()
}
