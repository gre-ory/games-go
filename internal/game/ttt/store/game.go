package store

import (
	util_store "github.com/gre-ory/games-go/internal/util/store"

	"github.com/gre-ory/games-go/internal/game/ttt/model"
)

type GameStore interface {
	ListStatus(status model.GameStatus) []*model.Game
	Set(game *model.Game) error
	Get(id model.GameId) (*model.Game, error)
	Delete(id model.GameId) error
}

func NewGameStore() GameStore {
	return util_store.NewGameMemoryStore[model.GameId, model.GameStatus, *model.Game]()
}
