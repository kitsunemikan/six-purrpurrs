package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	playersFlag = flag.String("players", "XO", "a list of players and their tokens")
	wFlag       = flag.Uint("w", 3, "board width")
	hFlag       = flag.Uint("h", 3, "board height")
	strikeFlag  = flag.Uint("strike", 3, "the number of marks in a row to win the game")
)

func main() {
	flag.Parse()

	w, h := int(*wFlag), int(*hFlag)
	conf := GameOptions{
		BoardSize:    Offset{w, h},
		StrikeLength: int(*strikeFlag),
		PlayerTokens: []rune(*playersFlag),
	}

	game := NewGame(conf)

	p := tea.NewProgram(NewGameplayModel(game))
	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "internal error: %v\n", err)
		os.Exit(1)
	}
}
