package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/kitsunemikan/six-purrpurrs/ai"
	"github.com/kitsunemikan/six-purrpurrs/game"
	"github.com/kitsunemikan/six-purrpurrs/gamecli"
	. "github.com/kitsunemikan/six-purrpurrs/geom"
)

var playerTypeGenerators = map[string]func() game.PlayerAgent{
	"local":  gamecli.NewLocalPlayer,
	"random": ai.NewRandomPlayer,
}

var (
	avatarsFlag     = flag.String("avatars", " ,.,X,O", "a list of strings used as player board marks, including unavailable and empty cell")
	playerTypesFlag = flag.String("playertypes", "local,random", fmt.Sprintf("specifies logic for each player (available options: %s)", availablePlayerTypes()))
	wFlag           = flag.Uint("w", 40, "screen width")
	hFlag           = flag.Uint("h", 20, "screen height")
	borderFlag      = flag.Uint("border", 7, "the width of a border around marked cells where players can make a move")
	strikeFlag      = flag.Uint("strike", 6, "the number of marks in a row to win the game")
	trackDepthFlag  = flag.Uint("trackDepth", 20, "The width of camera borders in % after which to follow player moves")
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

func CommaListUnfiltered(s string) []string {
	return strings.FieldsFunc(s, func(c rune) bool { return c == ',' })
}

func CommaList(s string) []string {
	fields := CommaListUnfiltered(s)
	for i, dirty := range fields {
		fields[i] = strings.TrimSpace(dirty)
	}

	return fields
}

func main() {
	flag.Parse()

	avatars := CommaListUnfiltered(*avatarsFlag)
	if len(avatars) < 3 {
		fmt.Fprintf(os.Stderr, "error: at least 3 avatars must be provided (for unavailable, unoccupied and first player cells). Provided %d\n", len(avatars))
		os.Exit(1)
	}

	theme := gamecli.DefaultBoardTheme
	theme.InvalidCell = avatars[0]
	theme.UnoccupiedCell = avatars[1]
	theme.PlayerCells = avatars[2:]

	// Create players
	playerTypes := CommaList(*playerTypesFlag)
	if len(theme.PlayerCells) != len(playerTypes) {
		fmt.Fprintf(os.Stderr, "error: mismatch between number of player avatars (%d) and number of player types (%d) specified\nnote: avatars: %s\nnote: player types: %s\n", len(theme.PlayerCells), len(playerTypes), *avatarsFlag, *playerTypesFlag)
		os.Exit(1)
	}

	players := make([]game.PlayerAgent, 0, len(playerTypes))
	for _, playerType := range playerTypes {
		if _, exists := playerTypeGenerators[playerType]; !exists {
			fmt.Fprintf(os.Stderr, "error: invalid player type supplied: '%s'\nnote: available types are: %s\n", playerType, availablePlayerTypes())
			os.Exit(1)
		}

		players = append(players, playerTypeGenerators[playerType]())
	}

	gameConf := game.GameOptions{
		Border:       int(*borderFlag),
		StrikeLength: int(*strikeFlag),
	}

	game := game.NewGame(gameConf)

	w, h := int(*wFlag), int(*hFlag)
	minDim := w
	if h < minDim {
		minDim = h
	}

	trackDepth := int(*trackDepthFlag) * minDim / 100
	modelConf := gamecli.GameplayModelConfig{
		Game:       game,
		Players:    players,
		Theme:      &theme,
		ScreenSize: Offset{X: w, Y: h},
		TrackDepth: trackDepth,
	}

	p := tea.NewProgram(gamecli.NewGameplayModel(modelConf))
	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "internal error: %v\n", err)
		os.Exit(1)
	}
}
