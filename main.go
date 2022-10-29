package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var playerTypes = map[string]func(){
	"local": func() {},
	"ai":    func() {},
}

var (
	avatarsFlag     = flag.String("avatars", "X,O", "a list of strings used as player board marks")
	playerTypesFlag = flag.String("playertypes", "local,ai", fmt.Sprintf("specifies logic for each player (available options: %s)", availablePlayerTypes()))
	wFlag           = flag.Uint("w", 3, "board width")
	hFlag           = flag.Uint("h", 3, "board height")
	strikeFlag      = flag.Uint("strike", 3, "the number of marks in a row to win the game")
)

func availablePlayerTypes() (list string) {
	typeID := 0
	for name := range playerTypes {
		list += name
		if typeID < len(playerTypes)-1 {
			list += " "
		}
		typeID++
	}

	return
}

func main() {
	flag.Parse()

	avatars := strings.FieldsFunc(*avatarsFlag, func(c rune) bool { return c == ',' })

	w, h := int(*wFlag), int(*hFlag)
	conf := GameOptions{
		BoardSize:    Offset{w, h},
		StrikeLength: int(*strikeFlag),
		PlayerTokens: avatars,
	}

	game := NewGame(conf)

	p := tea.NewProgram(NewGameplayModel(game))
	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "internal error: %v\n", err)
		os.Exit(1)
	}
}
