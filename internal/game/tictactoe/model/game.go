package model

import (
	"strings"

	"github.com/gre-ory/games-go/internal/util"
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
		Id:      NewGameId(),
		Players: make(map[PlayerId]*Player),
		Rows:    rows,
		Round:   -1,
	}
	return game
}

type Game struct {
	Id        GameId
	Stopped   bool
	Players   map[PlayerId]*Player
	PlayerIds []PlayerId
	Round     int
	Rows      map[int]*Row
}

func (g *Game) Started() bool {
	return g.Round >= 0
}

func (g *Game) WithPlayer(player *Player) *Game {
	g.Players[player.GetId()] = player
	player.SetGameId(g.Id)
	return g
}

func (g *Game) WithoutPlayer(playerId PlayerId) *Game {
	delete(g.Players, playerId)
	return g
}

func (g *Game) CanStart() bool {
	return len(g.Players) == 2
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

	playerId := player.GetId()
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

func (g *Game) NextRound() {
	g.Round++
	g.SetPlayingPlayer()
}

func (g *Game) SetPlayingPlayer() {
	currentPlayerId := g.getCurrentPlayerId()
	for _, player := range g.Players {
		if currentPlayerId == player.GetId() {
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
			return player.GetId()
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
	labels = append(labels, "cell")
	switch c.Symbol {
	case NO_SYMBOL:
		labels = append(labels, "empty")
	case PLAYER_ONE_SYMBOL:
		labels = append(labels, "symbol-1")
	case PLAYER_TWO_SYMBOL:
		labels = append(labels, "symbol-2")
	}
	return strings.Join(labels, " ")
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
