package game

import (
	"fmt"

	. "github.com/kitsunemikan/ttt-cli/geom"
)

// We use the fact that the only deltas possible are:
// 1. Unavailable -> Unoccupied
// 2. Unoccupied -> P1 or P2
//
// Both of these can be reversed in one-to-one correspondance
type cellDelta struct {
	Cell     Offset
	NewState CellState
}

type boardDelta struct {
	Cells         []cellDelta
	OldBoardBound Rect
	NewBoardBound Rect
}

type BoardState struct {
	board           map[Offset]CellState
	unoccupiedCells map[Offset]struct{}
	playerCells     [2]map[Offset]struct{}

	// delta is assumed to be immutable as well as boardDelta values.
	// So, Clone() will share delta values
	delta       []boardDelta
	moveHistory []PlayerMove

	// Precalculated disk offsets, immutable
	circleMask []Offset

	borderWidth int
	boardBound  Rect
}

func generateCircleMask(radius int) (mask []Offset) {
	if radius == 0 {
		return
	}

	mask = make([]Offset, 0, radius*radius)

	for dx := -radius; dx <= radius; dx++ {
		for dy := -radius; dy <= radius; dy++ {
			ds := Offset{X: dx, Y: dy}
			if !ds.IsInsideCircle(radius) {
				continue
			}

			mask = append(mask, ds)
		}
	}

	return
}

func NewBoardState(borderWidth int) *BoardState {
	bs := &BoardState{
		board:           make(map[Offset]CellState),
		unoccupiedCells: make(map[Offset]struct{}),

		circleMask: generateCircleMask(borderWidth),

		borderWidth: borderWidth,
		boardBound:  Rect{X: -borderWidth, Y: -borderWidth, W: 2*borderWidth + 1, H: 2*borderWidth + 1},
	}

	bs.playerCells[0] = make(map[Offset]struct{})
	bs.playerCells[1] = make(map[Offset]struct{})

	// Mark initial available cells
	for _, ds := range bs.circleMask {
		bs.MarkUnoccupied(ds)
	}

	return bs
}

// NewBoardStateFromCells expects a non-zero border width
func NewBoardStateFromCells(borderWidth int, cells map[Offset]CellState) *BoardState {
	bs := &BoardState{
		board: make(map[Offset]CellState, len(cells)),
		// Size's just a hint, I will trade performance for extra memory consumption
		// Assuming for one player move there are ~borderWidth*borderWidth new cells
		// It's basically almost the full len(cells)
		unoccupiedCells: make(map[Offset]struct{}, len(cells)),

		circleMask: generateCircleMask(borderWidth),

		borderWidth: borderWidth,
	}

	// Random "intuitive", but substantially smaller hint than full len(cells)
	bs.playerCells[0] = make(map[Offset]struct{}, len(cells)/borderWidth)
	bs.playerCells[1] = make(map[Offset]struct{}, len(cells)/borderWidth)

	// A questionable... I guess... way to get any element from a map
	minX, minY, maxX, maxY := 0, 0, 0, 0
	for cell := range cells {
		minX, maxX = cell.X, cell.Y
		minY, maxY = cell.Y, cell.Y
		break
	}

	for cell, state := range cells {
		bs.board[cell] = state
		switch state {
		case CellUnoccupied:
			bs.unoccupiedCells[cell] = struct{}{}
		case CellP1:
			bs.playerCells[0][cell] = struct{}{}
		case CellP2:
			bs.playerCells[1][cell] = struct{}{}
		default:
			panic(fmt.Sprintf("new board state from cells: encountered an invalid cell at %v (state=%v)", cell, state))
		}

		if cell.X < minX {
			minX = cell.X
		} else if cell.X > maxX {
			maxX = cell.X
		}

		if cell.Y < minY {
			minY = cell.Y
		} else if cell.Y > maxY {
			maxY = cell.Y
		}
	}

	bs.boardBound = Rect{X: minX, Y: minY, W: maxX - minX + 1, H: maxY - minY + 1}

	return bs
}

func (bs *BoardState) Clone() *BoardState {
	newBs := &BoardState{
		board:           make(map[Offset]CellState, len(bs.board)),
		unoccupiedCells: make(map[Offset]struct{}, len(bs.unoccupiedCells)),

		delta:       make([]boardDelta, len(bs.delta)),
		moveHistory: make([]PlayerMove, len(bs.moveHistory)),

		circleMask: bs.circleMask,

		borderWidth: bs.borderWidth,
		boardBound:  bs.boardBound,
	}

	newBs.playerCells[0] = make(map[Offset]struct{}, len(bs.playerCells[0]))
	newBs.playerCells[1] = make(map[Offset]struct{}, len(bs.playerCells[1]))

	for k, v := range bs.board {
		newBs.board[k] = v
	}

	for k, v := range bs.unoccupiedCells {
		newBs.unoccupiedCells[k] = v
	}

	for k, v := range bs.playerCells[0] {
		newBs.playerCells[0][k] = v
	}

	for k, v := range bs.playerCells[1] {
		newBs.playerCells[1][k] = v
	}

	copy(newBs.moveHistory, bs.moveHistory)
	copy(newBs.delta, bs.delta)

	return newBs
}

