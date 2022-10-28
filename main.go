package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Offset struct {
	X, Y int
}

type CellState int

const (
	Unoccupied CellState = iota
	Player1
	Player2
)

type Model struct {
	BoardSize Offset
	Board     map[Offset]CellState

	Selection     Offset
	CurrentPlayer int

	PlayerCount int
}

var initialModel = Model{
	BoardSize: Offset{3, 3},
	Board:     make(map[Offset]CellState),

	Selection:     Offset{1, 1},
	CurrentPlayer: 0,

	PlayerCount: 2,
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

			switch m.Board[Offset{x, y}] {
			case Unoccupied:
				view.WriteByte('.')
			case Player1:
				view.WriteByte('X')
			case Player2:
				view.WriteByte('O')
			}

			view.WriteRune(rightSide)
		}
		view.WriteByte('\n')
	}

	return view.String()
}

func main() {
	p := tea.NewProgram(initialModel)
	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "internal error: %v\n", err)
		os.Exit(1)
	}
}
