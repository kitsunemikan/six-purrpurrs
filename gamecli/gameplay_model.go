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

// A bubbletea event
type PlayerMoveMsg struct {
	ChosenCell Offset
}

type GameplayModelConfig struct {
	Game    *game.GameState
	Players []game.PlayerAgent

	Theme      *BoardTheme
	ScreenSize Offset
}

type GameplayModel struct {
	Game  *game.GameState
	board BoardModel
	help  help.Model

	MoveCommitted bool
	CurrentPlayer game.PlayerID
	Players       []game.PlayerAgent
}

func NewGameplayModel(config GameplayModelConfig) GameplayModel {
	if config.Game == nil {
		panic("new gameplay model: game state is nil")
	}

	if config.Theme == nil {
		panic("new gameplay model: board theme is nil")
	}

	if len(config.Players) == 0 {
		panic("new gameplay model: no player agents specified")
	}

	if config.ScreenSize.IsZero() {
		panic("new gameplay model: zero screen size")
	}

	board := NewBoardModel(config.ScreenSize)
	board.Board = config.Game.Board
	board.Theme = config.Theme
	board.CurrentPlayer = game.P1

	help := help.New()
	help.Styles = HelpStyle
	return GameplayModel{
		Game:    config.Game,
		Players: config.Players,
		board:   board,
		help:    help,

		CurrentPlayer: game.P1,
	}
}

func (m *GameplayModel) AwaitMove(player game.PlayerID) tea.Cmd {
	return func() tea.Msg {
		move := m.Players[player].MakeMove(m.Game.Board)
		return PlayerMoveMsg{move}
	}
}

func (m *GameplayModel) IsLocalPlayerTurn() bool {
	_, local := m.Players[m.CurrentPlayer].(*LocalPlayer)
	return local
}

func (m GameplayModel) Init() tea.Cmd {
	return m.AwaitMove(m.CurrentPlayer)
}

func (m GameplayModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can it can gracefully truncate
		// its view as needed
		m.help.Width = msg.Width

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keymap.Gameplay.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, keymap.Gameplay.Quit):
			return m, tea.Quit
		}

		if !m.IsLocalPlayerTurn() {
			break
		}

		switch {
		case key.Matches(msg, keymap.Gameplay.Left):
			m.board = m.board.MoveSelectionBy(Offset{X: -1, Y: 0}).NudgeToSelection()
			return m, nil

		case key.Matches(msg, keymap.Gameplay.Right):
			m.board = m.board.MoveSelectionBy(Offset{X: 1, Y: 0}).NudgeToSelection()
			return m, nil

		case key.Matches(msg, keymap.Gameplay.Up):
			m.board = m.board.MoveSelectionBy(Offset{X: 0, Y: -1}).NudgeToSelection()
			return m, nil

		case key.Matches(msg, keymap.Gameplay.Down):
			m.board = m.board.MoveSelectionBy(Offset{X: 0, Y: 1}).NudgeToSelection()
			return m, nil

		case key.Matches(msg, keymap.Gameplay.Select):
			if m.MoveCommitted {
				return m, nil
			}

			if m.Game.Cell(m.board.Selection()) != game.CellUnoccupied {
				return m, nil
			}

			localPlayer := m.Players[m.CurrentPlayer].(*LocalPlayer)
			localPlayer.CommitMove(m.board.Selection())

			m.MoveCommitted = true
			return m, nil
		}

	case PlayerMoveMsg:
		m.Game.MarkCell(msg.ChosenCell, m.CurrentPlayer)
		m.MoveCommitted = false

		if m.Game.Over() {
			return GameOverModel{Game: m.Game, Board: m.board, Help: m.help}, nil
		}

		m.CurrentPlayer = m.CurrentPlayer.Other()
		m.board.CurrentPlayer = m.CurrentPlayer

		m.board = m.board.NudgeCameraTo(msg.ChosenCell).SnapSelectionIntoCamera()

		return m, m.AwaitMove(m.CurrentPlayer)
	}

	return m, nil
}

func (m GameplayModel) View() string {
	m.board.SelectionVisible = m.IsLocalPlayerTurn()

	var view strings.Builder

	gameModel := GameModel{
		Game:  m.Game,
		Board: m.board,
	}

	view.WriteString(gameModel.View())
	view.WriteByte('\n')

	if m.IsLocalPlayerTurn() {
		view.WriteString("Current player: ")
		view.WriteString(m.board.Theme.PlayerCells[m.CurrentPlayer])
	} else {
		view.WriteString("Awaiting player ")
		view.WriteString(m.board.Theme.PlayerCells[m.CurrentPlayer])
		view.WriteString(" move...")
	}

	// view.WriteString(fmt.Sprintf("\nCamera bound: %v | Camera: %v", m.cameraBound, m.Camera))
	view.WriteString("\n\n")

	view.WriteString(m.help.View(keymap.Gameplay))
	view.WriteByte('\n')

	return view.String()
}