func (bs *BoardState) BorderWidth() int {
	return bs.borderWidth
}

func (bs *BoardState) Delta() []boardDelta {
	return bs.delta
}

func (bs *BoardState) MoveHistoryCopy() []PlayerMove {
	historyCopy := make([]PlayerMove, len(bs.moveHistory))

	copy(historyCopy, bs.moveHistory)

	return historyCopy
}

// Do not modify the result
func (bs *BoardState) AllCells() map[Offset]CellState {
	return bs.board
}

func (bs *BoardState) PlayerCells() [2]map[Offset]struct{} {
	return bs.playerCells
}

func (bs *BoardState) UnoccupiedCells() map[Offset]struct{} {
	return bs.unoccupiedCells
}

func (bs *BoardState) MoveCount() int {
	return len(bs.moveHistory)
}

func (bs *BoardState) Cell(pos Offset) CellState {
	state, available := bs.board[pos]
	if !available {
		return CellUnavailable
	}

	return state
}

func (bs *BoardState) BoardBound() Rect {
	return bs.boardBound
}

func (bs *BoardState) LatestMove() PlayerMove {
	if len(bs.moveHistory) == 0 {
		panic("game state: get last move: no moves have been yet made")
	}

	return bs.moveHistory[len(bs.moveHistory)-1]
}

// Will turn an unavailable cell into an unoccupied cell.
// Panics if cell is already available.
func (bs *BoardState) MarkUnoccupied(pos Offset) {
	previousState, exists := bs.board[pos]
	if exists {
		panic(fmt.Sprintf("board state: add unoccupied cell at %v: the cell is already present (state=%d)", pos, previousState))
	}

	bs.board[pos] = CellUnoccupied
	bs.unoccupiedCells[pos] = struct{}{}
}

func (bs *BoardState) MarkCell(pos Offset, player PlayerID) {
	if bs.board[pos] != CellUnoccupied {
		panic(fmt.Sprintf("Trying to mark an occupied cell at %#v", pos))
	}

	bs.board[pos] = CellState(player)
	delete(bs.unoccupiedCells, pos)
	bs.playerCells[player][pos] = struct{}{}

	bs.moveHistory = append(bs.moveHistory, PlayerMove{pos, player})

	delta := boardDelta{}
	delta.Cells = append(delta.Cells, cellDelta{Cell: pos, NewState: CellState(player)})
	delta.OldBoardBound = bs.boardBound

	// Update board bounding rectangle
	borderOffset := Offset{X: bs.borderWidth, Y: bs.borderWidth}
	newCellsBoundingRect := NewRectFromOffsets(pos.Sub(borderOffset), borderOffset.ScaleUp(2).AddXY(1, 1))

	bs.boardBound = bs.boardBound.GrowToContainRect(newCellsBoundingRect)
	delta.NewBoardBound = bs.boardBound

	// Create new available cells
	for _, ds := range bs.circleMask {
		curCell := pos.Add(ds)

		_, available := bs.board[curCell]
		if !available {
			bs.MarkUnoccupied(curCell)
			delta.Cells = append(delta.Cells, cellDelta{curCell, CellUnoccupied})
		}
	}

	bs.delta = append(bs.delta, delta)
}

func (bs *BoardState) UndoLastMove() {
	if bs.MoveCount() == 0 {
		panic("board state: undo last move: no move to undo")
	}

	bs.moveHistory = bs.moveHistory[:len(bs.moveHistory)-1]

	lastDelta := bs.delta[len(bs.delta)-1]
	bs.delta = bs.delta[:len(bs.delta)-1]

	bs.boardBound = lastDelta.OldBoardBound

	for _, dcell := range lastDelta.Cells {
		switch dcell.NewState {
		case CellUnoccupied:
			delete(bs.board, dcell.Cell)
			delete(bs.unoccupiedCells, dcell.Cell)

		case CellP1, CellP2:
			bs.board[dcell.Cell] = CellUnoccupied
			delete(bs.playerCells[dcell.NewState], dcell.Cell)
			bs.unoccupiedCells[dcell.Cell] = struct{}{}

		default:
			panic(fmt.Sprintf("board state: undo last move: invalid cell delta at %v, new state=%v", dcell.Cell, dcell.NewState))
		}
	}
}
