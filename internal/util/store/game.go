package store

type Game[IdT comparable, StatusT comparable] interface {
	GetId() IdT
	GetStatus() StatusT
}

type GameStore[IdT comparable, StatusT comparable, GameT Game[IdT, StatusT]] interface {
	ListStatus(status StatusT) []GameT
	Set(game GameT) error
	Get(id IdT) (GameT, error)
	Delete(id IdT) error
}

func NewGameMemoryStore[IdT comparable, StatusT comparable, GameT Game[IdT, StatusT]]() GameStore[IdT, StatusT, GameT] {
	return &gameMemoryStore[IdT, StatusT, GameT]{
		games: map[IdT]GameT{},
	}
}

type gameMemoryStore[IdT comparable, StatusT comparable, GameT Game[IdT, StatusT]] struct {
	games map[IdT]GameT
}

func (s *gameMemoryStore[IdT, StatusT, GameT]) ListStatus(status StatusT) []GameT {
	filtered := make([]GameT, 0, len(s.games))
	for _, game := range s.games {
		if game.GetStatus() == status {
			filtered = append(filtered, game)
		}
	}
	return filtered
}

func (s *gameMemoryStore[IdT, StatusT, GameT]) Set(game GameT) error {
	s.games[game.GetId()] = game
	return nil
}

func (s *gameMemoryStore[IdT, StatusT, GameT]) Get(id IdT) (GameT, error) {
	if game, ok := s.games[id]; ok {
		return game, nil
	}
	var empty GameT
	return empty, ErrGameNotFound
}

func (s *gameMemoryStore[IdT, StatusT, GameT]) Delete(id IdT) error {
	if _, ok := s.games[id]; ok {
		delete(s.games, id)
		return nil
	}
	return ErrGameNotFound
}
