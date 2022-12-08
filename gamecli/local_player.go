package gamecli

import (
	"github.com/kitsunemikan/six-purrpurrs/game"
	. "github.com/kitsunemikan/six-purrpurrs/geom"
)

type LocalPlayer struct {
	moves chan Offset
}

func NewLocalPlayer() game.PlayerAgent {
	return &LocalPlayer{
		moves: make(chan Offset),
	}
}

func (p *LocalPlayer) MakeMove(g *game.BoardState) Offset {
	return <-p.moves
}

func (p *LocalPlayer) CommitMove(pos Offset) {
	p.moves <- pos
}
