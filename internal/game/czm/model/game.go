package model

import (
	"html/template"
	"strings"
	"time"

	"github.com/gre-ory/games-go/internal/util"
	"github.com/gre-ory/games-go/internal/util/list"
	"github.com/gre-ory/games-go/internal/util/loc"
	"github.com/gre-ory/games-go/internal/util/websocket"
)

const (
	Game_MinPlayer = 2
	Game_MaxPlayer = 4
)

type GameId string

func NewGameId() GameId {
	return GameId(util.GenerateGameId())
}

func NewGame() *Game {
	game := &Game{
		id:           NewGameId(),
		CreatedAt:    time.Now(),
		Players:      make(map[PlayerId]*Player),
		DrawCardDeck: NewDrawCardDeck(),
		DiscardCardDecks: [NbCardDeck]CardDeck{
			NewDiscardCardDeck(),
			NewDiscardCardDeck(),
			NewDiscardCardDeck(),
			NewDiscardCardDeck(),
		},
		SelectedCardIndex:  -1,
		DrawMissionDeck:    NewDrawMissionDeck(),
		Missions:           [NbMission]Mission{},
		DiscardMissionDeck: NewDiscardMissionDeck(),
	}
	return game
}

type GameStatus int

const (
	Joinable GameStatus = iota
	NotJoinable
	Started
	Stopped
)

const (
	NbCardDeck = 4
	NbMission  = 4
)

type TopMissions [NbMission]Mission

type Game struct {
	id                    GameId
	CreatedAt             time.Time
	status                GameStatus
	WinnerIds             []PlayerId
	Players               map[PlayerId]*Player
	PlayerIds             []PlayerId
	Round                 int
	DrawCardDeck          CardDeck
	DiscardCardDecks      [NbCardDeck]CardDeck
	SelectedCardIndex     int
	DrawMissionDeck       MissionDeck
	Missions              [NbMission]Mission
	DiscardMissionDeck    MissionDeck
	ValidatedMissionIndex int
}

func (g *Game) Id() GameId {
	return g.id
}

func (g *Game) Status() GameStatus {
	return g.status
}

func (g *Game) SetStatus(status GameStatus) {
	g.status = status
}

func (g *Game) GetCreatedAt() time.Time {
	return g.CreatedAt
}

func (g *Game) Started() bool {
	return g.status == Started
}

func (g *Game) Stopped() bool {
	return g.status == Stopped
}

func (g *Game) WithPlayer(player *Player) *Game {
	g.Players[player.Id()] = player
	player.SetGameId(g.id)
	return g
}

func (g *Game) WithoutPlayer(player *Player) *Game {
	delete(g.Players, player.Id())
	player.UnsetGameId()
	return g
}

func (g *Game) UpdateStatus() {
	if !g.CanJoin() {
		g.SetStatus(NotJoinable)
		for _, player := range g.Players {
			player.Status = WaitingToStart
		}
	} else {
		g.SetStatus(Joinable)
		if g.CanStart() {
			for _, player := range g.Players {
				player.Status = WaitingToJoinOrStart
			}
		} else {
			for _, player := range g.Players {
				player.Status = WaitingToJoin
			}
		}
	}
}

func (g *Game) CanJoin() bool {
	return len(g.Players) < Game_MaxPlayer
}

func (g *Game) CanStart() bool {
	return len(g.Players) >= Game_MinPlayer
}

func (g *Game) HasPlayer(playerId PlayerId) bool {
	return list.Contains(g.PlayerIds, playerId)
}

func (g *Game) GetPlayer(id PlayerId) (*Player, error) {
	if player, ok := g.Players[id]; ok {
		return player, nil
	}
	return nil, ErrPlayerNotFound
}

func (g *Game) GetCurrentPlayerId() (PlayerId, error) {
	if !g.Started() {
		return "", ErrGameNotStarted
	}
	return g.getCurrentPlayerId(), nil
}

