package gamecli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/kitsunemikan/ttt-cli/game"
	"github.com/kitsunemikan/ttt-cli/gamecli/keymap"
	. "github.com/kitsunemikan/ttt-cli/geom"
)

type ReplayModelOptions struct {
	Game   *game.GameState
	Board  BoardModel
	Help   help.Model
	Parent tea.Model
}

type ReplayModel struct {
	game     *game.GameState
	board    BoardModel
	moves    []game.PlayerMove
	nextMove int

	help   help.Model
	parent tea.Model
}

func NewReplayModel(config ReplayModelOptions) ReplayModel {
	moveHistory := config.Game.MoveHistoryCopy()
	return ReplayModel{
		game:     config.Game,
		board:    config.Board,
		moves:    moveHistory,
		nextMove: len(moveHistory),

		help:   config.Help,
		parent: config.Parent,
	}
}

func (m ReplayModel) Init() tea.Cmd {
	return nil
}

func (m ReplayModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keymap.Replay.Quit):
			return m.parent, nil

		case key.Matches(msg, keymap.Replay.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil

		case key.Matches(msg, keymap.Replay.Left):
			m.board = m.board.MoveCameraBy(Offset{X: -1, Y: 0})
			return m, nil

		case key.Matches(msg, keymap.Replay.Right):
			m.board = m.board.MoveCameraBy(Offset{X: 1, Y: 0})
			return m, nil

		case key.Matches(msg, keymap.Replay.Up):
			m.board = m.board.MoveCameraBy(Offset{X: 0, Y: -1})
			return m, nil

		case key.Matches(msg, keymap.Replay.Down):
			m.board = m.board.MoveCameraBy(Offset{X: 0, Y: 1})
			return m, nil

		case key.Matches(msg, keymap.Replay.Forward):
			if m.nextMove == len(m.moves) {
				return m, nil
			}

			move := m.moves[m.nextMove]
			m.game.MarkCell(move.Cell, move.ID)
			m.nextMove++

			m.board = m.board.MoveSelectionTo(move.Cell).CenterOnSelection()

		case key.Matches(msg, keymap.Replay.Rewind):
			if m.nextMove == 0 {
				return m, nil
			}

			m.game.UndoLastMove()
			m.nextMove--
			m.board = m.board.CenterOnBoard()
		}
	}

	return m, nil
}

func (m ReplayModel) View() string {
	var view strings.Builder

	view.WriteString(m.board.View())
	view.WriteString("\n")
	view.WriteString(fmt.Sprintf("Move %d/%d\n", m.nextMove, len(m.moves)))
	view.WriteString("\n")
	view.WriteString(m.help.View(keymap.Replay))

	return view.String()
}
