package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
    Counter int
}

var initialModel = model{
    Counter: 0,
}

func (m model) Init() tea.Cmd {
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            return m, tea.Quit

        case "up", "k":
            m.Counter += 1
            return m, nil

        case "down", "j":
            m.Counter -= 1
            return m, nil
        }
    }

    return m, nil
}

func (m model) View() string {
    return fmt.Sprintf("Counter: %d\n", m.Counter)
}

func main() {
    p := tea.NewProgram(initialModel)
    if err := p.Start(); err != nil {
        fmt.Fprintf(os.Stderr, "internal error: %v\n", err)
        os.Exit(1)
    }
}
