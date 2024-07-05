package service

import (
	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/util/loc"

	share_model "github.com/gre-ory/games-go/internal/game/share/model"
	share_service "github.com/gre-ory/games-go/internal/game/share/service"
	share_websocket "github.com/gre-ory/games-go/internal/game/share/websocket"

	"github.com/gre-ory/games-go/internal/game/ttt/model"
	"github.com/gre-ory/games-go/internal/game/ttt/store"
)

type GameService interface {
	share_service.GameService[*model.Player, *model.Game]
	PlayPlayerGame(player *model.Player, x, y int) (*model.Game, error)
	WrapData(data share_websocket.Data, player *model.Player) (bool, any)
}

func NewGameService(logger *zap.Logger, gameStore store.GameStore, playerStore store.PlayerStore) GameService {
	plugin := NewGamePlugin()
	return &gameService{
		GameService: share_service.NewGameService(logger, plugin, gameStore, playerStore),
		logger:      logger,
	}
}

type gameService struct {
	share_service.GameService[*model.Player, *model.Game]
	logger *zap.Logger
}

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
		return nil, model.ErrWrongPlayer
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

func (s *gameService) WrapData(data share_websocket.Data, player *model.Player) (bool, any) {
	if player == nil {
		return true, data
	}
	localizer := loc.NewLocalizer(model.AppId, loc.Language(player.Language()), s.logger)
	// s.logger.Info(fmt.Sprintf("[wrap] player %v: lang=%s ( %s )", player.Id(), player.Language, localizer.Loc("GameTitle", "ABC")))
	data.With("lang", localizer)
	if !player.HasGameId() {
		return true, data
	}
	game, err := s.GetGame(player.GameId())
	if err != nil {
		return false, nil
	}
	return game.WrapData(data, player)
}

// //////////////////////////////////////////////
// game plugin

func NewGamePlugin() share_service.GamePlugin[*model.Player, *model.Game] {
	return &gamePlugin{}
}

type gamePlugin struct{}

func (p *gamePlugin) CanCreateGame(player *model.Player) error {
	return nil
}

func (p *gamePlugin) CreateGame(player *model.Player) (*model.Game, error) {
	return model.NewGame(3, 3), nil
}

func (p *gamePlugin) CanJoinGame(game *model.Game, player *model.Player) error {
	return nil
}

func (p *gamePlugin) JoinGame(game *model.Game, player *model.Player) (*model.Game, error) {
	game.AttachPlayer(player)
	return game, nil
}

func (p *gamePlugin) CanStartGame(game *model.Game) error {
	if !game.CanStart() {
		return model.ErrMissingPlayers
	}
	return nil
}

func (p *gamePlugin) StartGame(game *model.Game) (*model.Game, error) {
	game.SetRandomOrder()
	game.GetOrderedPlayer(0).SetSymbol(model.PLAYER_ONE_SYMBOL)
	game.GetOrderedPlayer(1).SetSymbol(model.PLAYER_TWO_SYMBOL)

	game.Start()
	game.SetPlayingRoundPlayer()

	return game, nil
}

func (p *gamePlugin) CanStopGame(game *model.Game) error {
	return nil
}

func (p *gamePlugin) CanLeaveGame(game *model.Game, player *model.Player) error {
	return nil
}

func (p *gamePlugin) LeaveGame(game *model.Game, player *model.Player) (*model.Game, error) {
	switch {
	case game.Status().IsStopped():
	case game.Status().IsStarted():
		// set other player as winner
		game.SetWinnerOthers(player.Id())
		game.Stop()
	default:
		game.DetachPlayer(player)
		if !game.HasPlayers() {
			game.MarkForDeletion()
		} else {
			game.UpdateStatus()
		}
	}
	return game, nil
}

func (p *gamePlugin) CanDeleteGame(game *model.Game, playerId share_model.PlayerId) error {
	return nil
}
