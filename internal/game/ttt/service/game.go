package service

import (
	"go.uber.org/zap"

	share_api "github.com/gre-ory/games-go/internal/game/share/api"
	share_model "github.com/gre-ory/games-go/internal/game/share/model"
	share_service "github.com/gre-ory/games-go/internal/game/share/service"

	"github.com/gre-ory/games-go/internal/game/ttt/model"
	"github.com/gre-ory/games-go/internal/game/ttt/store"
)

type GameService interface {
	share_service.GameService[*model.Player, *model.Game]
	PlayPlayerGame(player *model.Player, x, y int) (*model.Game, error)
}

func NewGameService(logger *zap.Logger, gameStore store.GameStore) GameService {
	plugin := NewGamePlugin()
	return &gameService{
		GameService: share_service.NewGameService(logger, plugin, gameStore),
		logger:      logger,
	}
}

type gameService struct {
	share_service.GameService[*model.Player, *model.Game]
	logger *zap.Logger
}

var _ share_api.GameService[*model.Player, *model.Game] = &gameService{}

func (s *gameService) PlayPlayerGame(player *model.Player, x, y int) (*model.Game, error) {
	game, err := s.GetGame(player.GameId())
	if err != nil {
		return nil, err
	}
	return s.PlayGame(game, player, x, y)
}

func (s *gameService) PlayGame(game *model.Game, player *model.Player, x, y int) (*model.Game, error) {
	if err := game.Status().CanPlay(); err != nil {
		return nil, err
	}
	if !player.Status().IsPlaying() {
		return nil, share_model.ErrWrongPlayer
	}

	err := game.Play(player, x, y)
	if err != nil {
		return nil, err
	}

	if yes, winnerId := game.HasWinner(); yes {
		game.SetWinners(winnerId)
		s.StopGame(game)
	} else if game.IsTie() {
		game.SetTie()
		s.StopGame(game)
	} else {
		game.NextRound()
		game.SetPlayingRoundPlayer()
	}

	return s.SaveGame(game)
}

// //////////////////////////////////////////////
// game plugin

func NewGamePlugin() share_service.GamePlugin[*model.Player, *model.Game] {
	return &gamePlugin{}
}

type gamePlugin struct{}

func (p *gamePlugin) CanCreateGame(user share_model.User) error {
	return nil
}

func (p *gamePlugin) CreateGame(user share_model.User) (*model.Game, *model.Player, error) {
	game := model.NewGame(3, 3)
	player := model.NewPlayerFromUser(game.Id(), user)
	return game, player, nil
}

func (p *gamePlugin) CanJoinGame(game *model.Game, user share_model.User) error {
	return nil
}

func (p *gamePlugin) JoinGame(game *model.Game, user share_model.User) (*model.Game, *model.Player, error) {
	player := model.NewPlayerFromUser(game.Id(), user)
	return game, player, nil
}

func (p *gamePlugin) CanStartGame(game *model.Game) error {
	if !game.CanStart() {
		return share_model.ErrMissingPlayers
	}
	return nil
}

func (p *gamePlugin) StartGame(game *model.Game) (*model.Game, error) {

	//
	// set random order
	//

	game.SetRandomOrder()
	game.OrderedPlayer(0).SetSymbol(model.PLAYER_ONE_SYMBOL)
	game.OrderedPlayer(1).SetSymbol(model.PLAYER_TWO_SYMBOL)

	//
	// set first playing player
	//

	game.FirstRound()
	game.SetPlayingRoundPlayer()

	return game, nil
}

func (p *gamePlugin) CanStopGame(game *model.Game) error {
	return nil
}

func (p *gamePlugin) StopGame(game *model.Game) (*model.Game, error) {
	return game, nil
}

func (p *gamePlugin) CanLeaveGame(game *model.Game, player *model.Player) error {
	return nil
}

func (p *gamePlugin) LeaveGame(game *model.Game, player *model.Player) (*model.Game, error) {
	switch {
	case game.IsStopped():
	case game.IsStarted():
		// set current player as looser
		game.SetLoosers(player.Id())
		game.SetStopped()
	default:
		game.DetachPlayer(player)
		if !game.HasPlayers() {
			game.MarkForDeletion()
		} else {
			game.UpdateJoinStatus()
		}
	}
	return game, nil
}

func (p *gamePlugin) CanDeleteGame(game *model.Game, playerId share_model.PlayerId) error {
	return nil
}
