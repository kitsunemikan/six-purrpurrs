package game_test

import (
	"testing"
	"testing/quick"

	"github.com/kitsunemikan/ttt-cli/ai"
	"github.com/kitsunemikan/ttt-cli/game"
	"github.com/kitsunemikan/ttt-cli/game/gametest"
	"github.com/sanity-io/litter"
)

func TestBoardStateRevertability(t *testing.T) {
	randomPlayer := ai.NewRandomPlayer()

	assertion := func(moveCount uint8) bool {
		// moveCount /= 4
		moveCount = 20
		if moveCount < 2 {
			moveCount = 2
		}

		board := game.NewBoardState(3)

		boardHistory := make([]*game.BoardState, moveCount)

		player := game.P1
		for i := 0; i < int(moveCount); i++ {
			boardHistory[i] = board.Clone()

			nextMove := randomPlayer.MakeMove(board)
			board.MarkCell(nextMove, player)

			player = player.Other()
		}

		if len(board.PlayerCells()[0]) == 0 {
			t.Logf("case failed: after %d moves: P1 board is empty (player cells array = %v)", moveCount, litter.Sdump(board.PlayerCells()))
			return false
		}

		if len(board.PlayerCells()[1]) == 0 {
			t.Logf("case failed: after %d moves: P2 board is empty (player cells array = %v)", moveCount, litter.Sdump(board.PlayerCells()))
			return false
		}

		for i := int(moveCount) - 1; i >= 0; i-- {
			board.UndoLastMove()

			if err := gametest.BoardStatesEqual(board, boardHistory[i]); err != nil {
				t.Logf("case failed on move %d/%d: %v", i, moveCount, err)
				return false
			}
		}

		return true
	}

	if err := quick.Check(assertion, nil); err != nil {
		checkErr := err.(*quick.CheckError)
		t.Errorf("#%d: failed with input %d", checkErr.Count, checkErr.In[0])
	}
}
