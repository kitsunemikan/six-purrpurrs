package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type GameOverModel struct {
	Game   *GameState
	Theme  *BoardTheme
	Camera Rect
}

func (m GameOverModel) Init() tea.Cmd {
	return nil
}

func (m GameOverModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	}

	return m, nil
}

func (m GameOverModel) View() string {
	cliBoard := m.Theme.BoardToText(m.Game.AllCells(), m.Camera)

	for _, cell := range m.Game.Solution() {
		cliBoard[cell] = winCellStyle.Render(cliBoard[cell])
	}

	var view strings.Builder
	for y := 0; y < m.Camera.H; y++ {
		for x := 0; x < m.Camera.W; x++ {
			curCell := m.Camera.ToWorldXY(x, y)

			view.WriteString(" ")
			if m.Game.Cell(curCell) != CellUnoccupied {
				view.WriteString(cliBoard[curCell])
			} else {
				view.WriteString(m.Theme.UnoccupiedCell)
			}
		}
		view.WriteByte('\n')
	}

	view.WriteByte('\n')

	if m.Game.Solution() == nil {
		view.WriteString("A draw...")
	} else {
		view.WriteString(m.Theme.PlayerCells[m.Game.Winner()-1])
		view.WriteString(" wins!")
	}

	view.WriteString(fmt.Sprintf("\n\nTotal number of moves made: %d\n\nPress any key to exit...\n", m.Game.MoveNumber()-1))

	return view.String()
}
