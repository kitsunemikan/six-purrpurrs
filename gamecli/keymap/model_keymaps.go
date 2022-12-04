package keymap

import "github.com/charmbracelet/bubbles/key"

type GameplayModel struct {
	Left   key.Binding
	Right  key.Binding
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Help   key.Binding
	Quit   key.Binding
}

func (k GameplayModel) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k GameplayModel) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Left, k.Right, k.Up, k.Down, k.Select},
		{k.Help, k.Quit},
	}
}

type GameOverModel struct {
	Quit key.Binding
}

func (k GameOverModel) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

func (k GameOverModel) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit},
	}
}
