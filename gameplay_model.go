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

	Camera     Offset
	ScreenSize Offset

	Selection Offset

	MoveCommitted bool
	CurrentPlayer PlayerID
	Players       map[PlayerID]PlayerAgent
}

func NewGameplayModel(game *GameState, players map[PlayerID]PlayerAgent, screenSize Offset) GameplayModel {
	return GameplayModel{
		Game:    game,
		Players: players,

		Camera:     screenSize.ScaleDown(-2),
		ScreenSize: screenSize,

		Selection:     Offset{0, 0},
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
			m.Selection.X -= 1
			m.Camera.X -= 1
			return m, nil

		case "right", "l":
			m.Selection.X += 1
			m.Camera.X += 1
			return m, nil

		case "up", "k":
			m.Selection.Y -= 1
			m.Camera.Y -= 1
			return m, nil

		case "down", "j":
			m.Selection.Y += 1
			m.Camera.Y += 1
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
			return GameOverModel{m.Game, m.ScreenSize}, nil
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
	cliBoard := make(map[Offset]string, m.ScreenSize.Area())

	for y := 0; y < m.ScreenSize.Y; y++ {
		for x := 0; x < m.ScreenSize.X; x++ {
			curCell := m.Camera.Add(Offset{x, y})
			cliBoard[curCell] = m.Game.PlayerToken(m.Game.Cell(curCell))
		}
	}

	if m.IsLocalPlayerTurn() && m.Game.Cell(m.Selection) == CellUnoccupied {
		candidates := m.Game.CandidateCellsAt(m.Selection, m.CurrentPlayer)

		for _, cell := range candidates {
			if !cell.IsInsideRect(m.Camera, m.Camera.Add(m.ScreenSize)) {
				continue
			}

			cliBoard[cell] = candidateStyle.Render(cliBoard[cell])
		}
	}

	var view strings.Builder
	UnoccupiedToken := m.Game.PlayerToken(CellUnoccupied)
	for y := 0; y < m.ScreenSize.Y; y++ {
		for x := 0; x < m.ScreenSize.X; x++ {
			curCell := m.Camera.Add(Offset{x, y})

			leftSide := " "
			rightSide := " "
			if m.IsLocalPlayerTurn() && curCell.IsEqual(m.Selection) {
				leftSide = "["
				rightSide = "]"
			}

			if m.Game.Cell(curCell) != CellUnoccupied {
				view.WriteString(inactiveTextStyle.Render(leftSide))
			} else {
				view.WriteString(leftSide)
			}

			if m.Game.Cell(curCell) == CellUnoccupied {
				view.WriteString(UnoccupiedToken)
			} else {
				view.WriteString(cliBoard[curCell])
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
