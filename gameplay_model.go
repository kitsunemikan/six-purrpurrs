package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// A bubbletea event
type PlayerMoveMsg struct {
	ChosenCell Offset
}

type GameplayModel struct {
	Game *GameState

	Camera Rect

	Selection Offset

	MoveCommitted bool
	CurrentPlayer PlayerID
	Players       map[PlayerID]PlayerAgent

	cameraBound Rect
}

func NewGameplayModel(game *GameState, players map[PlayerID]PlayerAgent, screenSize Offset) GameplayModel {
	return GameplayModel{
		Game:    game,
		Players: players,

		Camera:      NewRectFromOffsets(screenSize.ScaleDown(-2), screenSize),
		cameraBound: game.BoardBound(),

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
			if m.Selection.X > m.cameraBound.X {
				m.Selection.X -= 1
				m.Camera = m.Camera.CenterOn(m.Selection).SnapInto(m.cameraBound)
			}
			return m, nil

		case "right", "l":
			if m.Selection.X < m.cameraBound.X+m.cameraBound.W-1 {
				m.Selection.X += 1
				m.Camera = m.Camera.CenterOn(m.Selection).SnapInto(m.cameraBound)
			}
			return m, nil

		case "up", "k":
			if m.Selection.Y > m.cameraBound.Y {
				m.Selection.Y -= 1
				m.Camera = m.Camera.CenterOn(m.Selection).SnapInto(m.cameraBound)
			}
			return m, nil

		case "down", "j":
			if m.Selection.Y < m.cameraBound.Y+m.cameraBound.H-1 {
				m.Selection.Y += 1
				m.Camera = m.Camera.CenterOn(m.Selection).SnapInto(m.cameraBound)
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
			return GameOverModel{m.Game, m.Camera}, nil
		}

		if m.CurrentPlayer == m.Game.LastPlayer() {
			m.CurrentPlayer = 1
		} else {
			m.CurrentPlayer++
		}

		m.cameraBound = m.Game.BoardBound()
		m.Camera = m.Camera.CenterOn(m.Selection).SnapInto(m.cameraBound)

		return m, m.AwaitMove(m.CurrentPlayer)
	}

	return m, nil
}

func (m GameplayModel) View() string {
	cliBoard := m.Game.BoardToStrings(m.Camera)

	if m.IsLocalPlayerTurn() && m.Game.Cell(m.Selection) == CellUnoccupied {
		candidates := m.Game.CandidateCellsAt(m.Selection, m.CurrentPlayer)

		for _, cell := range candidates {
			if !m.Camera.IsOffsetInside(cell) {
				continue
			}

			cliBoard[cell] = candidateStyle.Render(cliBoard[cell])
		}
	}

	if m.Game.MoveNumber() > 1 {
		latestMove := m.Game.LatestMove()
		cliBoard[latestMove.Cell] = enemyCellStyle.Render(cliBoard[latestMove.Cell])
	}

	var view strings.Builder
	UnoccupiedToken := m.Game.PlayerToken(CellUnoccupied)
	for y := 0; y < m.Camera.H; y++ {
		for x := 0; x < m.Camera.W; x++ {
			curCell := m.Camera.ToWorldXY(x, y)

			leftSide := " "
			rightSide := ""
			if m.IsLocalPlayerTurn() {
				if curCell.IsEqual(m.Selection) {
					leftSide = "["
					rightSide = "]"
				} else if curCell.IsEqual(m.Selection.SubXY(1, 0)) {
					rightSide = ""
				} else if curCell.IsEqual(m.Selection.AddXY(1, 0)) {
					leftSide = ""
				}
			}

			if m.Game.Cell(curCell) != CellUnoccupied {
				view.WriteString(inactiveTextStyle.Render(leftSide))
				view.WriteString(cliBoard[curCell])
				view.WriteString(inactiveTextStyle.Render(rightSide))
			} else {
				view.WriteString(leftSide)
				view.WriteString(UnoccupiedToken)
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

	view.WriteString(fmt.Sprintf("\nCamera bound: %v | Camera: %v", m.cameraBound, m.Camera))
	view.WriteByte('\n')

	return view.String()
}
