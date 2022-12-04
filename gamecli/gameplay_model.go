package gamecli

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/kitsunemikan/ttt-cli/game"
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
	Theme *BoardTheme

	Camera Rect

	Selection Offset

	MoveCommitted bool
	CurrentPlayer game.PlayerID
	Players       []game.PlayerAgent

	cameraBound Rect
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

	return GameplayModel{
		Game:    config.Game,
		Players: config.Players,
		Theme:   config.Theme,

		// Board extends to negative integers, so board's center is at (0,0),
		// and not (screenWidth/2, screenHeight/2)
		Camera:      NewRectFromOffsets(config.ScreenSize.ScaleDown(-2), config.ScreenSize),
		cameraBound: config.Game.BoardBound(),

		Selection:     Offset{0, 0},
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

			if m.Game.Cell(m.Selection) != game.CellUnoccupied {
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
			return GameOverModel{m.Game, m.Theme, m.Camera}, nil
		}

		m.CurrentPlayer = m.CurrentPlayer.Other()

		m.cameraBound = m.Game.BoardBound()
		m.Camera = m.Camera.CenterOn(m.Selection).SnapInto(m.cameraBound)

		return m, m.AwaitMove(m.CurrentPlayer)
	}

	return m, nil
}

func (m GameplayModel) View() string {
	cliBoard := m.Theme.BoardToText(m.Game.AllCells(), m.Camera)

	styledCells := make(map[Offset]lipgloss.Style)

	if m.IsLocalPlayerTurn() && m.Game.Cell(m.Selection) == game.CellUnoccupied {
		candidates := m.Game.CandidateCellsAt(m.Selection, m.CurrentPlayer)

		for _, cell := range candidates {
			if !m.Camera.IsOffsetInside(cell) {
				continue
			}

			styledCells[cell] = m.Theme.CandidateCellStyle
		}
	}

	if m.Game.MoveNumber() > 1 {
		latestMove := m.Game.LatestMove()
		styledCells[latestMove.Cell] = m.Theme.LastEnemyCellStyle
	}

	for pos, str := range cliBoard {
		style, special := styledCells[pos]
		if special {
			cliBoard[pos] = style.Render(str)
			continue
		}

		cellState := m.Game.Cell(pos)
		if cellState == game.CellUnavailable || cellState == game.CellUnoccupied {
			continue
		}

		cliBoard[pos] = m.Theme.PlayerCellStyles[cellState].Render(str)
	}

	var view strings.Builder
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

			if m.Game.Cell(curCell) != game.CellUnoccupied {
				view.WriteString(m.Theme.SelectionInactiveStyle.Render(leftSide))
				view.WriteString(cliBoard[curCell])
				view.WriteString(m.Theme.SelectionInactiveStyle.Render(rightSide))
			} else {
				view.WriteString(leftSide)
				view.WriteString(m.Theme.UnoccupiedCell)
				view.WriteString(rightSide)
			}
		}
		view.WriteByte('\n')
	}

	view.WriteByte('\n')

	if m.IsLocalPlayerTurn() {
		view.WriteString("Current player: ")
		view.WriteString(m.Theme.PlayerCells[m.CurrentPlayer])
	} else {
		view.WriteString("Awaiting player ")
		view.WriteString(m.Theme.PlayerCells[m.CurrentPlayer])
		view.WriteString(" move...")
	}

	// view.WriteString(fmt.Sprintf("\nCamera bound: %v | Camera: %v", m.cameraBound, m.Camera))
	view.WriteByte('\n')

	return view.String()
}
