package gamecli

import (
	"github.com/charmbracelet/bubbles/key"
)

type GameplayKeymap struct {
	Left   key.Binding
	Right  key.Binding
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Help   key.Binding
	Quit   key.Binding
}

func (k GameplayKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k GameplayKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Left, k.Right, k.Up, k.Down, k.Select},
		{k.Help, k.Quit},
	}
}

type GameOverKeymap struct {
	Quit key.Binding
}

func (k GameOverKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

func (k GameOverKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit},
	}
}

var GlobalGameplayKeymap = GameplayKeymap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("enter/space", "make move"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

var GlobalGameOverKeymap = GameOverKeymap{
	Quit: key.NewBinding(
		key.WithKeys("enter", " ", "q", "esc", "ctrl+c"),
		key.WithHelp("enter/q", "quit"),
	),
}
