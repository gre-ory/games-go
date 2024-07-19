package store

import "github.com/gre-ory/games-go/internal/game/share/model"

// //////////////////////////////////////////////////
// game store

type GameStorable interface {
	Id() model.GameId
	Status() model.GameStatus
}

type GameStore[GameT GameStorable] interface {
	ListStatus(status model.GameStatus) []GameT
	Set(game GameT) error
	Get(id model.GameId) (GameT, error)
	Delete(id model.GameId) error
}

// //////////////////////////////////////////////////
// game memory store

func NewGameMemoryStore[GameT GameStorable]() GameStore[GameT] {
	return &gameMemoryStore[GameT]{
		games: map[model.GameId]GameT{},
	}
}

type gameMemoryStore[GameT GameStorable] struct {
	games map[model.GameId]GameT
	empty GameT
}

func (s *gameMemoryStore[GameT]) ListStatus(status model.GameStatus) []GameT {
	filtered := make([]GameT, 0, len(s.games))
	for _, game := range s.games {
		if game.Status() == status {
			filtered = append(filtered, game)
		}
	}
	return filtered
}

func (s *gameMemoryStore[GameT]) Set(game GameT) error {
	s.games[game.Id()] = game
	return nil
}

func (s *gameMemoryStore[GameT]) Get(id model.GameId) (GameT, error) {
	if game, ok := s.games[id]; ok {
		return game, nil
	}
	return s.empty, model.ErrGameNotFound
}

func (s *gameMemoryStore[GameT]) Delete(id model.GameId) error {
	if _, ok := s.games[id]; ok {
		delete(s.games, id)
		return nil
	}
	return model.ErrGameNotFound
}
