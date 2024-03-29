package store

import "github.com/gre-ory/games-go/internal/game/tictactoe/model"

type GameStore interface {
	ListNotStarted() []*model.Game
	Set(game *model.Game) error
	Get(id model.GameId) (*model.Game, error)
	Delete(id model.GameId) error
}

func NewGameStore() GameStore {
	return &gameStore{
		games: map[model.GameId]*model.Game{},
	}
}

type gameStore struct {
	games map[model.GameId]*model.Game
}

func (s *gameStore) ListNotStarted() []*model.Game {
	notStartedGames := make([]*model.Game, 0, len(s.games))
	for _, game := range s.games {
		if !game.Started() {
			notStartedGames = append(notStartedGames, game)
		}
	}
	return notStartedGames
}

func (s *gameStore) Set(game *model.Game) error {
	s.games[game.Id] = game
	return nil
}

func (s *gameStore) Get(id model.GameId) (*model.Game, error) {
	if game, ok := s.games[id]; ok {
		return game, nil
	}
	return nil, model.ErrGameNotFound
}

func (s *gameStore) Delete(id model.GameId) error {
	if _, ok := s.games[id]; ok {
		delete(s.games, id)
		return nil
	}
	return model.ErrGameNotFound
}