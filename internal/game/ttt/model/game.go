package model

import (
	"html/template"
	"strings"

	share_model "github.com/gre-ory/games-go/internal/game/share/model"
	share_websocket "github.com/gre-ory/games-go/internal/game/share/websocket"
)

const (
	NbPlayer = 2
)

func NewGame(nbRow, nbColumn int) *Game {
	rows := make(map[int]*Row, nbRow)
	for y := 1; y <= nbRow; y++ {
		rows[y] = NewRow(nbColumn)
	}
	game := &Game{
		Game: share_model.NewGame[*Player](NbPlayer, NbPlayer),
		Rows: rows,
	}
	return game
}

type Game struct {
	share_model.Game[*Player]
	Rows map[int]*Row
}

func (g *Game) WrapData(data share_websocket.Data, player *Player) (bool, any) {
	data = data.With("game", g)

	playerId := player.Id()
	if playerId == "" {
		return true, data
	}
	player, found := g.GetPlayer(playerId)
	if !found {
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

func (g *Game) HasWinner() (bool, share_model.PlayerId) {
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

func (g *Game) GetPlayerIdFromRune(symbol rune) share_model.PlayerId {
	for _, player := range g.GetPlayers() {
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
