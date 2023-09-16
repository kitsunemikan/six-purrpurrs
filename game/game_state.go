package game

import (
	. "github.com/kitsunemikan/six-purrpurrs/geom"
)

var solutionOffsets = []Offset{{X: 1, Y: 0}, {X: 1, Y: 1}, {X: 0, Y: 1}, {X: 1, Y: -1}}

type GameOptions struct {
	Border       int
	StrikeLength int
}

type GameState struct {
	// TODO: make private
	Conf GameOptions

	Board      *BoardState
	StrikeStat *StrikeSet
	solution   []Offset

	winner       PlayerID
	winnerMoveID int
}

func NewGame(conf GameOptions) *GameState {
	g := &GameState{
		Conf:       conf,
		Board:      NewBoardState(conf.Border),
		StrikeStat: NewStrikeSet(),
	}

	return g
}

// XXX: I didn't think when wrote this. Expects that no players yet won.
func (g *GameState) Clone() *GameState {
	strikeSet := NewStrikeSet()
	for _, move := range g.Board.moveHistory {
		strikeSet.MakeMove(move.Cell, move.Player)
	}

	return &GameState{
		Conf:       g.Conf,
		Board:      g.Board.Clone(),
		StrikeStat: strikeSet,
	}
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

func (g *GameState) CheckSolutionsAt(pos Offset) []Offset {
	strikes := g.StrikeStat.StrikesThrough(pos)

	for strikeID := range strikes {
		if strikes[strikeID].Len >= g.Conf.StrikeLength {
			return strikes[strikeID].AsCells()
		}
	}
	return nil
}

func (g *GameState) CandidateCellsAt(cell Offset, player PlayerID) []Offset {
	var candidates []Offset

	for _, dir := range StrikeDirs {
		// Forward direction
		afterCell := cell.Add(dir.Offset())
		afterStrike := g.StrikeStat.StrikesThrough(afterCell)[dir.FixedID]
		if afterStrike.Player == player {
			cells := afterStrike.AsCells()
			candidates = append(candidates, cells...)
		}

		// Backward direction
		beforeCell := cell.Sub(dir.Offset())
		beforeStrike := g.StrikeStat.StrikesThrough(beforeCell)[dir.FixedID]
		if beforeStrike.Player == player {
			cells := beforeStrike.AsCells()
			candidates = append(candidates, cells...)
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
	g.Board.MarkCell(pos, player)
	g.StrikeStat.MakeMove(pos, player)

	g.solution = g.CheckSolutionsAt(pos)
	if g.solution != nil {
		g.winner = player
		g.winnerMoveID = g.Board.MoveCount() - 1
	}
}

func (g *GameState) UndoLastMove() {
	lastMove := g.Board.LatestMove()
	g.StrikeStat.MarkUnoccupied(lastMove.Cell)

	g.Board.UndoLastMove()

	if g.Board.MoveCount() == g.winnerMoveID {
		g.solution = nil
		g.winnerMoveID = 0
	}
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
