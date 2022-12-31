package game_test

import (
	"fmt"
	"testing"

	"github.com/kitsunemikan/six-purrpurrs/ai"
	"github.com/kitsunemikan/six-purrpurrs/game"
	"github.com/kitsunemikan/six-purrpurrs/geom"
)

func BenchmarkGameBoardRandomPlayers(b *testing.B) {
	opt := game.GameOptions{
		Border:       7,
		StrikeLength: 6,
	}

	cases := []struct {
		moveCount int
	}{
		{moveCount: 100},
		{moveCount: 1000},
		{moveCount: 10000},
	}

	for _, data := range cases {
		name := fmt.Sprintf("%d random moves", data.moveCount)
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()

				gameState := game.NewGame(opt)
				p1 := ai.NewRandomPlayer()
				p2 := ai.NewRandomPlayer()
				currentPlayer := game.P1

				b.StartTimer()
				for moveID := 0; moveID < data.moveCount; moveID++ {
					var chosenCell geom.Offset
					switch currentPlayer {
					case game.P1:
						chosenCell = p1.MakeMove(gameState.Board)
					case game.P2:
						chosenCell = p2.MakeMove(gameState.Board)
					}

					gameState.MarkCell(chosenCell, currentPlayer)

					currentPlayer = currentPlayer.Other()
				}
			}
		})
	}
}
