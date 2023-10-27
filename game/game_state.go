package game

import (
	. "github.com/kitsunemikan/six-purrpurrs/geom"
)

var solutionOffsets = []Offset{{X: 1, Y: 0}, {X: 1, Y: 1}, {X: 0, Y: 1}, {X: 1, Y: -1}}

type GameOptions struct {
	Border int

	Victory VictoryChecker
}

type GameState struct {
	Board      *BoardState
	StrikeStat *StrikeSet

	victory VictoryChecker
}

func NewGame(conf GameOptions) *GameState {
	g := &GameState{
		Board:      NewBoardState(conf.Border),
		StrikeStat: NewStrikeSet(),

		victory: conf.Victory,
	}

	return g
}

func (g *GameState) VictoryChecker() VictoryChecker {
	return g.victory
}

func (g *GameState) Clone() *GameState {
	strikeSet := NewStrikeSet()
	for _, move := range g.Board.moveHistory {
		strikeSet.MakeMove(move.Cell, move.Player)
	}

	return &GameState{
		Board:      g.Board.Clone(),
		StrikeStat: strikeSet,
		victory:    g.victory.Clone(),
	}
}

func (g *GameState) MoveNumber() int {
	return g.Board.MoveCount() + 1
}

func (g *GameState) Cell(pos Offset) CellState {
	return g.Board.Cell(pos)
}

func (g *GameState) MoveHistoryCopy() []PlayerMove {
	return g.Board.MoveHistoryCopy()
}

func (g *GameState) CandidatesAroundFor(cell Offset, player PlayerID) []Offset {
	return g.victory.CandidatesAroundFor(g.StrikeStat, cell, player)
}

func (g *GameState) Over() bool {
	return g.victory.Reached()
}

func (g *GameState) BoardBound() Rect {
	return g.Board.BoardBound()
}

func (g *GameState) MarkCell(pos Offset, player PlayerID) {
	g.Board.MarkCell(pos, player)
	g.StrikeStat.MakeMove(pos, player)

	g.victory.CheckAt(g.StrikeStat, pos)
}

func (g *GameState) UndoLastMove() {
	lastMove := g.Board.LatestMove()
	g.StrikeStat.MarkUnoccupied(lastMove.Cell)

	g.Board.UndoLastMove()

	if g.victory.Reached() {
		g.victory.Reset()
	}
}

func (g *GameState) Winner() PlayerID {
	return g.victory.VictoriousPlayer()
}

func (g *GameState) VictoriousStrike() []Offset {
	return g.victory.VictoriousStrike()
}

func (g *GameState) LatestMove() PlayerMove {
	return g.Board.LatestMove()
}
