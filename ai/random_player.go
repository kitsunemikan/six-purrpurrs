package ai

import (
	"github.com/kitsunemikan/six-purrpurrs/game"
	. "github.com/kitsunemikan/six-purrpurrs/geom"
)

type RandomPlayer struct{}

func NewRandomPlayer() game.PlayerAgent {
	return &RandomPlayer{}
}

func (p *RandomPlayer) MakeMove(g *game.GameState) Offset {
	for cell := range g.Board.UnoccupiedCells() {
		return cell
	}

	panic("random player: no unoccupied cells were present at all!")
}