func (g *Game) IsCardSelectable(playerId PlayerId, cardIndex int) bool {
	currentPlayerId, err := g.GetCurrentPlayerId()
	if err != nil {
		return false
	}
	if playerId != currentPlayerId {
		return false
	}
	return cardIndex != g.SelectedCardIndex
}

func (g *Game) IsCardSelected(playerId PlayerId, cardIndex int) bool {
	currentPlayerId, err := g.GetCurrentPlayerId()
	if err != nil {
		return false
	}
	if playerId != currentPlayerId {
		return false
	}
	return cardIndex == g.SelectedCardIndex
}

func (g *Game) GetTopCards() TopCards {
	topCards := TopCards{}
	for index, discardDeck := range g.DiscardCardDecks {
		topCards[index] = discardDeck.GetTopCard()
	}
	return topCards
}

func (g *Game) GetOtherPlayerId(playerId PlayerId) (PlayerId, error) {
	for id := range g.Players {
		if id != playerId {
			return id, nil
		}
	}
	return "", ErrPlayerNotFound
}

func (g *Game) getCurrentPlayerId() PlayerId {
	return g.PlayerIds[g.Round%len(g.PlayerIds)]
}

func (g *Game) GetCurrentPlayer() (*Player, error) {
	currentPlayerId, err := g.GetCurrentPlayerId()
	if err != nil {
		return nil, err
	}
	if player, ok := g.Players[currentPlayerId]; ok {
		return player, err
	}
	return nil, ErrPlayerNotFound
}

func (g *Game) WrapData(data websocket.Data, player *Player) (bool, any) {
	data = data.With("game", g)

	playerId := player.Id()
	if playerId == "" {
		return true, data
	}
	player, err := g.GetPlayer(playerId)
	if err != nil {
		return false, nil
	}
	return true, data.With("player", player)
}

func (g *Game) Play(player *Player, x, y int) error {
	// if row, ok := g.Rows[y]; ok {
	// 	return row.Play(player, x)
	// }
	return ErrOutOfRowBound
}

func (g *Game) SetPlayingPlayer() {
	currentPlayerId := g.getCurrentPlayerId()
	for _, player := range g.Players {
		if currentPlayerId == player.Id() {
			player.Status = Playing
		} else {
			player.Status = WaitingToPlay
		}
	}
}

func (g *Game) HasValidatedMission() bool {
	return g.ValidatedMissionIndex >= 0
}

func (g *Game) GetValidatedMissionIndex() int {
	return g.ValidatedMissionIndex
}

func (g *Game) HasBronzeMedal() bool {
	return len(g.DiscardMissionDeck) > 15
}

func (g *Game) HasSilverMedal() bool {
	return len(g.DiscardMissionDeck) > 19
}

func (g *Game) HasGoldMedal() bool {
	return len(g.DiscardMissionDeck) > 23
}

func (g *Game) HasWinner() (bool, PlayerId) {
	// for x := 1; x <= 3; x++ {
	// 	same, symbol := g.HasSameSymbol(g.Rows[1].Cells[x], g.Rows[2].Cells[x], g.Rows[3].Cells[x])
	// 	if same {
	// 		return true, g.GetPlayerIdFromRune(symbol)
	// 	}
	// }
	// for y := 1; y <= 3; y++ {
	// 	same, symbol := g.HasSameSymbol(g.Rows[y].Cells[1], g.Rows[y].Cells[2], g.Rows[y].Cells[3])
	// 	if same {
	// 		return true, g.GetPlayerIdFromRune(symbol)
	// 	}
	// }
	// same, symbol := g.HasSameSymbol(g.Rows[1].Cells[1], g.Rows[2].Cells[2], g.Rows[3].Cells[3])
	// if same {
	// 	return true, g.GetPlayerIdFromRune(symbol)
	// }
	// same, symbol = g.HasSameSymbol(g.Rows[1].Cells[3], g.Rows[2].Cells[2], g.Rows[3].Cells[1])
	// if same {
	// 	return true, g.GetPlayerIdFromRune(symbol)
	// }
	return false, ""
}

