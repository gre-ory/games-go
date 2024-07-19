package service

import (
	"fmt"
	"sort"

	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/game/share/model"
	"github.com/gre-ory/games-go/internal/game/share/store"
)

// //////////////////////////////////////////////////
// game service

type GameService[PlayerT model.Player, GameT model.Game[PlayerT]] interface {
	GetPlayer(playerId model.PlayerId) (PlayerT, error)

	GetGame(gameId model.GameId) (GameT, error)
	GetJoinableGames() []GameT
	GetNonJoinableGames(userId model.UserId) []GameT
	SortGamesByCreationTime(games []GameT) []GameT
	FilterGamesByPlayer(games []GameT, playerId model.PlayerId) []GameT

	CreateGame(user model.User) (GameT, error)
	JoinGameId(gameId model.GameId, user model.User) (GameT, error)
	JoinGame(game GameT, user model.User) (GameT, error)
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
	RegisterOnLeaveGame(func(game GameT, userId model.UserId))
}

// //////////////////////////////////////////////////
// game plugin

type GamePlugin[PlayerT model.Player, GameT model.Game[PlayerT]] interface {
	CanCreateGame(user model.User) error
	CreateGame(user model.User) (GameT, PlayerT, error)

	CanJoinGame(game GameT, user model.User) error
	JoinGame(game GameT, user model.User) (GameT, PlayerT, error)

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
	onLeaveFns []func(game GameT, userId model.UserId)
	empty      GameT
}

// //////////////////////////////////////////////////
// get player

func (s *gameService[PlayerT, GameT]) GetPlayer(playerId model.PlayerId) (PlayerT, error) {
	game, err := s.GetGame(playerId.GameId())
	if err != nil {
		var empty PlayerT
		return empty, err
	}
	player, found := game.Player(playerId)
	if !found {
		var empty PlayerT
		return empty, model.ErrPlayerNotFound
	}
	return player, nil
}

// //////////////////////////////////////////////////
// get game

