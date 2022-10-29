package main

type LocalPlayer struct {
	moves chan Offset
}

func NewLocalPlayer() PlayerAgent {
	return &LocalPlayer{
		moves: make(chan Offset),
	}
}

func (p *LocalPlayer) MakeMove(g *GameState) Offset {
	return <-p.moves
}

func (p *LocalPlayer) CommitMove(pos Offset) {
	p.moves <- pos
}
