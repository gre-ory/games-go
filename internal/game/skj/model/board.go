package model

import "strings"

// //////////////////////////////////////////////////
// board

const (
	NbRow    = 3
	NbColumn = 4
)

type PlayerBoard struct {
	columns []*PlayerColumn
}

func NewPlayerBoard() *PlayerBoard {
	board := &PlayerBoard{
		columns: make([]*PlayerColumn, 0),
	}
	return board
}

func (board *PlayerBoard) AddColumn(column *PlayerColumn) {
	board.columns = append(board.columns, column)
}

func (board *PlayerBoard) Columns() []*PlayerColumn {
	return board.columns
}

func (board *PlayerBoard) IsFlipped() bool {
	for _, column := range board.columns {
		if !column.IsFlipped() {
			return false
		}
	}
	return true
}

func (board *PlayerBoard) Total() int {
	result := 0
	for _, column := range board.columns {
		result += column.Total()
	}
	return result
}

func (board *PlayerBoard) Flip(columnIndex, rowIndex int) error {
	if columnIndex < 0 || columnIndex >= NbColumn {
		return ErrInvalidColumn
	}
	return board.columns[columnIndex].Flip(rowIndex)
}

func (board *PlayerBoard) Labels() string {
	labels := make([]string, 0)
	labels = append(labels, "board")
	if board.IsFlipped() {
		labels = append(labels, "flipped")
	}
	return strings.Join(labels, " ")
}

// //////////////////////////////////////////////////
// player column

type PlayerColumn struct {
	columnNumber int
	cells        []*PlayerCell
}

func NewPlayerColumn(columnNumber int) *PlayerColumn {
	column := &PlayerColumn{
		columnNumber: columnNumber,
		cells:        make([]*PlayerCell, 0),
	}
	return column
}

func (column *PlayerColumn) AddCell(cell *PlayerCell) {
	column.cells = append(column.cells, cell)
}

func (column *PlayerColumn) Cells() []*PlayerCell {
	return column.cells
}

func (column *PlayerColumn) IsSkyjo() bool {
	firstCell := column.cells[0]
	for _, cell := range column.cells {
		if !cell.IsFlipped() {
			return false
		}
		if cell.card != firstCell.card {
			return false
		}
	}
	return true
}

func (column *PlayerColumn) IsFlipped() bool {
	for _, cell := range column.cells {
		if !cell.IsFlipped() {
			return false
		}
	}
	return true
}

func (column *PlayerColumn) Flip(rowIndex int) error {
	if rowIndex < 0 || rowIndex >= NbRow {
		return ErrInvalidRow
	}
	return column.cells[rowIndex].Flip()
}

func (column *PlayerColumn) Total() int {
	if column.IsSkyjo() {
		return 0
	}
	result := 0
	for _, cell := range column.cells {
		result += cell.Total()
	}
	return result
}

func (column *PlayerColumn) Labels() string {
	labels := make([]string, 0)
	labels = append(labels, "column")
	if column.IsFlipped() {
		labels = append(labels, "flipped")
	}
	if column.IsSkyjo() {
		labels = append(labels, "skyjo")
	}
	return strings.Join(labels, " ")
}

// //////////////////////////////////////////////////
// player cell

type PlayerCell struct {
	columnNumber int
	rowNumber    int
	card         Card
	flipped      bool
}

func NewPlayerCell(columnNumber, rowNumber int, card Card) *PlayerCell {
	return &PlayerCell{
		columnNumber: columnNumber,
		rowNumber:    rowNumber,
		card:         card,
		flipped:      false,
	}
}

func (cell *PlayerCell) Column() int {
	return cell.columnNumber
}

func (cell *PlayerCell) Row() int {
	return cell.rowNumber
}

func (cell *PlayerCell) Card() int {
	return int(cell.card)
}

func (cell *PlayerCell) IsVisible() bool {
	return !cell.flipped
}

func (cell *PlayerCell) IsFlipped() bool {
	return cell.flipped
}

func (cell *PlayerCell) CanFlip() bool {
	return !cell.flipped
}

func (cell *PlayerCell) Flip() error {
	if cell.flipped {
		return ErrCardAlreadyFlipped
	}
	cell.flipped = true
	return nil
}

func (cell *PlayerCell) Total() int {
	if cell.IsFlipped() {
		return 0
	}
	return int(cell.card)
}

func (cell *PlayerCell) Labels() string {
	labels := make([]string, 0)
	labels = append(labels, "cell")
	if cell.flipped {
		labels = append(labels, "flipped")
	}
	return strings.Join(labels, " ")
}
