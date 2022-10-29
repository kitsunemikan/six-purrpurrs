package main

import (
	"math/rand"
	"time"
)

type RandomPlayer struct{}

func NewRandomPlayer() PlayerAgent {
	rand.Seed(time.Now().Unix())
	return &RandomPlayer{}
}

func (p *RandomPlayer) MakeMove(g *GameState) Offset {
	time.Sleep(500 * time.Millisecond)

	var validMoves []Offset
	for x := 0; x < g.BoardSize().X; x++ {
		for y := 0; y < g.BoardSize().Y; y++ {
			if g.Cell(Offset{x, y}) == Unoccupied {
				validMoves = append(validMoves, Offset{x, y})
			}
		}
	}

	if len(validMoves) == 0 {
		panic("RandomPlayer: no unoccupied cells left")
	}

	moveID := rand.Intn(len(validMoves))
	return validMoves[moveID]
}
