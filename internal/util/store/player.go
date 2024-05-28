package store

type Player[IdT comparable] interface {
	Id() IdT
}

type PlayerStore[IdT comparable, PlayerT Player[IdT]] interface {
	Set(player PlayerT) error
	Get(id IdT) (PlayerT, error)
	Delete(id IdT) error
}

func NewPlayerMemoryStore[IdT comparable, PlayerT Player[IdT]]() PlayerStore[IdT, PlayerT] {
	return &playerMemoryStore[IdT, PlayerT]{
		players: map[IdT]PlayerT{},
	}
}

type playerMemoryStore[IdT comparable, PlayerT Player[IdT]] struct {
	players map[IdT]PlayerT
}

func (s *playerMemoryStore[IdT, PlayerT]) Set(player PlayerT) error {
	s.players[player.Id()] = player
	return nil
}

func (s *playerMemoryStore[IdT, PlayerT]) Get(id IdT) (PlayerT, error) {
	if player, ok := s.players[id]; ok {
		return player, nil
	}
	var empty PlayerT
	return empty, ErrPlayerNotFound
}

func (s *playerMemoryStore[IdT, PlayerT]) Delete(id IdT) error {
	if _, ok := s.players[id]; ok {
		delete(s.players, id)
		return nil
	}
	return ErrPlayerNotFound
}
