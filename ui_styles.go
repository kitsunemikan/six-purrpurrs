package main

import (
	"github.com/charmbracelet/lipgloss"
)

var defaultBoardTheme = BoardTheme{
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
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color("1")),

	SelectionInactiveStyle: lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("8")),
}
