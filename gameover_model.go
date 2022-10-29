package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type GameOverModel struct {
	Game *GameState
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
	for y := 0; y < m.Game.BoardSize().Y; y++ {
		for x := 0; x < m.Game.BoardSize().X; x++ {
			view.WriteRune(' ')

			curCell := Offset{x, y}
			if m.Game.Cell(curCell) == Unoccupied {
				view.WriteString(".")
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
		view.WriteRune(m.Game.PlayerToken(m.Game.Winner()))
		view.WriteString(" wins!")
	}

	view.WriteString("\n\nPress any key to exit...\n")

	return view.String()
}
