package gamecli

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/kitsunemikan/ttt-cli/game"
	"github.com/kitsunemikan/ttt-cli/gamecli/keymap"
	. "github.com/kitsunemikan/ttt-cli/geom"
)

type ReplayModel struct {
	Game  *game.GameState
	Board BoardModel
	Moves []game.PlayerMove
	Help  help.Model

	Parent tea.Model
}

func (m ReplayModel) Init() tea.Cmd {
	return nil
}

func (m ReplayModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keymap.Replay.Quit):
			return m.Parent, nil

		case key.Matches(msg, keymap.Replay.Help):
			m.Help.ShowAll = !m.Help.ShowAll
			return m, nil

		case key.Matches(msg, keymap.Gameplay.Left):
			m.Board = m.Board.MoveCameraBy(Offset{X: -1, Y: 0})
			return m, nil

		case key.Matches(msg, keymap.Gameplay.Right):
			m.Board = m.Board.MoveCameraBy(Offset{X: 1, Y: 0})
			return m, nil

		case key.Matches(msg, keymap.Gameplay.Up):
			m.Board = m.Board.MoveCameraBy(Offset{X: 0, Y: -1})
			return m, nil

		case key.Matches(msg, keymap.Gameplay.Down):
			m.Board = m.Board.MoveCameraBy(Offset{X: 0, Y: 1})
			return m, nil
		}
	}

	return m, nil
}

func (m ReplayModel) View() string {
	var view strings.Builder

	view.WriteString(m.Board.View())
	view.WriteString("\n")
	view.WriteString(m.Help.View(keymap.Replay))

	return view.String()
}
