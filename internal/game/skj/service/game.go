package service

import (
	"go.uber.org/zap"

	share_model "github.com/gre-ory/games-go/internal/game/share/model"
	share_service "github.com/gre-ory/games-go/internal/game/share/service"
	share_websocket "github.com/gre-ory/games-go/internal/game/share/websocket"

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
	return game, nil
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

	// check if all cells are flipped
	if board.IsFlipped() {

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
		return share_model.ErrMissingPlayers
	}
	return nil
}

func (p *gamePlugin) StartGame(game *model.Game) (*model.Game, error) {

	//
	// draw cards & build player boards
	//

	for _, player := range game.Players() {
		board := model.NewPlayerBoard()
		for columnIndex := 0; columnIndex < game.NbColumn; columnIndex++ {
			column := model.NewPlayerColumn(columnIndex + 1)
			for rowIndex := 0; rowIndex < game.NbRow; rowIndex++ {
				card, err := game.DrawDeck.Draw()
				if err != nil {
					return nil, err
				}
				cell := model.NewPlayerCell(columnIndex+1, rowIndex+1, card)
				column.AddCell(cell)
			}
			board.AddColumn(column)
		}
		game.AddBoard(player.Id(), board)
	}

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
