package gamecli

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/kitsunemikan/ttt-cli/game"
	. "github.com/kitsunemikan/ttt-cli/geom"
)

type BoardModel struct {
	Theme *BoardTheme
	Board *game.BoardState

	camera    Camera
	selection Offset

	SelectionVisible bool
	CurrentPlayer    game.PlayerID

	ForcedHighlight map[Offset]lipgloss.Style
}

func NewBoardModel(cameraSize Offset, trackDepth int) BoardModel {
	return BoardModel{
		camera: Camera{
			// Board extends to negative integers, so board's center is at (0,0),
			// and not (screenWidth/2, screenHeight/2)
			View:       NewRectFromOffsets(cameraSize.ScaleDown(-2), cameraSize),
			TrackDepth: trackDepth,
		},
	}
}

func (m BoardModel) MoveSelectionBy(ds Offset) BoardModel {
	cameraBound := m.Board.BoardBound()
	newSelection := m.selection.Add(ds)

	if newSelection.IsInsideRect(cameraBound) {
		m.selection = newSelection
	}

	return m
}

func (m BoardModel) MoveSelectionTo(pos Offset) BoardModel {
	cameraBound := m.Board.BoardBound()

	if pos.IsInsideRect(cameraBound) {
		m.selection = pos
	}

	return m
}

func (m BoardModel) MoveCameraBy(ds Offset) BoardModel {
	m.camera = m.camera.Move(ds).SnapIntoRect(m.Board.BoardBound())

	return m
}

func (m BoardModel) NudgeCameraTo(pos Offset) BoardModel {
	m.camera = m.camera.NudgeTo(pos)
	return m
}

func (m BoardModel) SnapSelectionIntoCamera() BoardModel {
	m.selection = m.selection.SnapIntoRect(m.camera.InnerView())

	return m
}

func (m BoardModel) NudgeToSelection() BoardModel {
	m.camera = m.camera.NudgeTo(m.selection).SnapIntoRect(m.Board.BoardBound())
	return m
}

func (m BoardModel) CenterOnBoard() BoardModel {
	m.camera = m.camera.SnapIntoRect(m.Board.BoardBound())
	return m
}

func (m BoardModel) Selection() Offset {
	return m.selection
}

func (m BoardModel) ModelDimensions() Offset {
	// 2 * camera dimensions, because we artificially stretch
	// the board, so that it appears more square when rendered
	return Offset{X: 2 * m.camera.View.W, Y: 2 * m.camera.View.H}
}

func (m BoardModel) View() string {
	cliBoard := m.Theme.BoardToText(m.Board.AllCells(), m.camera.View)

	// Repeated application of lipgloss render will produce incorrect results
	// Instead, we'll store the exact style for the cell in a map
	styledCells := make(map[Offset]lipgloss.Style)

	// Forced highlights (e.g., for pretty test fail outputs)
	for cell, style := range m.ForcedHighlight {
		styledCells[cell] = style
	}

	// Apply styles
	for pos, str := range cliBoard {
		style, special := styledCells[pos]
		if special {
			cliBoard[pos] = style.Render(str)
			continue
		}

		cellState := m.Board.Cell(pos)
		if cellState == game.CellUnavailable || cellState == game.CellUnoccupied {
			continue
		}

		cliBoard[pos] = m.Theme.PlayerCellStyles[cellState].Render(str)
	}

	var view strings.Builder
	for y := 0; y < m.camera.View.H; y++ {
		for x := 0; x < m.camera.View.W; x++ {
			curCell := m.camera.View.ToWorldXY(x, y)

			leftSide := " "
			rightSide := ""
			if m.SelectionVisible {
				if curCell.IsEqual(m.selection) {
					leftSide = "["
					rightSide = "]"
				} else if curCell.IsEqual(m.selection.SubXY(1, 0)) {
					// Because the right side will be '[' for the selected cell
					rightSide = ""
				} else if curCell.IsEqual(m.selection.AddXY(1, 0)) {
					// Because the left side will be ']' for the selected cell
					// Though this line is unnecessary...
					leftSide = ""
				}
			}

			if m.Board.Cell(curCell) != game.CellUnoccupied {
				view.WriteString(m.Theme.SelectionInactiveStyle.Render(leftSide))
				view.WriteString(cliBoard[curCell])
				view.WriteString(m.Theme.SelectionInactiveStyle.Render(rightSide))
			} else {
				view.WriteString(leftSide)
				view.WriteString(cliBoard[curCell])
				view.WriteString(rightSide)
			}
		}
		view.WriteByte('\n')
	}

	return view.String()
}
