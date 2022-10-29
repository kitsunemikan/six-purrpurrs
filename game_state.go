package main

import (
	"fmt"
)

type PlayerID int

const Unoccupied PlayerID = 0

var solutionOffsets = []Offset{{1, 0}, {1, 1}, {0, 1}, {1, -1}}

type GameOptions struct {
	BoardSize    Offset
	StrikeLength int
	PlayerTokens []string
}

type GameState struct {
	Conf GameOptions

	Board      map[Offset]PlayerID
	solution   []Offset
	winner     PlayerID
	MoveNumber int
}

func NewGame(conf GameOptions) *GameState {
	return &GameState{
		Conf:  conf,
		Board: make(map[Offset]PlayerID),
	}
}

func (g *GameState) LastPlayer() PlayerID {
	return PlayerID(len(g.Conf.PlayerTokens) - 1)
}

func (g *GameState) Cell(pos Offset) PlayerID {
	return g.Board[pos]
}

func (g *GameState) BoardSize() Offset {
	return g.Conf.BoardSize
}

func (g *GameState) PlayerToken(player PlayerID) string {
	if player < 0 || player > g.LastPlayer() {
		panic(fmt.Sprintf("model: player token for ID=%v: out of range (LastPlayerID=%v)", player, g.LastPlayer()))
	}

	return g.Conf.PlayerTokens[int(player)]
}

func (g *GameState) IsInsideBoard(pos Offset) bool {
	return pos.X >= 0 && pos.X < g.Conf.BoardSize.X && pos.Y >= 0 && pos.Y < g.Conf.BoardSize.Y
}

func (g *GameState) CheckSolutionsAt(pos Offset, player PlayerID) []Offset {
	solution := make([]Offset, 0, g.Conf.StrikeLength)

	for _, dir := range solutionOffsets {
		solution = solution[:0]

		if len(solution) != 0 {
			panic(42)
		}

		for i := 0; i < g.Conf.StrikeLength; i++ {
			curCell := Offset{pos.X + i*dir.X, pos.Y + i*dir.Y}
			if g.Board[curCell] != player {
				break
			}

			solution = append(solution, curCell)
		}

		for i := 1; i < g.Conf.StrikeLength; i++ {
			curCell := Offset{pos.X - i*dir.X, pos.Y - i*dir.Y}
			if g.Board[curCell] != player {
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

	if g.Board[pos] == player {
		candidates = append(candidates, pos)
	}

	for _, dir := range solutionOffsets {
		for i := 1; i < g.Conf.StrikeLength; i++ {
			curCell := pos.Add(dir.Scale(i))
			if g.Board[curCell] == player {
				candidates = append(candidates, curCell)
			} else {
				break
			}
		}

		for i := 1; i < g.Conf.StrikeLength; i++ {
			curCell := pos.Add(dir.Scale(-i))
			if g.Board[curCell] == player {
				candidates = append(candidates, curCell)
			} else {
				break
			}
		}
	}

	return candidates
}

func (g *GameState) NoMoreMoves() bool {
	return g.MoveNumber == g.Conf.BoardSize.X*g.Conf.BoardSize.Y
}

func (g *GameState) Over() bool {
	return g.solution != nil || g.NoMoreMoves()
}

func (g *GameState) MarkCell(pos Offset, player PlayerID) {
	if g.Board[pos] != Unoccupied {
		panic(fmt.Sprintf("Trying to mark an occupied cell at %#v", pos))
	}

	if g.solution != nil {
		panic(fmt.Sprintf("Trying to mark a cell at %#v, when the game is already over", pos))
	}

	g.Board[pos] = player
	g.MoveNumber++

	g.solution = g.CheckSolutionsAt(pos, player)
	if g.solution != nil {
		g.winner = player
	}
}

func (g *GameState) Winner() PlayerID {
	return g.winner
}

func (g *GameState) BoardToStrings() map[Offset]string {
	cliBoard := make(map[Offset]string)
	for pos, player := range g.Board {
		cliBoard[pos] = string(g.PlayerToken(player))
	}

	return cliBoard
}

func (g *GameState) Solution() []Offset {
	return g.solution
}