func (g *Game) GetPlayerIdFromRune(symbol rune) PlayerId {
	// for _, player := range g.Players {
	// 	if player.Symbol == symbol {
	// 		return player.Id()
	// 	}
	// }
	return ""
}

func (g *Game) IsTie() bool {
	// for _, row := range g.Rows {
	// 	for _, cell := range row.Cells {
	// 		if cell.IsEmpty() {
	// 			return false
	// 		}
	// 	}
	// }
	// if yes, _ := g.HasWinner(); yes {
	// 	return false
	// }
	return true
}

func (g *Game) Labels() string {
	labels := make([]string, 0)
	labels = append(labels, "game")
	switch g.Status() {
	case Joinable:
		labels = append(labels, "joinable")
	case NotJoinable:
		labels = append(labels, "not-joinable")
	case Started:
		labels = append(labels, "started")
	case Stopped:
		labels = append(labels, "stopped")
	}
	return strings.Join(labels, " ")
}

type PlayerResult int

const (
	PlayerResult_Undefined PlayerResult = iota
	PlayerResult_Win
	PlayerResult_Tie
	PlayerResult_Loose
)

func (g *Game) PlayerResult(playerId PlayerId) PlayerResult {
	if g.Status() == Stopped {
		if len(g.WinnerIds) == 0 {
			return PlayerResult_Tie
		} else if list.Contains(g.WinnerIds, playerId) {
			return PlayerResult_Win
		} else {
			return PlayerResult_Loose
		}
	}
	return PlayerResult_Undefined
}

func (g *Game) PlayerLabels(playerId PlayerId) string {
	if g.Status() == Stopped {
		labels := make([]string, 0)
		labels = append(labels, "player")
		switch g.PlayerResult(playerId) {
		case PlayerResult_Win:
			labels = append(labels, "win")
		case PlayerResult_Tie:
			labels = append(labels, "tie")
		case PlayerResult_Loose:
			labels = append(labels, "loose")
		}
		return strings.Join(labels, " ")
	}
	player, err := g.GetPlayer(playerId)
	if err != nil {
		return "error"
	}
	return player.Labels()
}

func (g *Game) YourPlayerMessage(localizer loc.Localizer, playerId PlayerId) template.HTML {
	if g.Status() == Stopped {
		switch g.PlayerResult(playerId) {
		case PlayerResult_Win:
			return localizer.Loc("YouWin")
		case PlayerResult_Tie:
			return localizer.Loc("YouTie")
		case PlayerResult_Loose:
			return localizer.Loc("YouLoose")
		}
		return ""
	}
	player, err := g.GetPlayer(playerId)
	if err != nil {
		return localizer.Loc("Error", err.Error())
	}
	return player.YourMessage(localizer)
}

func (g *Game) PlayerMessage(localizer loc.Localizer, playerId PlayerId) template.HTML {
	if g.Status() == Stopped {
		switch g.PlayerResult(playerId) {
		case PlayerResult_Win:
			return localizer.Loc("PlayerWin")
		case PlayerResult_Tie:
			return localizer.Loc("PlayerTie")
		case PlayerResult_Loose:
			return localizer.Loc("PlayerLoose")
		}
		return ""
	}
	player, err := g.GetPlayer(playerId)
	if err != nil {
		return localizer.Loc("Error", err.Error())
	}
	return player.Message(localizer)
}

func (g *Game) PlayerStatusIcon(playerId PlayerId) string {
	if g.Status() == Stopped {
		switch g.PlayerResult(playerId) {
		case PlayerResult_Win:
			return "icon-win"
		case PlayerResult_Tie:
			return "icon-tie"
		case PlayerResult_Loose:
			return "icon-loose"
		}
		return ""
	}
	player, err := g.GetPlayer(playerId)
	if err != nil {
		return ""
	}
	return player.StatusIcon()
}
