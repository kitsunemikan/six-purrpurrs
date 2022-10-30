package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type GameOverModel struct {
	Game       *GameState
	ScreenSize Offset
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
	cliBoard := m.Game.BoardToStrings()

	for _, cell := range m.Game.Solution() {
		cliBoard[cell] = winCellStyle.Render(cliBoard[cell])
	}

	var view strings.Builder
	UnoccupiedToken := m.Game.PlayerToken(CellUnoccupied)
	for y := 0; y < m.ScreenSize.Y; y++ {
		for x := 0; x < m.ScreenSize.X; x++ {
			view.WriteRune(' ')

			curCell := Offset{x, y}
			if m.Game.Cell(curCell) == CellUnoccupied {
				view.WriteString(UnoccupiedToken)
			} else {
				view.WriteString(cliBoard[Offset{x, y}])
			}

			view.WriteRune(' ')
		}
		view.WriteByte('\n')
	}

	view.WriteByte('\n')

	if m.Game.Solution() == nil {
		view.WriteString("A draw...")
	} else {
		view.WriteString(m.Game.PlayerToken(m.Game.Winner()))
		view.WriteString(" wins!")
	}

	view.WriteString(fmt.Sprintf("\n\nTotal number of moves made: %d\n\nPress any key to exit...\n", m.Game.MoveNumber))

	return view.String()
}
