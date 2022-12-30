package ai

import (
	"math/rand"
	"time"

	"github.com/kitsunemikan/six-purrpurrs/game"
	. "github.com/kitsunemikan/six-purrpurrs/geom"
)

type RandomPlayer struct{}

func NewRandomPlayer() game.PlayerAgent {
	rand.Seed(time.Now().Unix())
	return &RandomPlayer{}
}

func (p *RandomPlayer) MakeMove(b *game.BoardState) Offset {
	moveID := rand.Intn(len(b.UnoccupiedCells()))
	for cell := range b.UnoccupiedCells() {
		if moveID == 0 {
			return cell
		}

		moveID--
	}

	panic("random player: exhausted all unoccupied cells")
}
