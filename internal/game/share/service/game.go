package service

import (
	"sort"

	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/game/share/model"
	"github.com/gre-ory/games-go/internal/game/share/store"
)

// //////////////////////////////////////////////////
// game service

type GameService[PlayerT model.Player, GameT model.Game[PlayerT]] interface {
	GetGame(gameId model.GameId) (GameT, error)
	GetJoinableGames() []GameT
	GetNonJoinableGames(playerId model.PlayerId) []GameT
	SortGamesByCreationTime(games []GameT) []GameT
	FilterGamesByPlayer(games []GameT, playerId model.PlayerId) []GameT

	CreateGame(player PlayerT) (GameT, error)
	JoinGameId(gameId model.GameId, player PlayerT) (GameT, error)
	JoinGame(game GameT, player PlayerT) (GameT, error)
	StartPlayerGame(player PlayerT) (GameT, error)
	StartGame(game GameT) (GameT, error)
	LeavePlayerGame(player PlayerT) (GameT, error)
	LeaveGame(game GameT, player PlayerT) (GameT, error)
	StopGame(game GameT) (GameT, error)
	DeleteGameId(gameId model.GameId, playerId model.PlayerId) error
	DeleteGame(game GameT, playerId model.PlayerId) error

	SaveGame(game GameT) (GameT, error)

	RegisterOnJoinGame(func(game GameT, player PlayerT))
	RegisterOnGame(func(game GameT))
	RegisterOnLeaveGame(func(game GameT, player PlayerT))
}

// //////////////////////////////////////////////////
// game plugin

type GamePlugin[PlayerT model.Player, GameT model.Game[PlayerT]] interface {
	CanCreateGame(player PlayerT) error
	CreateGame(player PlayerT) (GameT, error)

	CanJoinGame(game GameT, player PlayerT) error
	JoinGame(game GameT, player PlayerT) (GameT, error)

	CanStartGame(game GameT) error
	StartGame(game GameT) (GameT, error)

	CanStopGame(game GameT) error
	StopGame(game GameT) (GameT, error)

	CanLeaveGame(game GameT, player PlayerT) error
	LeaveGame(game GameT, player PlayerT) (GameT, error)

	CanDeleteGame(game GameT, playerId model.PlayerId) error
}

func NewGameService[PlayerT model.Player, GameT model.Game[PlayerT]](logger *zap.Logger, plugin GamePlugin[PlayerT, GameT], gameStore store.GameStore[GameT]) GameService[PlayerT, GameT] {
	return &gameService[PlayerT, GameT]{
		logger:    logger,
		plugin:    plugin,
		gameStore: gameStore,
	}
}

type gameService[PlayerT model.Player, GameT model.Game[PlayerT]] struct {
	logger     *zap.Logger
	plugin     GamePlugin[PlayerT, GameT]
	gameStore  store.GameStore[GameT]
	onJoinFns  []func(game GameT, player PlayerT)
	onGameFns  []func(game GameT)
	onLeaveFns []func(game GameT, player PlayerT)
	empty      GameT
}

// //////////////////////////////////////////////////
// get game

func (s *gameService[PlayerT, GameT]) GetGame(id model.GameId) (GameT, error) {
	game, err := s.gameStore.Get(id)
	if err != nil {
		return s.empty, err
	}
	return game, nil
}

// //////////////////////////////////////////////////
// get joinable games

func (s *gameService[PlayerT, GameT]) GetJoinableGames() []GameT {
	games := make([]GameT, 0)
	games = append(games, s.gameStore.ListStatus(model.GameStatus_JoinableNotStartable)...)
	games = append(games, s.gameStore.ListStatus(model.GameStatus_JoinableAndStartable)...)
	return s.SortGamesByCreationTime(games)
}

// //////////////////////////////////////////////////
// get non-joinable games

func (s *gameService[PlayerT, GameT]) GetNonJoinableGames(playerId model.PlayerId) []GameT {
	games := make([]GameT, 0)
	games = append(games, s.gameStore.ListStatus(model.GameStatus_NotJoinableAndStartable)...)
	games = append(games, s.gameStore.ListStatus(model.GameStatus_Started)...)
	games = append(games, s.gameStore.ListStatus(model.GameStatus_Stopped)...)
	games = s.FilterGamesByPlayer(games, playerId)
	return s.SortGamesByCreationTime(games)
}

