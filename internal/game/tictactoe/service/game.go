package service

import (
	"github.com/gre-ory/games-go/internal/game/tictactoe/model"
	"github.com/gre-ory/games-go/internal/game/tictactoe/store"
	"github.com/gre-ory/games-go/internal/util/dict"
	"github.com/gre-ory/games-go/internal/util/list"
	"github.com/gre-ory/games-go/internal/util/websocket"
)

type GameService interface {
	GetJoinableGames() []*model.Game
	GetGame(id model.GameId) (*model.Game, error)
	NewGame() (*model.Game, error)
	JoinGame(id model.GameId, player *model.Player) (*model.Game, error)
	StartGame(player *model.Player) (*model.Game, error)
	PlayGame(player *model.Player, x, y int) (*model.Game, error)
	LeaveGame(player *model.Player) (*model.Game, error)
	DeleteGame(id model.GameId, playerId model.PlayerId) error
	WrapData(data websocket.Data, player *model.Player) (bool, any)
}

func NewGameService(gameStore store.GameStore, playerStore store.PlayerStore) GameService {
	return &gameService{
		gameStore:   gameStore,
		playerStore: playerStore,
	}
}

type gameService struct {
	gameStore   store.GameStore
	playerStore store.PlayerStore
}

func (s *gameService) GetGame(gameId model.GameId) (*model.Game, error) {
	game, err := s.gameStore.Get(gameId)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (s *gameService) GetJoinableGames() []*model.Game {
	return s.gameStore.ListNotStarted()
}

func (s *gameService) NewGame() (*model.Game, error) {
	game := model.NewGame(3, 3)
	return s.storeGame(game)
}

func (s *gameService) JoinGame(id model.GameId, player *model.Player) (*model.Game, error) {
	game, err := s.gameStore.Get(id)
	if err != nil {
		return nil, err
	}
	return s.joinGame(game, player)
}

func (s *gameService) joinGame(game *model.Game, player *model.Player) (*model.Game, error) {
	if game.Stopped {
		return nil, model.ErrGameStopped
	}
	if _, err := game.GetPlayer(player.Id()); err == nil {
		return game, nil
	}
	game = game.WithPlayer(player)

	if game.CanStart() {
		for _, player := range game.Players {
			player.Status = model.WaitingToStart
		}
	} else {
		for _, player := range game.Players {
			player.Status = model.WaitingToJoin
		}
	}

	return s.storeGame(game)
}

func (s *gameService) StartGame(player *model.Player) (*model.Game, error) {
	game, err := s.gameStore.Get(player.GameId())
	if err != nil {
		return nil, err
	}
	return s.startGame(game)
}

func (s *gameService) startGame(game *model.Game) (*model.Game, error) {
	if game.Stopped {
		return nil, model.ErrGameStopped
	}
	if !game.CanStart() {
		return nil, model.ErrMissingPlayers
	}

	ids := dict.ConvertToList(game.Players, dict.Key)
	list.Shuffle(ids)
	game.PlayerIds = ids

	game.Players[ids[0]].WithSymbol(model.PLAYER_ONE_SYMBOL)
	game.Players[ids[1]].WithSymbol(model.PLAYER_TWO_SYMBOL)

	game.NextRound()

	return s.storeGame(game)
}

func (s *gameService) PlayGame(player *model.Player, x, y int) (*model.Game, error) {
	game, err := s.gameStore.Get(player.GameId())
	if err != nil {
		return nil, err
	}
	return s.playGame(game, player, x, y)
}

func (s *gameService) playGame(game *model.Game, player *model.Player, x, y int) (*model.Game, error) {
	if !game.Started() {
		return nil, model.ErrGameNotStarted
	}
	if game.Stopped {
		return nil, model.ErrGameStopped
	}

	currentPlayer, err := game.GetCurrentPlayer()
	if err != nil {
		return nil, err
	}
	if player.Id() != currentPlayer.Id() {
		return nil, model.ErrWrongPlayer
	}

	err = game.Play(player, x, y)
	if err != nil {
		return nil, err
	}

	if yes, winnerId := game.HasWinner(); yes {
		s.stopGame(game, winnerId)
	} else if game.IsTie() {
		s.stopGame(game, "")
	} else {
		game.NextRound()
	}

	return s.storeGame(game)
}

func (s *gameService) LeaveGame(player *model.Player) (*model.Game, error) {
	game, err := s.gameStore.Get(player.GameId())
	if err != nil {
		return nil, err
	}
	return s.leaveGame(game, player)
}

func (s *gameService) leaveGame(game *model.Game, player *model.Player) (*model.Game, error) {
	if game.Stopped {
		player.UnsetGameId()
		return game, nil
	}
	if !game.Started() {
		game = game.WithoutPlayer(player)
		if len(game.Players) == 0 {
			return nil, s.deleteGame(game)
		} else {

			if game.CanStart() {
				for _, player := range game.Players {
					player.Status = model.WaitingToStart
				}
			} else {
				for _, player := range game.Players {
					player.Status = model.WaitingToJoin
				}
			}

			return s.storeGame(game)
		}
	}
	winnerId := game.PlayerIds[0]
	if game.PlayerIds[0] == player.Id() {
		winnerId = game.PlayerIds[1]
	}
	return s.stopGame(game, winnerId)
}

func (s *gameService) stopGame(game *model.Game, winnerId model.PlayerId) (*model.Game, error) {
	if game.Stopped {
		return nil, model.ErrGameStopped
	}
	if !game.Started() {
		return nil, model.ErrGameNotStarted
	}

	for id, player := range game.Players {
		if winnerId == "" {
			player.Status = model.Tie
		} else if id == winnerId {
			player.Status = model.Win
		} else {
			player.Status = model.Loose
		}
	}
	game.Stopped = true

	return s.storeGame(game)
}

func (s *gameService) storeGame(game *model.Game) (*model.Game, error) {
	err := s.gameStore.Set(game)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (s *gameService) DeleteGame(id model.GameId, playerId model.PlayerId) error {
	game, err := s.gameStore.Get(id)
	if err != nil {
		return err
	}
	if _, err := game.GetPlayer(playerId); err != nil {
		return err
	}
	return s.deleteGame(game)
}

func (s *gameService) deleteGame(game *model.Game) error {
	if game.Started() && !game.Stopped {
		return model.ErrGameNotStopped
	}
	return s.gameStore.Delete(game.Id)
}

func (s *gameService) WrapData(data websocket.Data, player *model.Player) (bool, any) {
	if player == nil {
		return true, data
	}
	if !player.HasGameId() {
		return true, data
	}
	game, err := s.gameStore.Get(player.GameId())
	if err != nil {
		return false, nil
	}
	return game.WrapData(data, player)
}
