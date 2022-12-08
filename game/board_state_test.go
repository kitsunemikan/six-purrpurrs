package game_test

import (
	"testing"
	"testing/quick"

	"github.com/kitsunemikan/six-purrpurrs/ai"
	"github.com/kitsunemikan/six-purrpurrs/game"
	"github.com/kitsunemikan/six-purrpurrs/game/gametest"
	"github.com/kitsunemikan/six-purrpurrs/gamecli"
	"github.com/sanity-io/litter"
)

func TestBoardStateRevertability(t *testing.T) {
	randomPlayer := ai.NewRandomPlayer()

	assertion := func(moveCount uint8) bool {
		moveCount /= 4
		if moveCount < 2 {
			moveCount = 2
		}

		board := game.NewBoardState(3)

		boardHistory := make([]*game.BoardState, moveCount+1)

		player := game.P1
		for i := 0; i < int(moveCount); i++ {
			boardHistory[i] = board.Clone()

			nextMove := randomPlayer.MakeMove(board)
			board.MarkCell(nextMove, player)

			player = player.Other()
		}

		// To print parent board if the first undo fails
		boardHistory[moveCount] = board.Clone()

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
				boardModel := gamecli.NewBoardModel(boardHistory[i+1].BoardBound().Dimensions(), 0)
				boardModel.Board = boardHistory[i+1]
				boardModel.Theme = &gamecli.DefaultBoardTheme

				parentBoard := boardModel.CenterOnBoard().View()

				t.Logf("case failed on move %d/%d:\nparent\n%v\n%v", i, moveCount, parentBoard, err)
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
