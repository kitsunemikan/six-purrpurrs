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
	time.Sleep(50 * time.Millisecond)

	var validMoves []Offset
	for cell, state := range g.Board {
		if state == CellUnoccupied {
			validMoves = append(validMoves, cell)
		}
	}

	if len(validMoves) == 0 {
		panic("RandomPlayer: no unoccupied cells left")
	}

	moveID := rand.Intn(len(validMoves))
	return validMoves[moveID]
}
