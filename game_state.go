package main

import (
	"fmt"
)

type PlayerID int

const (
	CellUnavailable PlayerID = -1
	CellUnoccupied  PlayerID = 0
)

var solutionOffsets = []Offset{{1, 0}, {1, 1}, {0, 1}, {1, -1}}

type PlayerMove struct {
	Cell Offset
	ID   PlayerID
}

type GameOptions struct {
	Border       int
	StrikeLength int
	PlayerTokens []string
}

type GameState struct {
	// TODO: make private
	Conf GameOptions

	Board    map[Offset]PlayerID
	solution []Offset
	winner   PlayerID

	circleMask  []Offset
	moveHistory []PlayerMove
	boardBound  Rect
}

func NewGame(conf GameOptions) *GameState {
	g := &GameState{
		Conf:       conf,
		Board:      make(map[Offset]PlayerID),
		circleMask: make([]Offset, 0, conf.Border*conf.Border),
		boardBound: Rect{X: -conf.Border, Y: -conf.Border, W: 2*conf.Border + 1, H: 2*conf.Border + 1},
	}

	// Generate circle mask
	for dx := -g.Conf.Border; dx <= g.Conf.Border; dx++ {
		for dy := -g.Conf.Border; dy <= g.Conf.Border; dy++ {
			ds := Offset{dx, dy}
			if !ds.IsInsideCircle(g.Conf.Border) {
				continue
			}

			g.circleMask = append(g.circleMask, ds)
		}
	}

	// Mark initial available cells
	for _, ds := range g.circleMask {
		g.Board[ds] = CellUnoccupied
	}

	return g
}

func (g *GameState) LastPlayer() PlayerID {
	return PlayerID(len(g.Conf.PlayerTokens) - 2)
}

func (g *GameState) MoveNumber() int {
	return len(g.moveHistory)
}

func (g *GameState) Cell(pos Offset) PlayerID {
	state, available := g.Board[pos]
	if !available {
		return CellUnavailable
	}

	return state
}

func (g *GameState) PlayerToken(player PlayerID) string {
	if player < -1 || player > g.LastPlayer() {
		panic(fmt.Sprintf("model: player token for ID=%v: out of range (LastPlayerID=%v)", player, g.LastPlayer()))
	}

	return g.Conf.PlayerTokens[int(player)+1]
}

func (g *GameState) IsInsideBoard(pos Offset) bool {
	// TODO: based on another border

	return true
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
			curCell := pos.Add(dir.ScaleUp(i))
			if g.Board[curCell] == player {
				candidates = append(candidates, curCell)
			} else {
				break
			}
		}

		for i := 1; i < g.Conf.StrikeLength; i++ {
			curCell := pos.Add(dir.ScaleUp(-i))
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
	return false
	// return g.MoveNumber == g.Conf.BoardSize.X*g.Conf.BoardSize.Y
}

func (g *GameState) Over() bool {
	return g.solution != nil || g.NoMoreMoves()
}

func (g *GameState) BoardBound() Rect {
	return g.boardBound
}

func (g *GameState) MarkCell(pos Offset, player PlayerID) {
	if g.Board[pos] != CellUnoccupied {
		panic(fmt.Sprintf("Trying to mark an occupied cell at %#v", pos))
	}

	if g.solution != nil {
		panic(fmt.Sprintf("Trying to mark a cell at %#v, when the game is already over", pos))
	}

	g.Board[pos] = player
	g.moveHistory = append(g.moveHistory, PlayerMove{pos, player})

	// Update board bounding rectangle

	borderOffset := Offset{g.Conf.Border, g.Conf.Border}
	newCellsBoundingRect := NewRectFromOffsets(pos.Sub(borderOffset), borderOffset.ScaleUp(2).AddXY(1, 1))
	g.boardBound = g.boardBound.GrowToContainRect(newCellsBoundingRect)

	// Create new available cells
	for _, ds := range g.circleMask {
		curCell := pos.Add(ds)

		_, available := g.Board[curCell]
		if !available {
			g.Board[curCell] = CellUnoccupied
		}
	}

	g.solution = g.CheckSolutionsAt(pos, player)
	if g.solution != nil {
		g.winner = player
	}
}

func (g *GameState) Winner() PlayerID {
	return g.winner
}

func (g *GameState) BoardToStrings(camera Rect) map[Offset]string {
	cliBoard := make(map[Offset]string, camera.Area())

	for x := 0; x < camera.W; x++ {
		for y := 0; y < camera.H; y++ {
			curCell := camera.ToWorldXY(x, y)
			cliBoard[curCell] = g.PlayerToken(g.Cell(curCell))
		}
	}
	return cliBoard
}

func (g *GameState) Solution() []Offset {
	return g.solution
}
