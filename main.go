package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var playerTypeGenerators = map[string]func() PlayerAgent{
	"local":  NewLocalPlayer,
	"random": NewRandomPlayer,
}

var (
	avatarsFlag     = flag.String("avatars", ".,X,O", "a list of strings used as player board marks, including empty cell")
	playerTypesFlag = flag.String("playertypes", "local,random", fmt.Sprintf("specifies logic for each player (available options: %s)", availablePlayerTypes()))
	wFlag           = flag.Uint("w", 3, "board width")
	hFlag           = flag.Uint("h", 3, "board height")
	strikeFlag      = flag.Uint("strike", 3, "the number of marks in a row to win the game")
)

func availablePlayerTypes() (list string) {
	typeID := 0
	for name := range playerTypeGenerators {
		list += name
		if typeID < len(playerTypeGenerators)-1 {
			list += ", "
		}
		typeID++
	}

	return
}

func CommaList(s string) []string {
	fields := strings.FieldsFunc(s, func(c rune) bool { return c == ',' })
	for i, dirty := range fields {
		fields[i] = strings.TrimSpace(dirty)
	}

	return fields
}

func main() {
	flag.Parse()

	// Initialize game
	avatars := CommaList(*avatarsFlag)

	w, h := int(*wFlag), int(*hFlag)
	conf := GameOptions{
		BoardSize:    Offset{w, h},
		StrikeLength: int(*strikeFlag),
		PlayerTokens: avatars,
	}

	game := NewGame(conf)

	// Create players
	playerTypes := CommaList(*playerTypesFlag)

	players := map[PlayerID]PlayerAgent{}
	for i, playerType := range playerTypes {
		if _, exists := playerTypeGenerators[playerType]; !exists {
			fmt.Fprintf(os.Stderr, "error: invalid player type supplied: '%s'\nnote: available types are: %s\n", playerType, availablePlayerTypes())
			os.Exit(1)
		}

		players[PlayerID(i+1)] = playerTypeGenerators[playerType]()
	}

	p := tea.NewProgram(NewGameplayModel(game, players))
	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "internal error: %v\n", err)
		os.Exit(1)
	}
}