func (s *gameService[PlayerT, GameT]) SortGamesByCreationTime(games []GameT) []GameT {
	sort.Slice(games, func(i, j int) bool {
		// sort by reverse creation time
		return games[i].CreatedAt().After(games[j].CreatedAt())
	})
	return games
}

func (s *gameService[PlayerT, GameT]) FilterGamesByPlayer(games []GameT, playerId model.PlayerId) []GameT {
	filtered := make([]GameT, 0, len(games))
	for _, game := range games {
		if game.HasPlayer(playerId) {
			filtered = append(filtered, game)
		}
	}
	return filtered
}

// //////////////////////////////////////////////////
// create game

func (s *gameService[PlayerT, GameT]) CreateGame(player PlayerT) (GameT, error) {

	//
	// preliminary checks
	//

	if err := s.plugin.CanCreateGame(player); err != nil {
		return s.empty, err
	}

	//
	// join game
	//

	game, err := s.plugin.CreateGame(player)
	if err != nil {
		return s.empty, err
	}

	game.AttachPlayer(player)
	game.UpdateJoinStatus()

	//
	// save game
	//

	game, err = s.SaveGame(game)
	if err != nil {
		return s.empty, err
	}

	//
	// callbacks
	//

	s.onJoinGame(game, player)

	return game, nil
}

// //////////////////////////////////////////////////
// join game

func (s *gameService[PlayerT, GameT]) JoinGameId(id model.GameId, player PlayerT) (GameT, error) {
	game, err := s.gameStore.Get(id)
	if err != nil {
		return s.empty, err
	}
	return s.JoinGame(game, player)
}

func (s *gameService[PlayerT, GameT]) JoinGame(game GameT, player PlayerT) (GameT, error) {

	//
	// check status
	//

	if err := game.Status().CanJoin(); err != nil {
		return s.empty, err
	}

	//
	// check player
	//

	if ok := game.HasPlayer(player.Id()); ok {
		return game, nil
	}

	//
	// preliminary checks
	//

	if err := s.plugin.CanJoinGame(game, player); err != nil {
		return s.empty, err
	}

	//
	// join game
	//

	game, err := s.plugin.JoinGame(game, player)
	if err != nil {
		return s.empty, err
	}
	if !player.HasGameId() {
		game.AttachPlayer(player)
	}
	game.UpdateJoinStatus()

	//
	// save game
	//

	game, err = s.SaveGame(game)
	if err != nil {
		return s.empty, err
	}

	//
	// callbacks
	//

	s.onJoinGame(game, player)

	return game, nil
}

// //////////////////////////////////////////////////
// start game

func (s *gameService[PlayerT, GameT]) StartPlayerGame(player PlayerT) (GameT, error) {
	if !player.HasGameId() {
		return s.empty, model.ErrPlayerNotInGame
	}
	game, err := s.gameStore.Get(player.GameId())
	if err != nil {
		return s.empty, err
	}
	return s.StartGame(game)
}

func (s *gameService[PlayerT, GameT]) StartGame(game GameT) (GameT, error) {

	//
	// preliminary checks
	//

	if err := game.Status().CanStart(); err != nil {
		return s.empty, err
	}
	if err := s.plugin.CanStartGame(game); err != nil {
		return s.empty, err
	}

	//
	// start game
	//

	game, err := s.plugin.StartGame(game)
	if err != nil {
		return s.empty, err
	}
	if !game.IsStarted() {
		game.SetStarted()
	}

	//
	// save game
	//

	game, err = s.SaveGame(game)
	if err != nil {
		return s.empty, err
	}

	//
	// callbacks
	//

	s.onGame(game)

	return game, nil
}

// //////////////////////////////////////////////////
// leave game

func (s *gameService[PlayerT, GameT]) LeavePlayerGame(player PlayerT) (GameT, error) {
	game, err := s.gameStore.Get(player.GameId())
	if err != nil {
		return s.empty, err
	}
	return s.LeaveGame(game, player)
}

