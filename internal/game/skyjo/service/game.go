package service

import (
	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/util/dict"
	"github.com/gre-ory/games-go/internal/util/list"
	"github.com/gre-ory/games-go/internal/util/websocket"

	share_service "github.com/gre-ory/games-go/internal/game/share/service"

	"github.com/gre-ory/games-go/internal/game/skyjo/model"
	"github.com/gre-ory/games-go/internal/game/skyjo/store"
)

type GameService interface {
	GetJoinableGames() []*model.Game
	GetNotJoinableGames(playerId model.PlayerId) []*model.Game
	GetGame(id model.GameId) (*model.Game, error)
	NewGame() (*model.Game, error)
	JoinGame(id model.GameId, player *model.Player) (*model.Game, error)
	StartGame(player *model.Player) (*model.Game, error)
	PlayGame(player *model.Player, x, y int) (*model.Game, error)
	LeaveGame(player *model.Player) (*model.Game, error)
	DeleteGame(id model.GameId, playerId model.PlayerId) error
	WrapData(data websocket.Data, player *model.Player) (bool, any)
}

func NewGameService(logger *zap.Logger, gameStore store.GameStore, playerStore store.PlayerStore) GameService {
	return &gameService{
		GameService: share_service.NewGameService(logger, gameStore),
		logger:      logger,
		gameStore:   gameStore,
		playerStore: playerStore,
	}
}

type gameService struct {
	share_service.GameService[model.PlayerId, model.GameId, *model.Player, *model.Game]
	logger      *zap.Logger
	gameStore   store.GameStore
	playerStore store.PlayerStore
}

func (s *gameService) GetJoinableGames() []*model.Game {
	games := s.gameStore.ListStatus(model.Joinable)
	return s.SortGames(games)
}

func (s *gameService) GetNotJoinableGames(playerId model.PlayerId) []*model.Game {
	games := make([]*model.Game, 0)
	games = append(games, s.gameStore.ListStatus(model.NotJoinable)...)
	games = append(games, s.gameStore.ListStatus(model.Started)...)
	games = append(games, s.gameStore.ListStatus(model.Stopped)...)
	games = s.FilterGamesByPlayer(games, playerId)
	return s.SortGames(games)
}

func (s *gameService) NewGame() (*model.Game, error) {
	game := model.NewGame(3, 3)
	return s.StoreGame(game)
}

func (s *gameService) JoinGame(id model.GameId, player *model.Player) (*model.Game, error) {
	return s.OnPlayerGame(player, s.joinGame)
}

func (s *gameService) joinGame(game *model.Game, player *model.Player) (*model.Game, error) {
	switch game.Status {
	case model.NotJoinable:
		return nil, model.ErrGameNotJoinable
	case model.Started:
		return nil, model.ErrGameAlreadyStarted
	case model.Stopped:
		return nil, model.ErrGameStopped
	default:
	}

	if _, err := game.GetPlayer(player.GetId()); err == nil {
		return game, nil
	}
	game = game.WithPlayer(player)
	game.UpdateStatus()
	return game, nil
}

func (s *gameService) StartGame(player *model.Player) (*model.Game, error) {
	return s.OnPlayerGame(player, s.startGame)
}

func (s *gameService) startGame(game *model.Game, player *model.Player) (*model.Game, error) {
	switch game.Status {
	case model.Started:
		return nil, model.ErrGameAlreadyStarted
	case model.Stopped:
		return nil, model.ErrGameStopped
	default:
	}
	if !game.CanStart() {
		return nil, model.ErrMissingPlayers
	}

	ids := dict.ConvertToList(game.Players, dict.Key)
	list.Shuffle(ids)
	game.PlayerIds = ids

	game.Players[ids[0]].WithSymbol(model.PLAYER_ONE_SYMBOL)
	game.Players[ids[1]].WithSymbol(model.PLAYER_TWO_SYMBOL)

	game.Status = model.Started
	game.Round = 1
	game.SetPlayingPlayer()

	return game, nil
}

func (s *gameService) PlayGame(player *model.Player, x, y int) (*model.Game, error) {
	return s.OnPlayerGame(player, func(game *model.Game, player *model.Player) (*model.Game, error) {
		return s.playGame(game, player, x, y)
	})
}

func (s *gameService) playGame(game *model.Game, player *model.Player, x, y int) (*model.Game, error) {
	switch game.Status {
	case model.Joinable, model.NotJoinable:
		return nil, model.ErrGameNotStarted
	case model.Stopped:
		return nil, model.ErrGameStopped
	default:
	}

	currentPlayer, err := game.GetCurrentPlayer()
	if err != nil {
		return nil, err
	}
	if player.GetId() != currentPlayer.GetId() {
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
		game.Round++
		game.SetPlayingPlayer()
	}

	return game, nil
}

func (s *gameService) LeaveGame(player *model.Player) (*model.Game, error) {
	return s.OnPlayerGame(player, s.leaveGame)
}

func (s *gameService) leaveGame(game *model.Game, player *model.Player) (*model.Game, error) {
	switch game.Status {
	case model.Joinable, model.NotJoinable:
		game = game.WithoutPlayer(player)
		player.Status = model.WaitingToJoin
		if len(game.Players) == 0 {
			return nil, s.deleteGame(game)
		} else {
			game.UpdateStatus()

			return game, nil
		}
	case model.Started:
		winnerId, err := game.GetOtherPlayerId(player.GetId())
		if err != nil {
			return nil, err
		}
		player.Status = model.WaitingToJoin
		player.UnsetGameId()
		return s.stopGame(game, winnerId)
	case model.Stopped:
		player.Status = model.WaitingToJoin
		player.UnsetGameId()
		return game, nil
	}

	return game, nil
}

func (s *gameService) stopGame(game *model.Game, winnerId model.PlayerId) (*model.Game, error) {
	switch game.Status {
	case model.Joinable, model.NotJoinable:
		return nil, model.ErrGameNotStarted
	case model.Stopped:
		return nil, model.ErrGameStopped
	default:
	}

	game.Status = model.Stopped
	if winnerId != "" {
		game.WinnerIds = []model.PlayerId{winnerId}
	}

	return game, nil
}

func (s *gameService) DeleteGame(id model.GameId, playerId model.PlayerId) error {
	game, err := s.GetGame(id)
	if err != nil {
		return err
	}
	if _, err := game.GetPlayer(playerId); err != nil {
		return err
	}
	return s.deleteGame(game)
}

func (s *gameService) deleteGame(game *model.Game) error {
	switch game.Status {
	case model.Started:
		return model.ErrGameNotStopped
	default:
	}
	return s.GameService.DeleteGame(game.Id)
}
