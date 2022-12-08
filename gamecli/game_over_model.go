package gamecli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/kitsunemikan/six-purrpurrs/game"
	"github.com/kitsunemikan/six-purrpurrs/gamecli/keymap"
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
		case key.Matches(msg, keymap.GameOver.Quit):
			return m, tea.Quit

		case key.Matches(msg, keymap.GameOver.WatchReplay):
			replayModel := NewReplayModel(ReplayModelOptions{
				Game:   m.Game,
				Board:  m.Board,
				Help:   m.Help,
				Parent: m,
			})

			cmd := replayModel.Init()
			return replayModel, cmd
		}
	}

	return m, nil
}

func (m GameOverModel) View() string {
	m.Board.SelectionVisible = false
	m.Help.ShowAll = false

	var view strings.Builder

	gameModel := GameModel{
		Game:  m.Game,
		Board: m.Board,
	}

	view.WriteString(gameModel.View())
	view.WriteByte('\n')

	if m.Game.Solution() == nil {
		view.WriteString("A draw...")
	} else {
		view.WriteString(m.Board.Theme.PlayerCells[m.Game.Winner()])
		view.WriteString(" wins!")
	}

	view.WriteString(fmt.Sprintf("\n\nTotal number of moves made: %d\n\n", m.Game.MoveNumber()-1))

	view.WriteString(m.Help.View(keymap.GameOver))
	view.WriteByte('\n')

	return view.String()
}