func (s *gameService[PlayerT, GameT]) LeaveGame(game GameT, player PlayerT) (GameT, error) {

	//
	// check status
	//

	if err := game.Status().CanLeave(); err != nil {
		return s.empty, err
	}

	//
	// check player
	//

	if ok := game.HasPlayer(player.Id()); !ok {
		return s.empty, model.ErrPlayerNotInGame
	}

	//
	// preliminary checks
	//

	if err := s.plugin.CanLeaveGame(game, player); err != nil {
		return s.empty, err
	}

	//
	// leave game
	//

	game, err := s.plugin.LeaveGame(game, player)
	if err != nil {
		return s.empty, err
	}
	if player.HasGameId() {
		game.UpdateJoinStatus()
	}
	game.UpdateJoinStatus()
	player.SetStatus(model.PlayerStatus_WaitingToJoin)

	//
	// save game
	//

	game, err = s.SaveGame(game)
	if err != nil {
		return s.empty, err
	}

	//
	// callbacks
	//

	s.onLeaveGame(game, player)

	return game, nil
}

// //////////////////////////////////////////////////
// stop game

func (s *gameService[PlayerT, GameT]) StopGame(game GameT) (GameT, error) {

	//
	// check status
	//

	if err := game.Status().CanStop(); err != nil {
		return s.empty, err
	}

	//
	// preliminary checks
	//

	if err := s.plugin.CanStopGame(game); err != nil {
		return s.empty, err
	}

	//
	// stop game
	//

	game, err := s.plugin.StopGame(game)
	if err != nil {
		return s.empty, err
	}
	if !game.IsStopped() {
		game.SetStopped()
	}

	//
	// store game
	//

	game, err = s.storeGame(game)
	if err != nil {
		return s.empty, err
	}

	//
	// callbacks
	//

	s.onGame(game)

	return game, nil
}

// //////////////////////////////////////////////////
// delete game

func (s *gameService[PlayerT, GameT]) DeleteGameId(id model.GameId, playerId model.PlayerId) error {
	game, err := s.gameStore.Get(id)
	if err != nil {
		return err
	}
	return s.DeleteGame(game, playerId)
}

func (s *gameService[PlayerT, GameT]) DeleteGame(game GameT, playerId model.PlayerId) error {

	//
	// check status
	//

	if err := game.Status().CanDelete(); err != nil {
		return err
	}

	//
	// check player
	//

	if !game.HasPlayer(playerId) {
		return model.ErrPlayerNotInGame
	}

	//
	// preliminary checks
	//

	if err := s.plugin.CanDeleteGame(game, playerId); err != nil {
		return err
	}

	//
	// delete game
	//

	return s.deleteGame(game)
}

// //////////////////////////////////////////////////
// save game

func (s *gameService[PlayerT, GameT]) SaveGame(game GameT) (GameT, error) {
	if game.Status().IsMarkedForDeletion() {
		return s.empty, s.deleteGame(game)
	}
	return s.storeGame(game)
}

func (s *gameService[PlayerT, GameT]) storeGame(game GameT) (GameT, error) {
	if err := s.gameStore.Set(game); err != nil {
		return s.empty, err
	}
	return game, nil
}

func (s *gameService[PlayerT, GameT]) deleteGame(game GameT) error {
	return s.gameStore.Delete(game.Id())
}

// //////////////////////////////////////////////////
// callbacks

func (s *gameService[PlayerT, GameT]) RegisterOnJoinGame(onJoinFn func(game GameT, player PlayerT)) {
	s.onJoinFns = append(s.onJoinFns, onJoinFn)
}

func (s *gameService[PlayerT, GameT]) onJoinGame(game GameT, player PlayerT) {
	for _, onJoinFn := range s.onJoinFns {
		onJoinFn(game, player)
	}
}

func (s *gameService[PlayerT, GameT]) RegisterOnGame(onGameFn func(game GameT)) {
	s.onGameFns = append(s.onGameFns, onGameFn)
}

func (s *gameService[PlayerT, GameT]) onGame(game GameT) {
	for _, onGameFn := range s.onGameFns {
		onGameFn(game)
	}
}

func (s *gameService[PlayerT, GameT]) RegisterOnLeaveGame(onLeaveFn func(game GameT, player PlayerT)) {
	s.onLeaveFns = append(s.onLeaveFns, onLeaveFn)
}

func (s *gameService[PlayerT, GameT]) onLeaveGame(game GameT, player PlayerT) {
	for _, onLeaveFn := range s.onLeaveFns {
		onLeaveFn(game, player)
	}
}
