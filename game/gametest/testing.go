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

	diffColor := lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("1"))

	boardModel := gamecli.NewBoardModel(got.BoardBound().Dimensions(), 0)
	boardModel.Theme = &gamecli.DefaultBoardTheme

	boardDiff, same := cellBoardDiff(diffColor, got.AllCells(), want.AllCells())
	if !same {
		gotBoard, wantBoard := drawDiffBoards(boardModel, boardDiff, got.AllCells(), want.AllCells())
		return fmt.Errorf("all-cell boards are different:\ngot\n%v\nwant\n%v\n", gotBoard, wantBoard)
	}

	gotUnoccupied := cellSetToCellBoard(got.UnoccupiedCells(), game.CellUnoccupied)
	wantUnoccupied := cellSetToCellBoard(want.UnoccupiedCells(), game.CellUnoccupied)
	boardDiff, same = cellBoardDiff(diffColor, gotUnoccupied, wantUnoccupied)
	if !same {
		gotBoard, wantBoard := drawDiffBoards(boardModel, boardDiff, gotUnoccupied, wantUnoccupied)
		return fmt.Errorf("unoccupied cell boards are different:\ngot\n%v\nwant\n%v\n", gotBoard, wantBoard)
	}

	gotCellsP1 := cellSetToCellBoard(got.PlayerCells()[0], game.CellP1)
	wantCellsP1 := cellSetToCellBoard(want.PlayerCells()[0], game.CellP1)

	boardDiff, same = cellBoardDiff(diffColor, gotCellsP1, wantCellsP1)
	if !same {
		gotBoard, wantBoard := drawDiffBoards(boardModel, boardDiff, gotCellsP1, wantCellsP1)
		return fmt.Errorf("P1 cell boards are different:\ngot\n%v\nwant\n%v\n", gotBoard, wantBoard)
	}

	gotCellsP2 := cellSetToCellBoard(got.PlayerCells()[1], game.CellP2)
	wantCellsP2 := cellSetToCellBoard(want.PlayerCells()[1], game.CellP2)

	boardDiff, same = cellBoardDiff(diffColor, gotCellsP2, wantCellsP2)
	if !same {
		gotBoard, wantBoard := drawDiffBoards(boardModel, boardDiff, gotCellsP2, wantCellsP2)
		return fmt.Errorf("P2 cell boards are different:\ngot\n%v\nwant\n%v\n", gotBoard, wantBoard)
	}

	gotDelta := got.Delta()
	wantDelta := want.Delta()
	if !reflect.DeepEqual(gotDelta, wantDelta) {
		return fmt.Errorf("deltas are different:\ngot\n%v\nwant\n%v\n", litter.Sdump(gotDelta), litter.Sdump(wantDelta))
	}

	gotHistory := got.MoveHistoryCopy()
	wantHistory := want.MoveHistoryCopy()
	if !reflect.DeepEqual(gotHistory, wantHistory) {
		return fmt.Errorf("deltas are different:\ngot\n%v\nwant\n%v\n", litter.Sdump(gotHistory), litter.Sdump(wantHistory))
	}

	if !got.BoardBound().IsEqual(want.BoardBound()) {
		return fmt.Errorf("board bounds are different:\ngot\n%v\nwant\n%v\n", litter.Sdump(got.BoardBound()), litter.Sdump(want.BoardBound()))
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

func drawDiffBoards(model gamecli.BoardModel, diff map[Offset]lipgloss.Style, got, want map[Offset]game.CellState) (string, string) {
	model.ForcedHighlight = diff

	gotBoard := game.NewBoardStateFromCells(1, got)
	model.Board = gotBoard
	gotStr := model.CenterOnBoard().View()

	wantBoard := game.NewBoardStateFromCells(1, want)
	model.Board = wantBoard
	wantStr := model.CenterOnBoard().View()

	return gotStr, wantStr
}
