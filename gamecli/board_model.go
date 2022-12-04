package gamecli

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/kitsunemikan/ttt-cli/game"
	. "github.com/kitsunemikan/ttt-cli/geom"
)

type BoardModel struct {
	Theme *BoardTheme
	Game  *game.GameState

	camera    Rect
	selection Offset

	SelectionVisible bool
	CurrentPlayer    game.PlayerID
}

func NewBoardModel(cameraSize Offset) BoardModel {
	return BoardModel{
		// Board extends to negative integers, so board's center is at (0,0),
		// and not (screenWidth/2, screenHeight/2)
		camera: NewRectFromOffsets(cameraSize.ScaleDown(-2), cameraSize),
	}
}

func (m BoardModel) MoveSelectionBy(ds Offset) BoardModel {
	cameraBound := m.Game.BoardBound()
	newSelection := m.selection.Add(ds)

	if newSelection.IsInsideRect(cameraBound) {
		m.selection = newSelection
	}

	return m
}

func (m BoardModel) MoveCameraBy(ds Offset) BoardModel {
	cameraBound := m.Game.BoardBound()
	newCamera := m.camera.Move(ds)

	if newCamera.IsInsideRect(cameraBound) {
		m.camera = newCamera
	}

	return m
}

func (m BoardModel) CenterOnSelection() BoardModel {
	m.camera = m.camera.CenterOn(m.selection).SnapInto(m.Game.BoardBound())
	return m
}

func (m BoardModel) Selection() Offset {
	return m.selection
}

func (m BoardModel) View() string {
	cliBoard := m.Theme.BoardToText(m.Game.AllCells(), m.camera)

	// Repeated application of lipgloss render will produce incorrect results
	// Instead, we'll store the exact style for the cell in a map
	styledCells := make(map[Offset]lipgloss.Style)

	// Highlight candidates, if selection is visible
	if m.SelectionVisible && m.Game.Cell(m.selection) == game.CellUnoccupied {
		candidates := m.Game.CandidateCellsAt(m.selection, m.CurrentPlayer)

		for _, cell := range candidates {
			if !cell.IsInsideRect(m.camera) {
				continue
			}

			styledCells[cell] = m.Theme.CandidateCellStyle
		}
	}

	// Highlight last enemy cell
	if m.Game.MoveNumber() > 1 {
		latestMove := m.Game.LatestMove()
		styledCells[latestMove.Cell] = m.Theme.LastEnemyCellStyle
	}

	// Victory cells
	for _, cell := range m.Game.Solution() {
		styledCells[cell] = m.Theme.VictoryCellStyle
	}

	// Apply styles
	for pos, str := range cliBoard {
		style, special := styledCells[pos]
		if special {
			cliBoard[pos] = style.Render(str)
			continue
		}

		cellState := m.Game.Cell(pos)
		if cellState == game.CellUnavailable || cellState == game.CellUnoccupied {
			continue
		}

		cliBoard[pos] = m.Theme.PlayerCellStyles[cellState].Render(str)
	}

	var view strings.Builder
	for y := 0; y < m.camera.H; y++ {
		for x := 0; x < m.camera.W; x++ {
			curCell := m.camera.ToWorldXY(x, y)

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

			if m.Game.Cell(curCell) != game.CellUnoccupied {
				view.WriteString(m.Theme.SelectionInactiveStyle.Render(leftSide))
				view.WriteString(cliBoard[curCell])
				view.WriteString(m.Theme.SelectionInactiveStyle.Render(rightSide))
			} else {
				view.WriteString(leftSide)
				view.WriteString(m.Theme.UnoccupiedCell)
				view.WriteString(rightSide)
			}
		}
		view.WriteByte('\n')
	}

	return view.String()
}
