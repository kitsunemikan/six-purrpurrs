package keymap

import "github.com/charmbracelet/bubbles/key"

var (
	Left = key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	)
	Right = key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	)
	Up = key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	)
	Down = key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	)

	Select = key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("enter/space", "make move"),
	)
	Help = key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	)
	Quit = key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	)
	QuitOrSelect = key.NewBinding(
		key.WithKeys("enter", " ", "q", "esc", "ctrl+c"),
		key.WithHelp("enter/q", "quit"),
	)
)

var Gameplay = GameplayModel{
	Left:   Left,
	Right:  Right,
	Up:     Up,
	Down:   Down,
	Select: Select,
	Help:   Help,
	Quit:   Quit,
}

var GameOver = GameOverModel{
	Quit: QuitOrSelect,
}
