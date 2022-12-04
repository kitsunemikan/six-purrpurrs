package gamecli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/kitsunemikan/ttt-cli/game"
)

type GameOverModel struct {
	Game  *game.GameState
	Board BoardModel
	Help  help.Model
}

func (m GameOverModel) Init() tea.Cmd {
	return nil
}

func (m GameOverModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, GlobalGameOverKeymap.Quit):
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m GameOverModel) View() string {
	m.Board.SelectionVisible = false
	m.Help.ShowAll = false

	var view strings.Builder

	view.WriteString(m.Board.View())
	view.WriteByte('\n')

	if m.Game.Solution() == nil {
		view.WriteString("A draw...")
	} else {
		view.WriteString(m.Board.Theme.PlayerCells[m.Game.Winner()])
		view.WriteString(" wins!")
	}

	view.WriteString(fmt.Sprintf("\n\nTotal number of moves made: %d\n\n", m.Game.MoveNumber()-1))

	view.WriteString(m.Help.View(GlobalGameOverKeymap))
	view.WriteByte('\n')

	return view.String()
}