func (s *gameService[PlayerT, GameT]) GetGame(id model.GameId) (GameT, error) {
	return s.gameStore.Get(id)
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

func (s *gameService[PlayerT, GameT]) GetNonJoinableGames(userId model.UserId) []GameT {
	games := make([]GameT, 0)
	games = append(games, s.gameStore.ListStatus(model.GameStatus_NotJoinableAndStartable)...)
	games = append(games, s.gameStore.ListStatus(model.GameStatus_Started)...)
	games = append(games, s.gameStore.ListStatus(model.GameStatus_Stopped)...)
	games = s.FilterGamesByUser(games, userId)
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

func (s *gameService[PlayerT, GameT]) FilterGamesByUser(games []GameT, userId model.UserId) []GameT {
	filtered := make([]GameT, 0, len(games))
	for _, game := range games {
		if game.HasUser(userId) {
			filtered = append(filtered, game)
		}
	}
	return filtered
}

// //////////////////////////////////////////////////
// create game

func (s *gameService[PlayerT, GameT]) CreateGame(user model.User) (GameT, error) {

	var game GameT
	var player PlayerT
	var err error

	s.logger.Info(fmt.Sprintf("[DEBUG] >>> create-game :: user %s", user.Id()))
	defer func() {
		s.logger.Info(fmt.Sprintf("[DEBUG] <<< create-game :: game %s %s :: player %s %s", game.Id(), game.Status().String(), player.Id(), player.Status().String()))
	}()

	//
	// preliminary checks
	//

	if err = s.plugin.CanCreateGame(user); err != nil {
		return s.empty, err
	}

	//
	// join game
	//

	game, player, err = s.plugin.CreateGame(user)
	if err != nil {
		return s.empty, err
	}
	if !game.HasPlayer(player.Id()) {
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
// join game

func (s *gameService[PlayerT, GameT]) JoinGameId(id model.GameId, user model.User) (GameT, error) {
	game, err := s.gameStore.Get(id)
	if err != nil {
		return s.empty, err
	}
	return s.JoinGame(game, user)
}

func (s *gameService[PlayerT, GameT]) JoinGame(game GameT, user model.User) (GameT, error) {

	var player PlayerT
	var err error

	s.logger.Info(fmt.Sprintf("[DEBUG] >>> join-game :: game %s %s :: user %s", game.Id(), game.Status().String(), user.Id()))
	defer func() {
		s.logger.Info(fmt.Sprintf("[DEBUG] <<< join-game :: game %s %s :: player %s %s", game.Id(), game.Status().String(), player.Id(), player.Status().String()))
	}()

	//
	// check status
	//

	if err := game.Status().CanJoin(); err != nil {
		return s.empty, err
	}

	//
	// check player
	//

	if ok := game.HasUser(user.Id()); ok {
		return game, nil
	}

	//
	// preliminary checks
	//

	if err := s.plugin.CanJoinGame(game, user); err != nil {
		return s.empty, err
	}

	//
	// join game
	//

	game, player, err = s.plugin.JoinGame(game, user)
	if err != nil {
		return s.empty, err
	}
	if !game.HasPlayer(player.Id()) {
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
	game, err := s.gameStore.Get(player.GameId())
	if err != nil {
		return s.empty, err
	}
	return s.StartGame(game)
}

func (s *gameService[PlayerT, GameT]) StartGame(game GameT) (GameT, error) {

	s.logger.Info(fmt.Sprintf("[DEBUG] >>> start-game :: game %s %s", game.Id(), game.Status().String()))
	defer func() {
		s.logger.Info(fmt.Sprintf("[DEBUG] <<< start-game :: game %s %s", game.Id(), game.Status().String()))
	}()

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

	s.logger.Info(fmt.Sprintf("[DEBUG] >>> leave-game :: game %s %s :: player %s %s", game.Id(), game.Status().String(), player.Id(), player.Status().String()))
	defer func() {
		s.logger.Info(fmt.Sprintf("[DEBUG] <<< leave-game :: game %s %s :: player %s %s", game.Id(), game.Status().String(), player.Id(), player.Status().String()))
	}()

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
	game.UpdateJoinStatus()

	//
	// save game
	//

	if game.IsMarkedForDeletion() {
		err = s.deleteGame(game)
		if err != nil {
			return s.empty, err
		}
	} else {
		game, err = s.SaveGame(game)
		if err != nil {
			return s.empty, err
		}
	}

	//
	// callbacks
	//

	userId := player.Id().UserId()
	s.onLeaveGame(game, userId)

	return game, nil
}

// //////////////////////////////////////////////////
// stop game

func (s *gameService[PlayerT, GameT]) StopGame(game GameT) (GameT, error) {

	s.logger.Info(fmt.Sprintf("[DEBUG] >>> stop-game :: game %s %s", game.Id(), game.Status().String()))
	defer func() {
		s.logger.Info(fmt.Sprintf("[DEBUG] <<< stop-game :: game %s %s", game.Id(), game.Status().String()))
	}()

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
	for _, player := range game.Players() {
		if !player.Status().HasPlayed() {
			player.SetStatus(model.PlayerStatus_Played)
		}
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

	s.logger.Info(fmt.Sprintf("[DEBUG] >>> delete-game :: game %s %s :: player %s", game.Id(), game.Status().String(), playerId))
	defer func() {
		s.logger.Info(fmt.Sprintf("[DEBUG] <<< delete-game :: game %s %s :: player %s", game.Id(), game.Status().String(), playerId))
	}()

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

func (s *gameService[PlayerT, GameT]) RegisterOnLeaveGame(onLeaveFn func(game GameT, userId model.UserId)) {
	s.onLeaveFns = append(s.onLeaveFns, onLeaveFn)
}

func (s *gameService[PlayerT, GameT]) onLeaveGame(game GameT, userId model.UserId) {
	for _, onLeaveFn := range s.onLeaveFns {
		onLeaveFn(game, userId)
	}
}
