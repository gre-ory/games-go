package service

import (
	"fmt"

	"go.uber.org/zap"

	share_model "github.com/gre-ory/games-go/internal/game/share/model"
	share_service "github.com/gre-ory/games-go/internal/game/share/service"
	"github.com/gre-ory/games-go/internal/util/loc"

	"github.com/gre-ory/games-go/internal/game/skj/model"
	"github.com/gre-ory/games-go/internal/game/skj/store"
)

type GameService interface {
	share_service.GameService[*model.Player, *model.Game]

	DrawDiscardCard(player *model.Player) (*model.Game, error)
	DrawCard(player *model.Player) (*model.Game, error)
	PutCard(player *model.Player, columnNumber, rowNumber int) (*model.Game, error)
	DiscardCard(player *model.Player) (*model.Game, error)
	FlipCard(player *model.Player, columnNumber, rowNumber int) (*model.Game, error)
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

func (s *gameService) DrawDiscardCard(player *model.Player) (*model.Game, error) {
	game, err := s.getPlayGame(player)
	if err != nil {
		return nil, err
	}
	if game.SelectedCard != nil {
		return nil, model.ErrAlreadySelectedCard
	}
	card, err := game.DiscardDeck.Draw()
	if err != nil {
		return nil, err
	}
	game.SelectedCard = &card

	player.SetYourMessage(loc.NewMessage("YouPutting").With("card", card.Value()))
	player.SetMessage(loc.NewMessage("PlayerHasDrawn").With("card", card.Value()))

	return game, nil
}

func (s *gameService) DrawCard(player *model.Player) (*model.Game, error) {
	game, err := s.getPlayGame(player)
	if err != nil {
		return nil, err
	}
	if game.SelectedCard != nil {
		return nil, model.ErrAlreadySelectedCard
	}
	card, err := game.DrawDeck.Draw()
	if err != nil {
		return nil, err
	}

	game.SelectedCard = &card

	player.SetYourMessage(loc.NewMessage("YouPuttingOrDiscarding").With("card", card.Value()))
	player.SetMessage(loc.NewMessage("PlayerHasDrawn").With("card", card.Value()))

	return game, nil
}

func (s *gameService) PutCard(player *model.Player, columnNumber, rowNumber int) (*model.Game, error) {
	game, err := s.getPlayGame(player)
	if err != nil {
		return nil, err
	}
	if game.SelectedCard == nil {
		return nil, model.ErrMissingSelectedCard
	}
	board, err := s.getBoard(game, player)
	if err != nil {
		return nil, err
	}
	cardToDiscard, err := board.Put(*game.SelectedCard, columnNumber-1, rowNumber-1)
	if err != nil {
		return nil, err
	}

	game.DiscardDeck.Add(cardToDiscard)
	game.SelectedCard = nil

	player.ResetMessages()

	return s.endPlayerRound(player, game, board)
}

func (s *gameService) DiscardCard(player *model.Player) (*model.Game, error) {
	game, err := s.getPlayGame(player)
	if err != nil {
		return nil, err
	}
	if game.SelectedCard == nil {
		return nil, model.ErrMissingSelectedCard
	}
	game.DiscardDeck.Add(*game.SelectedCard)

	player.SetYourMessage(loc.NewMessage("YouFlipping"))
	player.SetMessage(loc.NewMessage("PlayerHasDiscarded").With("card", game.SelectedCard.Value()))

	game.SelectedCard = nil
	game.ShouldFlip = true

	return game, nil
}

func (s *gameService) FlipCard(player *model.Player, columnNumber, rowNumber int) (*model.Game, error) {
	game, err := s.getPlayGame(player)
	if err != nil {
		return nil, err
	}
	if !game.ShouldFlip {
		return nil, model.ErrNotShouldFlip
	}
	board, err := s.getBoard(game, player)
	if err != nil {
		return nil, err
	}
	err = board.Flip(columnNumber-1, rowNumber-1)
	if err != nil {
		return nil, err
	}
	game.ShouldFlip = false

	return s.endPlayerRound(player, game, board)
}

func (s *gameService) endPlayerRound(player *model.Player, game *model.Game, board *model.PlayerBoard) (*model.Game, error) {

	if !game.LastTurn {
		// check if all cells are flipped
		if board.IsFlipped() {
			s.logger.Info("[DEBUG] last turn")
			game.LastTurn = true
		}
	}
	if game.LastTurn {
		s.logger.Info(fmt.Sprintf("[DEBUG] player %s: END", player.Id()))
		player.SetStatus(share_model.PlayerStatus_Played)
		s.logger.Info(fmt.Sprintf("[DEBUG] player %s: %s", player.Id(), player.Status().String()))
	}
	s.logger.Info(fmt.Sprintf("[DEBUG] game ended? %t", game.IsEnded()))
	for _, player := range game.Players() {
		s.logger.Info(fmt.Sprintf("[DEBUG] > player %s: %s", player.Id(), player.Status().String()))
	}
	if game.IsEnded() {
		s.logger.Info("[DEBUG] flip all")
		game.FlipAll()
		winnerId := game.WinnerId()
		if winnerId != "" {
			s.logger.Info(fmt.Sprintf("[DEBUG] winner is %s", winnerId))
			game.SetWinners(winnerId)
		}
		s.logger.Info("[DEBUG] game stopped")
		game.SetStopped()
	} else {
		s.logger.Info("[DEBUG] next player round")
		game.NextRound()
		game.SetPlayingRoundPlayer()
		nextPlayer := game.RoundPlayer()
		nextPlayer.SetYourMessage(loc.NewMessage("YouDrawing"))
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

func (s *gameService) getBoard(game *model.Game, player *model.Player) (*model.PlayerBoard, error) {
	if board, found := game.GetBoard(player.Id()); found {
		return board, nil
	}
	return nil, model.ErrPlayerBoardNotFound
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
	game := model.NewGame()
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
	// draw cards & build player boards
	//

	for _, player := range game.Players() {
		board := player.Board()
		for columnIndex := 0; columnIndex < game.NbColumn; columnIndex++ {
			column := board.NewColumn()
			for rowIndex := 0; rowIndex < game.NbRow; rowIndex++ {
				card, err := game.DrawDeck.Draw()
				if err != nil {
					return nil, err
				}
				column.NewCell(card)
			}
		}
	}

	//
	// discard first card
	//

	firstCard, err := game.DrawDeck.Draw()
	if err != nil {
		return nil, err

	}
	game.DiscardDeck.Add(firstCard)

	//
	// set random order
	//

	game.SetRandomOrder()

	//
	// set first playing player
	//

	game.FirstRound()
	game.SetPlayingRoundPlayer()

	firstlayer := game.RoundPlayer()
	firstlayer.SetYourMessage(loc.NewMessage("YouDrawing"))

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
