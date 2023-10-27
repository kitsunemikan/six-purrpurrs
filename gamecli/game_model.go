package gamecli

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/kitsunemikan/six-purrpurrs/game"
	. "github.com/kitsunemikan/six-purrpurrs/geom"
)

// GameModel is a bubble that displays game board
// and additionally highlights:
// * Strike candidates
// * Latest marked cell
// * A winning strike
type GameModel struct {
	Game  *game.GameState
	Board BoardModel
}

func (m GameModel) View() string {
	// Repeated application of lipgloss render will produce incorrect results
	// Instead, we'll store the exact style for the cell in a map
	styledCells := make(map[Offset]lipgloss.Style)

	// Highlight candidates, if selection is visible
	selection := m.Board.Selection()
	if m.Board.SelectionVisible && m.Game.Cell(selection) == game.CellUnoccupied {
		candidates := m.Game.CandidatesAroundFor(selection, m.Board.CurrentPlayer)

		for _, cell := range candidates {
			styledCells[cell] = m.Board.Theme.CandidateCellStyle
		}
	}

	// Highlight last enemy cell
	if m.Game.MoveNumber() > 1 {
		latestMove := m.Game.LatestMove()
		styledCells[latestMove.Cell] = m.Board.Theme.LastEnemyCellStyle
	}

	// Victory cells
	for _, cell := range m.Game.VictoriousStrike() {
		styledCells[cell] = m.Board.Theme.VictoryCellStyle
	}

	m.Board.ForcedHighlight = styledCells

	return m.Board.View()
}
