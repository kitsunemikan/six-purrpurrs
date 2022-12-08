package game

import (
	"fmt"

	. "github.com/kitsunemikan/six-purrpurrs/geom"
)

// CellState has a property of, when positive, being equal
// to a correct PlayerID that occupies it. Special meanings
// are negative. This way, CellState can be used for indexing
// slices related to players without any additional manipulation
type CellState int

const (
	CellUnavailable CellState = iota - 2
	CellUnoccupied
	CellP1
	CellP2
)

func (cs CellState) IsOccupiedBy(player PlayerID) bool {
	return cs >= 0 && PlayerID(cs) == player
}

type PlayerID int

const (
	P1 PlayerID = iota
	P2
)

func (p PlayerID) Other() PlayerID {
	if p == P1 {
		return P2
	} else if p == P2 {
		return P1
	}

	panic(fmt.Sprintf("PlayerID: get other player: player is invalid (value=%d)", p))
}

var solutionOffsets = []Offset{{X: 1, Y: 0}, {X: 1, Y: 1}, {X: 0, Y: 1}, {X: 1, Y: -1}}

type PlayerMove struct {
	Cell Offset
	ID   PlayerID
}

type GameOptions struct {
	Border       int
	StrikeLength int
}

type GameState struct {
	// TODO: make private
	Conf GameOptions

	Board    *BoardState
	solution []Offset
	winner   PlayerID
}

func NewGame(conf GameOptions) *GameState {
	g := &GameState{
		Conf:  conf,
		Board: NewBoardState(conf.Border),
	}

	return g
}

func (g *GameState) AllCells() map[Offset]CellState {
	return g.Board.AllCells()
}

func (g *GameState) MoveNumber() int {
	return g.Board.MoveCount() + 1
}

func (g *GameState) Cell(pos Offset) CellState {
	return g.Board.Cell(pos)
}

func (g *GameState) IsInsideBoard(pos Offset) bool {
	// TODO: based on another border

	return true
}

func (g *GameState) MoveHistoryCopy() []PlayerMove {
	return g.Board.MoveHistoryCopy()
}

func (g *GameState) CheckSolutionsAt(pos Offset, player PlayerID) []Offset {
	solution := make([]Offset, 0, g.Conf.StrikeLength)

	for _, dir := range solutionOffsets {
		solution = solution[:0]

		if len(solution) != 0 {
			panic(42)
		}

		for i := 0; i < g.Conf.StrikeLength; i++ {
			curCell := Offset{X: pos.X + i*dir.X, Y: pos.Y + i*dir.Y}
			if !g.Board.Cell(curCell).IsOccupiedBy(player) {
				break
			}

			solution = append(solution, curCell)
		}

		for i := 1; i < g.Conf.StrikeLength; i++ {
			curCell := Offset{X: pos.X - i*dir.X, Y: pos.Y - i*dir.Y}
			if !g.Board.Cell(curCell).IsOccupiedBy(player) {
				break
			}

			solution = append(solution, curCell)
		}

		if len(solution) >= g.Conf.StrikeLength {
			return solution
		}
	}

	return nil
}

func (g *GameState) CandidateCellsAt(pos Offset, player PlayerID) []Offset {
	candidates := make([]Offset, 0, 2*len(solutionOffsets)*(g.Conf.StrikeLength-1)+1)

	if g.Board.Cell(pos).IsOccupiedBy(player) {
		candidates = append(candidates, pos)
	}

	for _, dir := range solutionOffsets {
		for i := 1; i < g.Conf.StrikeLength; i++ {
			curCell := pos.Add(dir.ScaleUp(i))
			if g.Board.Cell(curCell).IsOccupiedBy(player) {
				candidates = append(candidates, curCell)
			} else {
				break
			}
		}

		for i := 1; i < g.Conf.StrikeLength; i++ {
			curCell := pos.Add(dir.ScaleUp(-i))
			if g.Board.Cell(curCell).IsOccupiedBy(player) {
				candidates = append(candidates, curCell)
			} else {
				break
			}
		}
	}

	return candidates
}

func (g *GameState) NoMoreMoves() bool {
	return false
	// return g.MoveNumber == g.Conf.BoardSize.X*g.Conf.BoardSize.Y
}

func (g *GameState) Over() bool {
	return g.solution != nil || g.NoMoreMoves()
}

func (g *GameState) BoardBound() Rect {
	return g.Board.BoardBound()
}

func (g *GameState) MarkCell(pos Offset, player PlayerID) {
	if g.solution != nil {
		panic(fmt.Sprintf("Trying to mark a cell at %#v, when the game is already over", pos))
	}

	g.Board.MarkCell(pos, player)

	g.solution = g.CheckSolutionsAt(pos, player)
	if g.solution != nil {
		g.winner = player
	}
}

func (g *GameState) UndoLastMove() {
	g.Board.UndoLastMove()
	g.solution = nil
}

func (g *GameState) Winner() PlayerID {
	return g.winner
}

func (g *GameState) Solution() []Offset {
	return g.solution
}

func (g *GameState) LatestMove() PlayerMove {
	return g.Board.LatestMove()
}
