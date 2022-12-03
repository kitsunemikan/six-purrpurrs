package ai

import (
	"math/rand"
	"time"

	"github.com/kitsunemikan/ttt-cli/game"
	. "github.com/kitsunemikan/ttt-cli/geom"
)

type RandomPlayer struct{}

func NewRandomPlayer() game.PlayerAgent {
	rand.Seed(time.Now().Unix())
	return &RandomPlayer{}
}

func (p *RandomPlayer) MakeMove(g *game.GameState) Offset {
	time.Sleep(50 * time.Millisecond)

	var validMoves []Offset
	for cell := range g.Board.UnoccupiedCells() {
		validMoves = append(validMoves, cell)
	}

	if len(validMoves) == 0 {
		panic("RandomPlayer: no unoccupied cells left")
	}

	moveID := rand.Intn(len(validMoves))
	return validMoves[moveID]
}
