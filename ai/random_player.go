package ai

import (
	"github.com/kitsunemikan/six-purrpurrs/game"
	. "github.com/kitsunemikan/six-purrpurrs/geom"
)

type RandomPlayer struct{}

func NewRandomPlayer() game.PlayerAgent {
	return &RandomPlayer{}
}

func (p *RandomPlayer) MakeMove(b *game.BoardState) Offset {
	for cell := range b.UnoccupiedCells() {
		return cell
	}

	panic("random player: no unoccupied cells were present at all!")
}
