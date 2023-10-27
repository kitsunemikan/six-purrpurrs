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
		key.WithKeys("right", "f"),
		key.WithHelp("f/→", "forward"),
	)
	Rewind = key.NewBinding(
		key.WithKeys("left", "r"),
		key.WithHelp("r/←", "rewind"),
	)

	ReplayCameraUp = key.NewBinding(
		key.WithKeys("shift+up", "k"),
		key.WithHelp("k/shift+↑", "forward"),
	)
	ReplayCameraDown = key.NewBinding(
		key.WithKeys("shift+down", "j"),
		key.WithHelp("j/shift+↓", "rewind"),
	)
	ReplayCameraRight = key.NewBinding(
		key.WithKeys("shift+right", "l"),
		key.WithHelp("l/shift+→", "forward"),
	)
	ReplayCameraLeft = key.NewBinding(
		key.WithKeys("shift+left", "h"),
		key.WithHelp("h/shift+←", "rewind"),
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
	Left:    ReplayCameraLeft,
	Right:   ReplayCameraRight,
	Up:      ReplayCameraUp,
	Down:    ReplayCameraDown,
	Forward: Forward,
	Rewind:  Rewind,
	Help:    Help,
	Quit:    Quit,
}
