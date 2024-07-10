package service

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"

	share_model "github.com/gre-ory/games-go/internal/game/share/model"
	share_service "github.com/gre-ory/games-go/internal/game/share/service"
	share_websocket "github.com/gre-ory/games-go/internal/game/share/websocket"

	"github.com/gre-ory/games-go/internal/game/czm/model"
	"github.com/gre-ory/games-go/internal/game/czm/store"
)

type GameService interface {
	share_service.GameService[*model.Player, *model.Game]

	SelectCard(player *model.Player, cardNumber int) (*model.Game, error)
	PlayCard(player *model.Player, discardNumber int) (*model.Game, error)

	WrapData(data share_websocket.Data, player *model.Player) (bool, any)
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

func (s *gameService) SelectCard(player *model.Player, cardNumber int) (*model.Game, error) {
	game, err := s.getPlayGame(player)
	if err != nil {
		return nil, err
	}

	_, err = player.SelectCard(cardNumber)
	if err != nil {
		return nil, err
	}

	game.SelectedCardNumber = cardNumber

	return game, nil
}

func (s *gameService) PlayCard(player *model.Player, discardNumber int) (*model.Game, error) {
	game, err := s.getPlayGame(player)
	if err != nil {
		return nil, err
	}
	if game.SelectedCardNumber == 0 {
		return nil, model.ErrNoSelectedCard
	}

	selectedCard, err := player.PlayCard(game.SelectedCardNumber)
	if err != nil {
		return nil, err
	}

	if discardNumber < 1 || model.NbCardDeck < discardNumber {
		return nil, model.ErrInvalidDiscardNumber
	}

	game.SelectedCardNumber = 0
	game.DiscardCardDecks[discardNumber].Add(selectedCard)

	newCard, err := game.DrawCardDeck.Draw()
	if err != nil {
		return nil, err
	}
	player.WithCard(newCard)

	// TODO
	// if yes, winnerId := game.HasWinner(); yes {
	// 	s.stopGame(game, winnerId)
	// } else if game.IsTie() {
	// 	s.stopGame(game, "")
	// } else {
	// 	game.Round++
	// 	game.SetPlayingPlayer()
	// }

	return game, nil
}

func (s *gameService) ValidateMissions(player *model.Player, cardIndex int, discardIndex int) (*model.Game, error) {
	game, err := s.getPlayGame(player)
	if err != nil {
		return nil, err
	}

	topCards := game.GetTopCards()

	game.ValidatedMissionIndex = -1
	for index, mission := range game.Missions {
		if mission == nil {
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

func (s *gameService) getPlayGame(player *model.Player) (*model.Game, error) {
	game, err := s.GetGame(player.GameId())
	if err != nil {
		return nil, err
	}
	if err := game.Status().CanPlay(); err != nil {
		return nil, err
	}
	if !player.Status().IsPlaying() {
		return nil, share_model.ErrWrongPlayer
	}
	return game, nil
}

func (s *gameService) WrapData(data share_websocket.Data, player *model.Player) (bool, any) {
	if player == nil {
		return true, data
	}
	data.With("lang", model.App.PlayerLocalizer(player))
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
	return model.NewGame(), nil
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
		return share_model.ErrMissingPlayers
	}
	return nil
}

func (p *gamePlugin) StartGame(game *model.Game) (*model.Game, error) {

	//
	// set random order
	//

	game.SetRandomOrder()

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
		// set other player as winner
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
