package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/kitsunemikan/six-purrpurrs/ai"
	"github.com/kitsunemikan/six-purrpurrs/game"
	"github.com/kitsunemikan/six-purrpurrs/gamecli"
	. "github.com/kitsunemikan/six-purrpurrs/geom"
)

func NewAIPlayer(id game.PlayerID) game.PlayerAgent {
	return ai.NewDefaultAIPlayer(id)
}

func NewLocalPlayer(id game.PlayerID) game.PlayerAgent {
	return gamecli.NewLocalPlayer()
}

func NewRandomPlayer(id game.PlayerID) game.PlayerAgent {
	return ai.NewRandomPlayer()
}

func NewObstructingPlayer(id game.PlayerID) game.PlayerAgent {
	return ai.NewObstructingPlayer(id)
}

var playerTypeGenerators = map[string]func(game.PlayerID) game.PlayerAgent{
	"local":       NewLocalPlayer,
	"random":      NewRandomPlayer,
	"ai":          NewAIPlayer,
	"obstructing": NewObstructingPlayer,
}

var (
	unavailableCellFlag = flag.String("unavailablecell", " ", "a character to denote a yet locked cell")
	availableCellFlag   = flag.String("availablecell", ".", "a character to denote a cell available for a move")
	p1CellFlag          = flag.String("p1avatar", "X", "a character to denote the first player on the board")
	p2CellFlag          = flag.String("p2avatar", "O", "a character to denote the second player on the board")
	p1TypeFlag          = flag.String("p1", "local", fmt.Sprintf("specifies logic for the first player (available: %s)", availablePlayerTypes()))
	p2TypeFlag          = flag.String("p2", "random", fmt.Sprintf("specifies logic for the second player (available: %s)", availablePlayerTypes()))
	wFlag               = flag.Uint("w", 40, "screen width")
	hFlag               = flag.Uint("h", 20, "screen height")
	borderFlag          = flag.Uint("border", 7, "the width of a border around marked cells where players can make a move")
	strikeFlag          = flag.Uint("strike", 6, "the number of marks in a row to win the game")
	trackDepthFlag      = flag.Uint("trackDepth", 20, "The width of camera borders in % after which to follow player moves")
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

func main() {
	flag.Parse()

	theme := gamecli.DefaultBoardTheme
	theme.InvalidCell = *unavailableCellFlag
	theme.UnoccupiedCell = *availableCellFlag
	theme.PlayerCells = []string{*p1CellFlag, *p2CellFlag}

	// Create players
	players := make([]game.PlayerAgent, 2)
	if _, exists := playerTypeGenerators[*p1TypeFlag]; !exists {
		fmt.Fprintf(os.Stderr, "error: invalid player type supplied: '%s'\nnote: available types are: %s\n", *p1TypeFlag, availablePlayerTypes())
		os.Exit(1)
	}

	if _, exists := playerTypeGenerators[*p2TypeFlag]; !exists {
		fmt.Fprintf(os.Stderr, "error: invalid player type supplied: '%s'\nnote: available types are: %s\n", *p2TypeFlag, availablePlayerTypes())
		os.Exit(1)
	}

	players[0] = playerTypeGenerators[*p1TypeFlag](game.P1)
	players[1] = playerTypeGenerators[*p2TypeFlag](game.P2)

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
