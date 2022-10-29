package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var winCellStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("10"))

var inactiveTextStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("8"))

var candidateStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("177"))

var (
	playersFlag = flag.String("players", "XO", "a list of players and their tokens")
	wFlag       = flag.Uint("w", 3, "board width")
	hFlag       = flag.Uint("h", 3, "board height")
	strikeFlag  = flag.Uint("strike", 3, "the number of marks in a row to win the game")
)

type Model struct {
	Game *GameState

	Selection     Offset
	CurrentPlayer PlayerID
}

func initialModel(game *GameState) Model {
	return Model{
		Game: game,

		Selection:     Offset{game.BoardSize().X / 2, game.BoardSize().Y / 2},
		CurrentPlayer: 1,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

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
			if m.Game.Over() {
				return m, tea.Quit
			}

			if m.Game.Cell(m.Selection) != Unoccupied {
				return m, nil
			}

			m.Game.MarkCell(m.Selection, m.CurrentPlayer)

			if m.Game.Over() {
				// TODO: return GameOverModel
				return m, nil
			}

			if m.CurrentPlayer == m.Game.LastPlayer() {
				m.CurrentPlayer = 1
			} else {
				m.CurrentPlayer++
			}
		}
	}

	return m, nil
}

func (m *Model) solutionView(cliBoard map[Offset]string) string {
	for _, cell := range m.Game.Solution() {
		cliBoard[cell] = winCellStyle.Render(cliBoard[cell])
	}

	var view strings.Builder
	for y := 0; y < m.Game.BoardSize().Y; y++ {
		for x := 0; x < m.Game.BoardSize().X; x++ {
			view.WriteRune(' ')

			curCell := Offset{x, y}
			if m.Game.Cell(curCell) == Unoccupied {
				view.WriteString(".")
			} else {
				view.WriteString(cliBoard[Offset{x, y}])
			}

			view.WriteRune(' ')
		}
		view.WriteByte('\n')
	}

	view.WriteByte('\n')

	if m.Game.Solution() == nil {
		view.WriteString("A draw...")
	} else {
		view.WriteRune(m.Game.PlayerToken(m.CurrentPlayer))
		view.WriteString(" wins!")
	}

	view.WriteByte('\n')

	return view.String()
}

func (m *Model) selectionView(cliBoard map[Offset]string) string {
	if m.Game.Cell(m.Selection) == Unoccupied {
		candidates := m.Game.CandidateCellsAt(m.Selection, m.CurrentPlayer)

		for _, cell := range candidates {
			cliBoard[cell] = candidateStyle.Render(cliBoard[cell])
		}
	}

	var view strings.Builder
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
				view.WriteString(".")
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

	view.WriteString("Current player: ")
	view.WriteRune(m.Game.PlayerToken(m.CurrentPlayer))

	view.WriteByte('\n')

	return view.String()
}

func (m Model) View() string {
	cliBoard := m.Game.BoardToStrings()

	if m.Game.Over() {
		return m.solutionView(cliBoard)
	}

	return m.selectionView(cliBoard)
}

func main() {
	flag.Parse()

	w, h := int(*wFlag), int(*hFlag)
	conf := GameOptions{
		BoardSize:    Offset{w, h},
		StrikeLength: int(*strikeFlag),
		PlayerTokens: []rune(*playersFlag),
	}

	game := NewGame(conf)

	p := tea.NewProgram(initialModel(game))
	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "internal error: %v\n", err)
		os.Exit(1)
	}
}
