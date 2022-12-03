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

func (p *RandomPlayer) MakeMove(b *game.BoardState) Offset {
	validMoves := make([]Offset, 0, len(b.UnoccupiedCells()))
	for cell := range b.UnoccupiedCells() {
		validMoves = append(validMoves, cell)
	}

	if len(validMoves) == 0 {
		panic("RandomPlayer: no unoccupied cells left")
	}

	moveID := rand.Intn(len(validMoves))
	return validMoves[moveID]
}
