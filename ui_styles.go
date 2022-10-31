package main

import (
	"github.com/charmbracelet/lipgloss"
)

var winCellStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("10"))

var inactiveTextStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("8"))

var candidateStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("177"))

var enemyCellStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("210"))
