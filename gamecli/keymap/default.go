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

	WatchReplay = key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "watch replay"),
	)
	Forward = key.NewBinding(
		key.WithKeys("shift+right", "f"),
		key.WithHelp("f/shift+→", "forward"),
	)
	Rewind = key.NewBinding(
		key.WithKeys("shift+left", "r"),
		key.WithHelp("r/shift+←", "rewind"),
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
	WatchReplay: WatchReplay,
	Quit:        QuitOrSelect,
}

var Replay = ReplayModel{
	Left:    Left,
	Right:   Right,
	Up:      Up,
	Down:    Down,
	Forward: Forward,
	Rewind:  Rewind,
	Help:    Help,
	Quit:    Quit,
}
