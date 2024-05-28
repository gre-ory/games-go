package service

import (
	"fmt"
	"sort"
	"time"

	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/util/loc"
	"github.com/gre-ory/games-go/internal/util/websocket"

	"github.com/gre-ory/games-go/internal/game/share/model"
)

type Player[PlayerIdT comparable, GameIdT comparable] interface {
	Id() PlayerIdT
	HasGameId() bool
	GameId() GameIdT
	GetLanguage() string
}

type Game[PlayerIdT comparable, GameIdT comparable, PlayerT Player[PlayerIdT, GameIdT]] interface {
	HasPlayer(playerId PlayerIdT) bool
	GetPlayer(playerId PlayerIdT) (PlayerT, error)
	GetCreatedAt() time.Time
	WrapData(data websocket.Data, player PlayerT) (bool, any)
}

type GameStore[GameIdT comparable, GameT any] interface {
	Get(gameId GameIdT) (GameT, error)
	Set(game GameT) error
	Delete(gameId GameIdT) error
}

type GameService[PlayerIdT comparable, GameIdT comparable, PlayerT Player[PlayerIdT, GameIdT], GameT Game[PlayerIdT, GameIdT, PlayerT]] interface {
	GetGame(gameId GameIdT) (GameT, error)
	StoreGame(game GameT) (GameT, error)
	DeleteGame(gameId GameIdT) error
	OnGame(gameId GameIdT, onGame func(game GameT) (GameT, error)) (GameT, error)
	OnPlayerGame(player PlayerT, onPlayerGame func(game GameT, player PlayerT) (GameT, error)) (GameT, error)
	SortGames(games []GameT) []GameT
	FilterGamesByPlayer(games []GameT, playerId PlayerIdT) []GameT
	WrapData(data websocket.Data, player PlayerT) (bool, any)
}

func NewGameService[PlayerIdT comparable, GameIdT comparable, PlayerT Player[PlayerIdT, GameIdT], GameT Game[PlayerIdT, GameIdT, PlayerT]](
	logger *zap.Logger,
	appId string,
	gameStore GameStore[GameIdT, GameT],
) GameService[PlayerIdT, GameIdT, PlayerT, GameT] {
	return &gameService[PlayerIdT, GameIdT, PlayerT, GameT]{
		logger:    logger,
		appId:     appId,
		gameStore: gameStore,
	}
}

type gameService[PlayerIdT comparable, GameIdT comparable, PlayerT Player[PlayerIdT, GameIdT], GameT Game[PlayerIdT, GameIdT, PlayerT]] struct {
	logger    *zap.Logger
	appId     string
	gameStore GameStore[GameIdT, GameT]
}

func (s *gameService[PlayerIdT, GameIdT, PlayerT, GameT]) GetGame(gameId GameIdT) (GameT, error) {
	return s.gameStore.Get(gameId)
}

func (s *gameService[PlayerIdT, GameIdT, PlayerT, GameT]) StoreGame(game GameT) (GameT, error) {
	err := s.gameStore.Set(game)
	if err != nil {
		var empty GameT
		return empty, err
	}
	return game, nil
}

func (s *gameService[PlayerIdT, GameIdT, PlayerT, GameT]) DeleteGame(gameId GameIdT) error {
	return s.gameStore.Delete(gameId)
}

func (s *gameService[PlayerIdT, GameIdT, PlayerT, GameT]) OnGame(gameId GameIdT, onGame func(game GameT) (GameT, error)) (GameT, error) {
	game, err := s.gameStore.Get(gameId)
	if err != nil {
		var empty GameT
		return empty, err
	}
	game, err = onGame(game)
	if err != nil {
		return game, err
	}
	return s.StoreGame(game)
}

func (s *gameService[PlayerIdT, GameIdT, PlayerT, GameT]) OnPlayerGame(player PlayerT, onPlayerGame func(game GameT, player PlayerT) (GameT, error)) (GameT, error) {
	if !player.HasGameId() {
		var empty GameT
		return empty, model.ErrMissingGameId
	}
	s.logger.Info("OnPlayerGame", zap.Any("game-id", player.GameId()), zap.Any("player-id", player.Id()))
	game, err := s.gameStore.Get(player.GameId())
	if err != nil {
		var empty GameT
		return empty, err
	}
	game, err = onPlayerGame(game, player)
	if err != nil {
		return game, err
	}
	return s.StoreGame(game)
}

func (s *gameService[PlayerIdT, GameIdT, PlayerT, GameT]) SortGames(games []GameT) []GameT {
	sort.Slice(games, func(i, j int) bool {
		// sort by reverse creation time
		return games[i].GetCreatedAt().After(games[j].GetCreatedAt())
	})
	return games
}

func (s *gameService[PlayerIdT, GameIdT, PlayerT, GameT]) FilterGamesByPlayer(games []GameT, playerId PlayerIdT) []GameT {
	filtered := make([]GameT, 0, len(games))
	for _, game := range games {
		if game.HasPlayer(playerId) {
			filtered = append(filtered, game)
		}
	}
	return filtered
}

func (s *gameService[PlayerIdT, GameIdT, PlayerT, GameT]) WrapData(data websocket.Data, player PlayerT) (bool, any) {
	language := player.GetLanguage()
	if language != "" {
		localizer := loc.NewLocalizer(s.appId, language, s.logger)
		s.logger.Info(fmt.Sprintf("[wrap] player %v: app=%s, lang=%s", player.Id(), s.appId, language))
		data.With("lang", localizer)
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
