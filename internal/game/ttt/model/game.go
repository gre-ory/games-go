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

type GameId string

func NewGameId() GameId {
	return GameId(util.GenerateGameId())
}

func NewGame(nbRow, nbColumn int) *Game {
	rows := make(map[int]*Row, nbRow)
	for y := 1; y <= nbRow; y++ {
		rows[y] = NewRow(nbColumn)
	}
	game := &Game{
		id:        NewGameId(),
		CreatedAt: time.Now(),
		Players:   make(map[PlayerId]*Player),
		Rows:      rows,
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

type Game struct {
	id        GameId
	CreatedAt time.Time
	status    GameStatus
	WinnerIds []PlayerId
	Players   map[PlayerId]*Player
	PlayerIds []PlayerId
	Round     int
	Rows      map[int]*Row
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
		g.status = NotJoinable
		for _, player := range g.Players {
			player.Status = WaitingToStart
		}
	} else {
		g.status = Joinable
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
	return len(g.Players) < 2
}

func (g *Game) CanStart() bool {
	return len(g.Players) == 2
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
	if row, ok := g.Rows[y]; ok {
		return row.Play(player, x)
	}
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

func (g *Game) HasWinner() (bool, PlayerId) {
	for x := 1; x <= 3; x++ {
		same, symbol := g.HasSameSymbol(g.Rows[1].Cells[x], g.Rows[2].Cells[x], g.Rows[3].Cells[x])
		if same {
			return true, g.GetPlayerIdFromRune(symbol)
		}
	}
	for y := 1; y <= 3; y++ {
		same, symbol := g.HasSameSymbol(g.Rows[y].Cells[1], g.Rows[y].Cells[2], g.Rows[y].Cells[3])
		if same {
			return true, g.GetPlayerIdFromRune(symbol)
		}
	}
	same, symbol := g.HasSameSymbol(g.Rows[1].Cells[1], g.Rows[2].Cells[2], g.Rows[3].Cells[3])
	if same {
		return true, g.GetPlayerIdFromRune(symbol)
	}
	same, symbol = g.HasSameSymbol(g.Rows[1].Cells[3], g.Rows[2].Cells[2], g.Rows[3].Cells[1])
	if same {
		return true, g.GetPlayerIdFromRune(symbol)
	}
	return false, ""
}

func (g *Game) GetPlayerIdFromRune(symbol rune) PlayerId {
	for _, player := range g.Players {
		if player.Symbol == symbol {
			return player.Id()
		}
	}
	return ""
}

func (g *Game) HasSameSymbol(cells ...*Cell) (bool, rune) {
	symbol := NO_SYMBOL
	for index, cell := range cells {
		if cell.IsEmpty() {
			return false, NO_SYMBOL
		} else if index == 0 {
			symbol = cell.Symbol
		} else if cell.Symbol != symbol {
			return false, NO_SYMBOL
		}
	}
	return true, symbol
}

func (g *Game) IsTie() bool {
	for _, row := range g.Rows {
		for _, cell := range row.Cells {
			if cell.IsEmpty() {
				return false
			}
		}
	}
	if yes, _ := g.HasWinner(); yes {
		return false
	}
	return true
}

func (g *Game) Labels() string {
	labels := make([]string, 0)
	labels = append(labels, "game")
	switch g.status {
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
	if g.status == Stopped {
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
	if g.status == Stopped {
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
	if g.status == Stopped {
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
	if g.status == Stopped {
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
	if g.status == Stopped {
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

// //////////////////////////////////////////////////
// row

type Row struct {
	Cells map[int]*Cell
}

func NewRow(nbCell int) *Row {
	cells := make(map[int]*Cell, nbCell)
	for x := 1; x <= nbCell; x++ {
		cells[x] = NewCell()
	}
	return &Row{
		Cells: cells,
	}
}

func (r *Row) Play(player *Player, x int) error {
	if cell, ok := r.Cells[x]; ok {
		return cell.Play(player)
	}
	return ErrOutOfColumnBound
}

// //////////////////////////////////////////////////
// cell

type Cell struct {
	Symbol rune
}

var (
	NO_SYMBOL         rune = ' '
	PLAYER_ONE_SYMBOL rune = 'X'
	PLAYER_TWO_SYMBOL rune = 'O'
)

func NewCell() *Cell {
	return &Cell{
		Symbol: NO_SYMBOL,
	}
}

func (c *Cell) Play(player *Player) error {
	if c.Symbol != NO_SYMBOL {
		return ErrAlreadyPlayOnCell
	}
	c.Symbol = player.Symbol
	return nil
}

func (c *Cell) Labels() string {
	labels := make([]string, 0)
	switch c.Symbol {
	case PLAYER_ONE_SYMBOL:
		labels = append(labels, "symbol-1")
	case PLAYER_TWO_SYMBOL:
		labels = append(labels, "symbol-2")
	}
	return strings.Join(labels, " ")
}

func (c *Cell) IconHtml() template.HTML {
	switch c.Symbol {
	case PLAYER_ONE_SYMBOL:
		return "<div class=\"icon-cross\"></div>"
	case PLAYER_TWO_SYMBOL:
		return "<div class=\"icon-circle\"></div>"
	}
	return ""
}

func (c *Cell) String() string {
	return string(c.Symbol)
}

func (c *Cell) IsEmpty() bool {
	return c.Symbol == NO_SYMBOL
}

func (c *Cell) IsPlayer(player *Player) bool {
	return c.Symbol == player.Symbol
}
