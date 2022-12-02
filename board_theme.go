package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type BoardTheme struct {
	InvalidCell    string
	UnoccupiedCell string
	PlayerCells    []string

	PlayerCellStyles []lipgloss.Style

	CandidateCellStyle lipgloss.Style
	VictoryCellStyle   lipgloss.Style
	LastEnemyCellStyle lipgloss.Style

	SelectionInactiveStyle lipgloss.Style
}

// ApplyCyclingStyles will render each text using the corresponding style mod
// style array length. The input array is not changed
func ApplyCyclingStyles(text []string, styles []lipgloss.Style) []string {
	styledText := make([]string, 0, len(text))

	for i, str := range text {
		styled := styles[i%len(styles)].Render(str)
		styledText = append(styledText, styled)
	}

	return styledText
}

func (ts *BoardTheme) BoardToText(board map[Offset]PlayerID, camera Rect) map[Offset]string {
	cliBoard := make(map[Offset]string, camera.Area())

	for x := 0; x < camera.W; x++ {
		for y := 0; y < camera.H; y++ {
			curCell := camera.ToWorldXY(x, y)

			player, present := board[curCell]
			if !present {
				cliBoard[curCell] = ts.InvalidCell
				continue
			}

			if player == 0 {
				cliBoard[curCell] = ts.UnoccupiedCell
				continue
			}

			if int(player)-1 >= len(ts.PlayerCells) {
				panic(fmt.Sprintf("board theme: no player style definition for ID=%v: out of range (PlayerCount=%v)", player, len(ts.PlayerCells)))
			}

			cliBoard[curCell] = ts.PlayerCells[player-1]
		}
	}
	return cliBoard
}
