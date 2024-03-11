package store

import (
	"github.com/gre-ory/games-go/internal/game/tictactoe/model"
)

type PlayerStore interface {
	Set(player *model.Player) error
	Get(id model.PlayerId) (*model.Player, error)
	Delete(id model.PlayerId) error
}

func NewPlayerStore() PlayerStore {
	return &playerStore{
		players: map[model.PlayerId]*model.Player{},
	}
}

type playerStore struct {
	players map[model.PlayerId]*model.Player
}

func (s *playerStore) Set(player *model.Player) error {
	s.players[player.Id()] = player
	return nil
}

func (s *playerStore) Get(id model.PlayerId) (*model.Player, error) {
	if player, ok := s.players[id]; ok {
		return player, nil
	}
	return nil, model.ErrPlayerNotFound
}

func (s *playerStore) Delete(id model.PlayerId) error {
	if _, ok := s.players[id]; ok {
		delete(s.players, id)
		return nil
	}
	return model.ErrPlayerNotFound
}
