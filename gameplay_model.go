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

		if !m.IsLocalPlayerTurn() {
			break
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

			if m.MoveCommitted {
				return m, nil
			}

			if m.Game.Cell(m.Selection) != CellUnoccupied {
				return m, nil
			}

			localPlayer := m.Players[m.CurrentPlayer].(*LocalPlayer)
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
	cliBoard := make(map[Offset]string, m.Game.BoardSize().Area())
	// bottomRight is one-cell beyond valid range
	topLeft := m.Selection.Add(m.Game.BoardSize().ScaleDown(-2))
	bottomRight := m.Selection.Add(m.Game.BoardSize().Add(Offset{1, 1}).ScaleDown(2))

	for y := topLeft.Y; y < bottomRight.Y; y++ {
		for x := topLeft.X; x < bottomRight.X; x++ {
			curCell := Offset{x, y}
			cliBoard[curCell] = m.Game.PlayerToken(m.Game.Cell(curCell))
		}
	}

	if m.IsLocalPlayerTurn() && m.Game.Cell(m.Selection) == CellUnoccupied {
		candidates := m.Game.CandidateCellsAt(m.Selection, m.CurrentPlayer)

		for _, cell := range candidates {
			if !cell.IsInsideRect(topLeft, bottomRight) {
				continue
			}

			cliBoard[cell] = candidateStyle.Render(cliBoard[cell])
		}
	}

	var view strings.Builder
	UnoccupiedToken := m.Game.PlayerToken(CellUnoccupied)
	for y := topLeft.Y; y < bottomRight.Y; y++ {
		for x := topLeft.X; x < bottomRight.X; x++ {
			leftSide := " "
			rightSide := " "
			if m.IsLocalPlayerTurn() && x == m.Selection.X && y == m.Selection.Y {
				leftSide = "["
				rightSide = "]"
			}

			curCell := Offset{x, y}

			if m.Game.Cell(curCell) != CellUnoccupied {
				view.WriteString(inactiveTextStyle.Render(leftSide))
			} else {
				view.WriteString(leftSide)
			}

			if m.Game.Cell(curCell) == CellUnoccupied {
				view.WriteString(UnoccupiedToken)
			} else {
				view.WriteString(cliBoard[Offset{x, y}])
			}

			if m.Game.Cell(curCell) != CellUnoccupied {
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
