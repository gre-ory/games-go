package service

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/util/dict"
	"github.com/gre-ory/games-go/internal/util/list"
	"github.com/gre-ory/games-go/internal/util/websocket"

	share_service "github.com/gre-ory/games-go/internal/game/share/service"

	"github.com/gre-ory/games-go/internal/game/czm/model"
	"github.com/gre-ory/games-go/internal/game/czm/store"
)

type GameService interface {
	GetJoinableGames() []*model.Game
	GetNotJoinableGames(playerId model.PlayerId) []*model.Game
	GetGame(id model.GameId) (*model.Game, error)
	CreateGame(player *model.Player) (*model.Game, error)
	JoinGame(id model.GameId, player *model.Player) (*model.Game, error)
	StartGame(player *model.Player) (*model.Game, error)
	SelectCard(player *model.Player, cardIndex int) (*model.Game, error)
	PlayCard(player *model.Player, cardIndex int, discardIndex int) (*model.Game, error)
	LeaveGame(player *model.Player) (*model.Game, error)
	DeleteGame(id model.GameId, playerId model.PlayerId) error
	WrapData(data websocket.Data, player *model.Player) (bool, any)
}

func NewGameService(logger *zap.Logger, gameStore store.GameStore, playerStore store.PlayerStore) GameService {
	return &gameService{
		GameService: share_service.NewGameService(logger, model.AppId, gameStore),
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
	game := model.NewGame()
	return s.StoreGame(game)
}

func (s *gameService) CreateGame(player *model.Player) (*model.Game, error) {
	game, err := s.NewGame()
	if err != nil {
		return nil, err
	}
	return s.joinGame(game, player)
}

func (s *gameService) JoinGame(id model.GameId, player *model.Player) (*model.Game, error) {
	return s.OnGame(id, func(game *model.Game) (*model.Game, error) {
		return s.joinGame(game, player)
	})
}

func (s *gameService) joinGame(game *model.Game, player *model.Player) (*model.Game, error) {
	switch game.Status() {
	case model.NotJoinable:
		return nil, model.ErrGameNotJoinable
	case model.Started:
		return nil, model.ErrGameAlreadyStarted
	case model.Stopped:
		return nil, model.ErrGameStopped
	default:
	}

	if _, err := game.GetPlayer(player.Id()); err == nil {
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
	switch game.Status() {
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

	for i := 0; i < 4; i++ {
		for _, player := range game.Players {
			card, err := game.DrawCardDeck.Draw()
			if err != nil {
				return nil, err
			}
			player.WithCard(card)
		}

		card, err := game.DrawCardDeck.Draw()
		if err != nil {
			return nil, err
		}
		game.DiscardCardDecks[i].Add(card)
	}

	game.SetStatus(model.Started)
	game.Round = 1
	game.SetPlayingPlayer()

	return game, nil
}

func (s *gameService) SelectCard(player *model.Player, cardIndex int) (*model.Game, error) {
	return s.OnPlayerGame(player, func(game *model.Game, player *model.Player) (*model.Game, error) {
		return s.selectCard(game, player, cardIndex)
	})
}

func (s *gameService) selectCard(game *model.Game, player *model.Player, cardIndex int) (*model.Game, error) {
	switch game.Status() {
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
	if player.Id() != currentPlayer.Id() {
		return nil, model.ErrWrongPlayer
	}

	_, err = player.SelectCard(cardIndex)
	if err != nil {
		return nil, err
	}

	game.SelectedCardIndex = cardIndex

	return game, nil
}

func (s *gameService) PlayCard(player *model.Player, cardIndex int, discardIndex int) (*model.Game, error) {
	return s.OnPlayerGame(player, func(game *model.Game, player *model.Player) (*model.Game, error) {
		return s.playCard(game, player, cardIndex, discardIndex)
	})
}

func (s *gameService) playCard(game *model.Game, player *model.Player, cardIndex int, discardIndex int) (*model.Game, error) {
	switch game.Status() {
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
	if player.Id() != currentPlayer.Id() {
		return nil, model.ErrWrongPlayer
	}

	card, err := player.PlayCard(cardIndex)
	if err != nil {
		return nil, err
	}

	if discardIndex < 0 || model.NbCardDeck <= discardIndex {
		return nil, model.ErrInvalidDiscardIndex
	}

	game.SelectedCardIndex = -1
	game.DiscardCardDecks[discardIndex].Add(card)

	newCard, err := game.DrawCardDeck.Draw()
	if err != nil {
		return nil, err
	}
	player.WithCard(newCard)

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

func (s *gameService) ValidateMissions(player *model.Player, cardIndex int, discardIndex int) (*model.Game, error) {
	return s.OnPlayerGame(player, func(game *model.Game, player *model.Player) (*model.Game, error) {
		return s.validateMissions(game, player)
	})
}

func (s *gameService) validateMissions(game *model.Game, player *model.Player) (*model.Game, error) {
	switch game.Status() {
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
	if player.Id() != currentPlayer.Id() {
		return nil, model.ErrWrongPlayer
	}

	topCards := game.GetTopCards()

	game.ValidatedMissionIndex = -1
	for index, mission := range game.Missions {
		if mission != nil {
			continue
		}
		if mission.IsCompleted(topCards) {
			game.ValidatedMissionIndex = index
			game.DiscardMissionDeck.Add(mission)
			// draw new mission
			newMission, err := game.DrawMissionDeck.Draw()
			if errors.Is(err, model.ErrEmptyMissionDeck) {
				// TODO win
				return nil, err
			}
			game.Missions[index] = newMission
			break
		}
	}
	return game, nil
}

func (s *gameService) LeaveGame(player *model.Player) (*model.Game, error) {
	return s.OnPlayerGame(player, s.leaveGame)
}

func (s *gameService) leaveGame(game *model.Game, player *model.Player) (*model.Game, error) {
	s.logger.Info("leaveGame", zap.Any("game", game), zap.Any("player", player))
	switch game.Status() {
	case model.Joinable, model.NotJoinable:
		game = game.WithoutPlayer(player)
		player.Status = model.WaitingToJoin
		if len(game.Players) == 0 {
			s.logger.Info(" -> delete game", zap.Any("game", game), zap.Any("player", player))
			return nil, s.deleteGame(game)
		} else {
			game.UpdateStatus()

			return game, nil
		}
	case model.Started:
		winnerId, err := game.GetOtherPlayerId(player.Id())
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
	switch game.Status() {
	case model.Joinable, model.NotJoinable:
		return nil, model.ErrGameNotStarted
	case model.Stopped:
		return nil, model.ErrGameStopped
	default:
	}

	game.SetStatus(model.Stopped)
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
	switch game.Status() {
	case model.Started:
		return model.ErrGameNotStopped
	default:
	}
	return s.gameStore.Delete(game.Id())
}
