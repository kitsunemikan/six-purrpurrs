package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

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

type Model struct {
	BoardSize Offset
	Board     map[Offset]PlayerID

	Selection     Offset
	CurrentPlayer PlayerID

	playerTokens []rune
}

func initialModel() Model {
	w, h := int(*wFlag), int(*hFlag)
	return Model{
		BoardSize: Offset{w, h},
		Board:     make(map[Offset]PlayerID),

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
			if m.Board[m.Selection] != Unoccupied {
				return m, nil
			}

			m.Board[m.Selection] = m.CurrentPlayer
			if m.CurrentPlayer == m.LastPlayer() {
				m.CurrentPlayer = 1
			} else {
				m.CurrentPlayer++
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
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

			if m.Board[Offset{x, y}] == 0 {
				view.WriteByte('.')
			} else {
				view.WriteRune(m.PlayerToken(m.Board[Offset{x, y}]))
			}

			view.WriteRune(rightSide)
		}
		view.WriteByte('\n')
	}

	view.WriteByte('\n')
	view.WriteString("Current player: ")
	view.WriteRune(m.PlayerToken(m.CurrentPlayer))
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
