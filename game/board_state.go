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

	delta       []boardDelta
	moveHistory []PlayerMove

	// Precalculated disk offsets
	circleMask []Offset

	borderWidth int
	boardBound  Rect
}

func NewBoardState(borderWidth int) *BoardState {
	bs := &BoardState{
		board:           make(map[Offset]CellState),
		unoccupiedCells: make(map[Offset]struct{}),

		circleMask: make([]Offset, 0, borderWidth*borderWidth),

		borderWidth: borderWidth,
		boardBound:  Rect{X: -borderWidth, Y: -borderWidth, W: 2*borderWidth + 1, H: 2*borderWidth + 1},
	}

	bs.playerCells[0] = make(map[Offset]struct{})
	bs.playerCells[1] = make(map[Offset]struct{})

	// Generate circle mask
	for dx := -borderWidth; dx <= borderWidth; dx++ {
		for dy := -borderWidth; dy <= borderWidth; dy++ {
			ds := Offset{X: dx, Y: dy}
			if !ds.IsInsideCircle(borderWidth) {
				continue
			}

			bs.circleMask = append(bs.circleMask, ds)
		}
	}

	// Mark initial available cells
	for _, ds := range bs.circleMask {
		bs.MarkUnoccupied(ds)
	}

	return bs
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
}
