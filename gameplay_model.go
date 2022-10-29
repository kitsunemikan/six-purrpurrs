package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// A bubbletea event
type PlayerMoveMsg struct {
	ChosenCell Offset
}

type GameplayModel struct {
	Game *GameState

	Selection Offset

	MoveCommitted bool
	CurrentPlayer PlayerID
	Players       map[PlayerID]PlayerAgent
}

func NewGameplayModel(game *GameState, players map[PlayerID]PlayerAgent) GameplayModel {
	return GameplayModel{
		Game:    game,
		Players: players,

		Selection:     Offset{game.BoardSize().X / 2, game.BoardSize().Y / 2},
		CurrentPlayer: 1,
	}
}

func (m *GameplayModel) AwaitMove(player PlayerID) tea.Cmd {
	return func() tea.Msg {
		move := m.Players[player].MakeMove(m.Game)
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
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

		switch msg.String() {
		case "left", "h":
			if m.Selection.X > 0 {
				m.Selection.X -= 1
			}
			return m, nil

		case "right", "l":
			if m.Selection.X < m.Game.BoardSize().X-1 {
				m.Selection.X += 1
			}
			return m, nil

		case "up", "k":
			if m.Selection.Y > 0 {
				m.Selection.Y -= 1
			}
			return m, nil

		case "down", "j":
			if m.Selection.Y < m.Game.BoardSize().Y-1 {
				m.Selection.Y += 1
			}
			return m, nil

		case "enter", " ":
			localPlayer, local := m.Players[m.CurrentPlayer].(*LocalPlayer)
			if !local {
				break
			}

			if m.MoveCommitted {
				return m, nil
			}

			if m.Game.Cell(m.Selection) != Unoccupied {
				return m, nil
			}

			localPlayer.CommitMove(m.Selection)

			m.MoveCommitted = true
			return m, nil
		}

	case PlayerMoveMsg:
		m.Game.MarkCell(msg.ChosenCell, m.CurrentPlayer)
		m.MoveCommitted = false

		if m.Game.Over() {
			return GameOverModel{m.Game}, nil
		}

		if m.CurrentPlayer == m.Game.LastPlayer() {
			m.CurrentPlayer = 1
		} else {
			m.CurrentPlayer++
		}

		return m, m.AwaitMove(m.CurrentPlayer)
	}

	return m, nil
}

func (m GameplayModel) View() string {
	cliBoard := m.Game.BoardToStrings()

	if m.Game.Cell(m.Selection) == Unoccupied {
		candidates := m.Game.CandidateCellsAt(m.Selection, m.CurrentPlayer)

		for _, cell := range candidates {
			cliBoard[cell] = candidateStyle.Render(cliBoard[cell])
		}
	}

	var view strings.Builder
	UnoccupiedToken := m.Game.PlayerToken(Unoccupied)
	for y := 0; y < m.Game.BoardSize().Y; y++ {
		for x := 0; x < m.Game.BoardSize().X; x++ {
			leftSide := " "
			rightSide := " "
			if x == m.Selection.X && y == m.Selection.Y {
				leftSide = "["
				rightSide = "]"
			}

			curCell := Offset{x, y}

			if m.Game.Cell(curCell) != Unoccupied {
				view.WriteString(inactiveTextStyle.Render(leftSide))
			} else {
				view.WriteString(leftSide)
			}

			if m.Game.Cell(curCell) == Unoccupied {
				view.WriteString(UnoccupiedToken)
			} else {
				view.WriteString(cliBoard[Offset{x, y}])
			}

			if m.Game.Cell(curCell) != Unoccupied {
				view.WriteString(inactiveTextStyle.Render(rightSide))
			} else {
				view.WriteString(rightSide)
			}
		}
		view.WriteByte('\n')
	}

	view.WriteByte('\n')

	if m.IsLocalPlayerTurn() {
		view.WriteString("Current player: ")
		view.WriteString(m.Game.PlayerToken(m.CurrentPlayer))
	} else {
		view.WriteString("Awaiting player ")
		view.WriteString(m.Game.PlayerToken(m.CurrentPlayer))
		view.WriteString(" move...")
	}

	view.WriteByte('\n')

	return view.String()
}
