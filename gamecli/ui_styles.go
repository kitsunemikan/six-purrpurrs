package gamecli

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

var DefaultBoardTheme = BoardTheme{
	PlayerCellStyles: []lipgloss.Style{
		// SlateBlue1
		lipgloss.NewStyle().Foreground(lipgloss.Color("99")),
		// Orange3
		lipgloss.NewStyle().Foreground(lipgloss.Color("172")),
	},

	CandidateCellStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("177")),

	VictoryCellStyle: lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("10")),

	LastEnemyCellStyle: lipgloss.NewStyle().
		Background(lipgloss.Color("88")),

	SelectionInactiveStyle: lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("8")),
}

var helpKeyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#8a8a8a"))

var helpDescStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))

var helpSepStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#4A4A4A"))

var HelpStyle = help.Styles{
	ShortKey:       helpKeyStyle,
	ShortDesc:      helpDescStyle,
	ShortSeparator: helpSepStyle,
	Ellipsis:       helpSepStyle.Copy(),
	FullKey:        helpKeyStyle.Copy(),
	FullDesc:       helpDescStyle.Copy(),
	FullSeparator:  helpSepStyle.Copy(),
}
