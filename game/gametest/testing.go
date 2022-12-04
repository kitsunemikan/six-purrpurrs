package gametest

import (
	"fmt"
	"reflect"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"

	"github.com/kitsunemikan/ttt-cli/game"
	"github.com/kitsunemikan/ttt-cli/gamecli"
	. "github.com/kitsunemikan/ttt-cli/geom"
	"github.com/sanity-io/litter"
)

func init() {
	lipgloss.SetColorProfile(termenv.ANSI256)
}

func BoardStatesEqual(got, want *game.BoardState) error {
	if got.BorderWidth() != want.BorderWidth() {
		return fmt.Errorf("got border width %v, want %v", got.BorderWidth(), want.BorderWidth())
	}

	boardModel := gamecli.NewBoardModel(got.BoardBound().Dimensions())
	boardModel.Theme = &gamecli.DefaultBoardTheme

	diffColor := lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("1"))
	boardDiff, same := cellBoardDiff(diffColor, got.AllCells(), want.AllCells())
	if !same {
		boardModel.ForcedHighlight = boardDiff

		boardModel.Board = got
		gotBoard := boardModel.CenterOnBoard().View()

		boardModel.Board = want
		wantBoard := boardModel.CenterOnBoard().View()
		return fmt.Errorf("compare board states: all cell boards are different:\ngot\n%v\nwant\n%v\n", gotBoard, wantBoard)
	}

	if !reflect.DeepEqual(got.UnoccupiedCells(), want.UnoccupiedCells()) {
		return fmt.Errorf("compare board states: unoccupied cell boards are different:\n%v\nand\n%v\n", litter.Sdump(got.UnoccupiedCells()), litter.Sdump(want.UnoccupiedCells()))
	}

	if !reflect.DeepEqual(got.PlayerCells()[0], want.PlayerCells()[0]) {
		return fmt.Errorf("compare board states: P1 cell boards are different:\n%v\nand\n%v\n", litter.Sdump(got.PlayerCells()[0]), litter.Sdump(want.PlayerCells()[0]))
	}

	gotCellsP2 := cellSetToCellBoard(got.PlayerCells()[1], game.CellP2)
	wantCellsP2 := cellSetToCellBoard(want.PlayerCells()[1], game.CellP2)
	boardDiff, same = cellBoardDiff(diffColor, gotCellsP2, wantCellsP2)
	if !same {
		return fmt.Errorf("compare board states: P2 cell boards are different:\n%v\nand\n%v\n", litter.Sdump(got.PlayerCells()[1]), litter.Sdump(want.PlayerCells()[1]))
	}

	gotDelta := got.Delta()
	wantDelta := want.Delta()
	if !reflect.DeepEqual(gotDelta, wantDelta) {
		return fmt.Errorf("compare board states: deltas are different:\n%v\nand\n%v\n", litter.Sdump(gotDelta), litter.Sdump(wantDelta))
	}

	gotHistory := got.MoveHistoryCopy()
	wantHistory := want.MoveHistoryCopy()
	if !reflect.DeepEqual(gotHistory, wantHistory) {
		return fmt.Errorf("compare board states: deltas are different:\n%v\nand\n%v\n", litter.Sdump(gotHistory), litter.Sdump(wantHistory))
	}

	if !got.BoardBound().IsEqual(want.BoardBound()) {
		return fmt.Errorf("compare board states: board bounds are different:\n%v\nand\n%v\n", litter.Sdump(got.BoardBound()), litter.Sdump(want.BoardBound()))
	}

	return nil
}

func cellBoardDiff(style lipgloss.Style, got map[Offset]game.CellState, want map[Offset]game.CellState) (map[Offset]lipgloss.Style, bool) {
	diff := make(map[Offset]lipgloss.Style)

	for gotCell, gotState := range got {
		wantState, exists := want[gotCell]

		if !exists {
			diff[gotCell] = style
			continue
		}

		if gotState != wantState {
			diff[gotCell] = style
		}
	}

	for wantCell := range want {
		_, exists := got[wantCell]

		if !exists {
			diff[wantCell] = style
		}
	}

	return diff, len(diff) == 0
}

func cellSetToCellBoard(set map[Offset]struct{}, state game.CellState) map[Offset]game.CellState {
	board := make(map[Offset]game.CellState, len(set))
	for cell := range set {
		board[cell] = state
	}

	return board
}
