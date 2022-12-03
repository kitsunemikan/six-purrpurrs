package game_test

import (
	"testing"
	"testing/quick"

	"github.com/kitsunemikan/ttt-cli/ai"
	"github.com/kitsunemikan/ttt-cli/game"
)

func TestBoardStateRevertability(t *testing.T) {
	randomPlayer := ai.NewRandomPlayer()

	assertion := func(moveCount uint8) bool {
		board := game.NewBoardState(3)

		boardHistory := make([]*game.BoardState, moveCount)

		player := game.P1
		for i := 0; i < int(moveCount); i++ {
			boardHistory[i] = board.Clone()

			nextMove := randomPlayer.MakeMove(board)
			board.MarkCell(nextMove, player)

			player = player.Other()
		}

		for i := int(moveCount) - 1; i >= 0; i-- {
			board.UndoLastMove()

			if !game.BoardStatesEqual(board, boardHistory[i]) {
				return true
			}
		}

		return true
	}

	if err := quick.Check(assertion, nil); err != nil {
		checkErr := err.(*quick.CheckError)
		t.Errorf("#%d: failed with input %d", checkErr.Count, checkErr.In[0])
	}
}
