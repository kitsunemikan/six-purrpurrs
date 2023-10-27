package ai

import (
	"math/rand"
	"time"

	"github.com/kitsunemikan/six-purrpurrs/game"
	"github.com/kitsunemikan/six-purrpurrs/geom"
)

type ObstructivePlayer struct {
	Me   game.PlayerID
	rand *rand.Rand
}

func NewObstructivePlayer(id game.PlayerID) game.PlayerAgent {
	source := rand.NewSource(time.Now().UnixMicro())
	return &ObstructivePlayer{
		Me:   id,
		rand: rand.New(source),
	}
}

func (p *ObstructivePlayer) MakeMove(g *game.GameState) geom.Offset {
	// Collect shifts
	dirs := make([]int, len(game.StrikeDirs))
	for i := range dirs {
		dirs[i] = i
	}

	for i := 0; i < len(dirs); i++ {
		swapID := p.rand.Intn(len(dirs)-1) + 1
		dirs[0], dirs[swapID] = dirs[swapID], dirs[0]
	}

	for opponentCell := range g.Board.PlayerCells()[p.Me.Other()] {
		for i := 0; i < len(dirs); i++ {
			cell := opponentCell.Add(game.StrikeDirs[dirs[i]].Offset())
			if _, ok := g.Board.UnoccupiedCells()[cell]; ok {
				return cell
			}

			cell = opponentCell.Sub(game.StrikeDirs[dirs[i]].Offset())
			if _, ok := g.Board.UnoccupiedCells()[cell]; ok {
				return cell
			}
		}
	}

	// If all opponent's cells are obstructed, choose unoccupied at random
	for cell := range g.Board.UnoccupiedCells() {
		return cell
	}

	panic("obstructing player: no unoccupied cells were present at all!")
}
