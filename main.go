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

type Offset struct {
	X, Y int
}

func (a Offset) Add(b Offset) Offset {
	return Offset{a.X + b.X, a.Y + b.Y}
}

func (a Offset) Scale(c int) Offset {
	return Offset{c * a.X, c * a.Y}
}

type PlayerID int

const Unoccupied PlayerID = 0

var (
	playersFlag = flag.String("players", "XO", "a list of players and their tokens")
	wFlag       = flag.Uint("w", 3, "board width")
	hFlag       = flag.Uint("h", 3, "board height")
	strikeFlag  = flag.Uint("strike", 3, "the number of marks in a row to win the game")
)

var solutionOffsets = []Offset{{1, 0}, {1, 1}, {0, 1}, {1, -1}}

type GameOptions struct {
	BoardSize    Offset
	StrikeLength int
	PlayerTokens []rune
}

type Model struct {
	Conf GameOptions

	Board    map[Offset]PlayerID
	Solution []Offset

	Selection     Offset
	CurrentPlayer PlayerID
}

func initialModel() Model {
	w, h := int(*wFlag), int(*hFlag)
	return Model{
		Conf: GameOptions{
			StrikeLength: int(*strikeFlag),
			BoardSize:    Offset{w, h},
			PlayerTokens: []rune(*playersFlag),
		},

		Board: make(map[Offset]PlayerID),

		Selection:     Offset{w / 2, h / 2},
		CurrentPlayer: 1,
	}
}

func (m *Model) LastPlayer() PlayerID {
	return PlayerID(len(m.Conf.PlayerTokens))
}

func (m *Model) PlayerToken(player PlayerID) rune {
	if player < 1 || player > m.LastPlayer() {
		panic(fmt.Sprintf("model: player token for ID=%v: out of range (LastPlayerID=%v)", player, m.LastPlayer()))
	}

	return m.Conf.PlayerTokens[int(player)-1]
}

func (m *Model) IsInsideBoard(pos Offset) bool {
	return pos.X >= 0 && pos.X < m.Conf.BoardSize.X && pos.Y >= 0 && pos.Y < m.Conf.BoardSize.Y
}

func (m *Model) CheckSolutionsAt(pos Offset, player PlayerID) []Offset {
	solution := make([]Offset, 0, m.Conf.StrikeLength)

	for _, dir := range solutionOffsets {
		solution = solution[:0]

		if len(solution) != 0 {
			panic(42)
		}

		for i := 0; i < m.Conf.StrikeLength; i++ {
			curCell := Offset{pos.X + i*dir.X, pos.Y + i*dir.Y}
			if m.Board[curCell] != player {
				break
			}

			solution = append(solution, curCell)
		}

		for i := 1; i < m.Conf.StrikeLength; i++ {
			curCell := Offset{pos.X - i*dir.X, pos.Y - i*dir.Y}
			if m.Board[curCell] != player {
				break
			}

			solution = append(solution, curCell)
		}

		if len(solution) >= m.Conf.StrikeLength {
			return solution
		}
	}

	return nil
}

func (m *Model) CandidateCellsAt(pos Offset, player PlayerID) []Offset {
	candidates := make([]Offset, 0, 2*len(solutionOffsets)*(m.Conf.StrikeLength-1)+1)

	if m.Board[pos] == player {
		candidates = append(candidates, pos)
	}

	for _, dir := range solutionOffsets {
		for i := 1; i < m.Conf.StrikeLength; i++ {
			curCell := pos.Add(dir.Scale(i))
			if m.Board[curCell] == player {
				candidates = append(candidates, curCell)
			} else {
				break
			}
		}

		for i := 1; i < m.Conf.StrikeLength; i++ {
			curCell := pos.Add(dir.Scale(-i))
			if m.Board[curCell] == player {
				candidates = append(candidates, curCell)
			} else {
				break
			}
		}
	}

	return candidates
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
			if m.Selection.X < m.Conf.BoardSize.X-1 {
				m.Selection.X += 1
			}
			return m, nil

		case "up", "k":
			if m.Selection.Y > 0 {
				m.Selection.Y -= 1
			}
			return m, nil

		case "down", "j":
			if m.Selection.Y < m.Conf.BoardSize.Y-1 {
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
		cliBoard[cell] = winCellStyle.Render(cliBoard[cell])
	}

	var view strings.Builder
	for y := 0; y < m.Conf.BoardSize.Y; y++ {
		for x := 0; x < m.Conf.BoardSize.X; x++ {
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
	if m.Board[m.Selection] == Unoccupied {
		candidates := m.CandidateCellsAt(m.Selection, m.CurrentPlayer)

		for _, cell := range candidates {
			cliBoard[cell] = candidateStyle.Render(cliBoard[cell])
		}
	}

	var view strings.Builder
	for y := 0; y < m.Conf.BoardSize.Y; y++ {
		for x := 0; x < m.Conf.BoardSize.X; x++ {
			leftSide := " "
			rightSide := " "
			if x == m.Selection.X && y == m.Selection.Y {
				leftSide = "["
				rightSide = "]"
			}

			curCell := Offset{x, y}

			if m.Board[curCell] != 0 {
				view.WriteString(inactiveTextStyle.Render(leftSide))
			} else {
				view.WriteString(leftSide)
			}

			if m.Board[curCell] == 0 {
				view.WriteString(".")
			} else {
				view.WriteString(cliBoard[Offset{x, y}])
			}

			if m.Board[curCell] != 0 {
				view.WriteString(inactiveTextStyle.Render(rightSide))
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
}

func main() {
	flag.Parse()

	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "internal error: %v\n", err)
		os.Exit(1)
	}
}
