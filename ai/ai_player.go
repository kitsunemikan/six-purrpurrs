package ai

import (
	"math/rand"
	"time"

	"github.com/kitsunemikan/ttt-cli/game"
	. "github.com/kitsunemikan/ttt-cli/geom"
)

type BoardRank struct {
	Me   []int
	Them []int
}

type moveOutcome struct {
	Cell     Offset
	BestRank BoardRank
}

type AIPlayer struct {
	id   game.PlayerID
	rand *rand.Rand
}

func NewDefaultAIPlayer(id game.PlayerID) game.PlayerAgent {
	return &AIPlayer{
		id:   id,
		rand: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (p *AIPlayer) minimax(g *game.BoardState, player game.PlayerID, depth int) moveOutcome {
	// TODO: add UnoccupiedCells method
	validMoves := make([]Offset, 0, len(g.AllCells()))
	for cell, state := range g.AllCells() {
		if state == game.CellUnoccupied {
			validMoves = append(validMoves, cell)
		}
	}

	ranks := make([]moveOutcome, len(validMoves))
	for i, move := range validMoves {
		if depth == 1 {
			// Calculate rank with move
		} else {
			// Change g
			g.MarkCell(move, player)
			ranks[i] = p.minimax(g, player.Other(), depth-1)
		}
	}

	return moveOutcome{
		Cell: Offset{X: 0, Y: 0},
	}
}

func (p *AIPlayer) MakeMove(b *game.BoardState) Offset {
	board := make(map[Offset]game.CellState, len(b.AllCells()))
	for k, v := range b.AllCells() {
		board[k] = v
	}

	// TODO: we don't need some stuff in GameState, refactor into BoardState
	// TODO: clone method
	boardCopy := b.Clone()

	bestMove := p.minimax(boardCopy, p.id, 4)

	return bestMove.Cell
}
