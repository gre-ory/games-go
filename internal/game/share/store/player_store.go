package store

import "github.com/gre-ory/games-go/internal/game/share/model"

// //////////////////////////////////////////////////
// player store

type PlayerStorable interface {
	Id() model.PlayerId
	Status() model.PlayerStatus
}

type PlayerStore[PlayerT PlayerStorable] interface {
	Set(player PlayerT) error
	Get(id model.PlayerId) (PlayerT, error)
	Delete(id model.PlayerId) error
}

// //////////////////////////////////////////////////
// player memory store

func NewPlayerMemoryStore[PlayerT PlayerStorable]() PlayerStore[PlayerT] {
	return &playerMemoryStore[PlayerT]{
		players: map[model.PlayerId]PlayerT{},
	}
}

type playerMemoryStore[PlayerT PlayerStorable] struct {
	players map[model.PlayerId]PlayerT
}

func (s *playerMemoryStore[PlayerT]) Set(player PlayerT) error {
	s.players[player.Id()] = player
	return nil
}

func (s *playerMemoryStore[PlayerT]) Get(id model.PlayerId) (PlayerT, error) {
	if player, ok := s.players[id]; ok {
		return player, nil
	}
	var empty PlayerT
	return empty, model.ErrPlayerNotFound
}

func (s *playerMemoryStore[PlayerT]) Delete(id model.PlayerId) error {
	if _, ok := s.players[id]; ok {
		delete(s.players, id)
		return nil
	}
	return model.ErrPlayerNotFound
}
