package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var winCell = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("10"))

var inactiveText = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("8"))

type Offset struct {
	X, Y int
}

type PlayerID int

const Unoccupied PlayerID = 0

var (
	playersFlag = flag.String("players", "XO", "a list of players and their tokens")
	wFlag       = flag.Uint("w", 3, "board width")
	hFlag       = flag.Uint("h", 3, "board height")
)

var solutionOffsets = []Offset{{1, 0}, {1, 1}, {0, 1}, {1, -1}}

type Model struct {
	StrikeLength int
	BoardSize    Offset
	Board        map[Offset]PlayerID
	Solution     []Offset

	Selection     Offset
	CurrentPlayer PlayerID

	playerTokens []rune
}

func initialModel() Model {
	w, h := int(*wFlag), int(*hFlag)
	return Model{
		StrikeLength: 3,
		BoardSize:    Offset{w, h},
		Board:        make(map[Offset]PlayerID),

		Selection:     Offset{w / 2, h / 2},
		CurrentPlayer: 1,

		playerTokens: []rune(*playersFlag),
	}
}

func (m *Model) LastPlayer() PlayerID {
	return PlayerID(len(m.playerTokens))
}

func (m *Model) PlayerToken(player PlayerID) rune {
	if player < 1 || player > m.LastPlayer() {
		panic(fmt.Sprintf("model: player token for ID=%v: out of range (LastPlayerID=%v)", player, m.LastPlayer()))
	}

	return m.playerTokens[int(player)-1]
}

func (m *Model) IsInsideBoard(pos Offset) bool {
	return pos.X >= 0 && pos.X < m.BoardSize.X && pos.Y >= 0 && pos.Y < m.BoardSize.Y
}

func (m *Model) CheckSolutionsAt(pos Offset, player PlayerID) []Offset {
	solution := make([]Offset, 0, m.StrikeLength)

	for _, dir := range solutionOffsets {
		solution = solution[:0]

		if len(solution) != 0 {
			panic(42)
		}

		for i := 0; i < m.StrikeLength; i++ {
			curCell := Offset{pos.X + i*dir.X, pos.Y + i*dir.Y}
			if m.Board[curCell] != player {
				break
			}

			solution = append(solution, curCell)
		}

		for i := 1; i < m.StrikeLength; i++ {
			curCell := Offset{pos.X - i*dir.X, pos.Y - i*dir.Y}
			if m.Board[curCell] != player {
				break
			}

			solution = append(solution, curCell)
		}

		if len(solution) >= m.StrikeLength {
			return solution
		}
	}

	return nil
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
			if m.Selection.X < m.BoardSize.X-1 {
				m.Selection.X += 1
			}
			return m, nil

		case "up", "k":
			if m.Selection.Y > 0 {
				m.Selection.Y -= 1
			}
			return m, nil

		case "down", "j":
			if m.Selection.Y < m.BoardSize.Y-1 {
				m.Selection.Y += 1
			}
			return m, nil

		case "enter", " ":
			if m.Solution != nil {
				return m, tea.Quit
			}

			if m.Board[m.Selection] != Unoccupied {
				return m, nil
			}

			m.Board[m.Selection] = m.CurrentPlayer

			m.Solution = m.CheckSolutionsAt(m.Selection, m.CurrentPlayer)
			if m.Solution != nil {
				return m, nil
			}

			if m.CurrentPlayer == m.LastPlayer() {
				m.CurrentPlayer = 1
			} else {
				m.CurrentPlayer++
			}
		}
	}

	return m, nil
}

func (m *Model) solutionView(cliBoard map[Offset]string) string {
	for _, cell := range m.Solution {
		cliBoard[cell] = winCell.Render(cliBoard[cell])
	}

	var view strings.Builder
	for y := 0; y < m.BoardSize.Y; y++ {
		for x := 0; x < m.BoardSize.X; x++ {
			view.WriteRune(' ')

			curCell := Offset{x, y}
			if m.Board[curCell] == 0 {
				view.WriteString(".")
			} else {
				view.WriteString(cliBoard[Offset{x, y}])
			}

			view.WriteRune(' ')
		}
		view.WriteByte('\n')
	}

	view.WriteByte('\n')

	view.WriteRune(m.PlayerToken(m.CurrentPlayer))
	view.WriteString(" wins!")

	view.WriteByte('\n')

	return view.String()
}

func (m *Model) selectionView(cliBoard map[Offset]string) string {
	var view strings.Builder
	for y := 0; y < m.BoardSize.Y; y++ {
		for x := 0; x < m.BoardSize.X; x++ {
			leftSide := " "
			rightSide := " "
			if x == m.Selection.X && y == m.Selection.Y {
				leftSide = "["
				rightSide = "]"
			}

			curCell := Offset{x, y}

			if m.Board[curCell] != 0 {
				view.WriteString(inactiveText.Render(leftSide))
			} else {
				view.WriteString(leftSide)
			}

			if m.Board[curCell] == 0 {
				view.WriteString(".")
			} else {
				view.WriteString(cliBoard[Offset{x, y}])
			}

			if m.Board[curCell] != 0 {
				view.WriteString(inactiveText.Render(rightSide))
			} else {
				view.WriteString(rightSide)
			}
		}
		view.WriteByte('\n')
	}

	view.WriteByte('\n')

	view.WriteString("Current player: ")
	view.WriteRune(m.PlayerToken(m.CurrentPlayer))

	view.WriteByte('\n')

	return view.String()
}

func (m Model) View() string {
	cliBoard := make(map[Offset]string)
	for pos, player := range m.Board {
		cliBoard[pos] = string(m.PlayerToken(player))
	}

	if m.Solution != nil {
		return m.solutionView(cliBoard)
	}

	return m.selectionView(cliBoard)

	var view strings.Builder
	for y := 0; y < m.BoardSize.Y; y++ {
		for x := 0; x < m.BoardSize.X; x++ {
			leftSide := ' '
			rightSide := ' '
			if x == m.Selection.X && y == m.Selection.Y {
				leftSide = '['
				rightSide = ']'
			}
			view.WriteRune(leftSide)

			curCell := Offset{x, y}
			if m.Board[curCell] == 0 {
				view.WriteString(".")
			} else {
				view.WriteString(cliBoard[Offset{x, y}])
			}

			view.WriteRune(rightSide)
		}
		view.WriteByte('\n')
	}

	view.WriteByte('\n')

	if m.Solution != nil {
		view.WriteRune(m.PlayerToken(m.CurrentPlayer))
		view.WriteString(" wins!")
	} else {
		view.WriteString("Current player: ")
		view.WriteRune(m.PlayerToken(m.CurrentPlayer))
	}

	view.WriteByte('\n')

	return view.String()
}

func main() {
	flag.Parse()

	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "internal error: %v\n", err)
		os.Exit(1)
	}
}
