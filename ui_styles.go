package main

import (
	"github.com/charmbracelet/lipgloss"
)

var winCellStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("0")).
	Background(lipgloss.Color("2"))

var inactiveTextStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("8"))

var candidateStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("177"))

var lastEnemyCellStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("0")).
	Background(lipgloss.Color("1"))
