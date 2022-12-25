package game_test

import (
	"testing"

	"github.com/kitsunemikan/six-purrpurrs/game"
	"github.com/kitsunemikan/six-purrpurrs/geom"

	"github.com/maxatome/go-testdeep/td"
)

// TODO: test
func StrikeFromStr(start geom.Offset, dir game.StrikeDir, desc string) (strike game.Strike) {
	if desc == "" {
		panic("strike from str: empty description")
	}

	strike.Start = start
	strike.Dir = dir

	i := 0
	if desc[0] == '.' {
		strike.ExtendableBefore = true
		i = 1
	}

	if desc[i] == 'X' {
		strike.Player = game.P1
	} else if desc[i] == 'O' {
		strike.Player = game.P2
	} else {
		panic("strike from str: unknown player avatar, should be either X or O")
	}

	for i < len(desc) && desc[i] != '.' {
		strike.Len++
		i++
	}

	if i < len(desc) {
		strike.ExtendableAfter = true
	}

	return
}

func TestStrikeSet(t *testing.T) {
	t.Run("empty strike set returns no strikes", func(t *testing.T) {
		set := &game.StrikeSet{}

		got := set.Strikes()
		if got != nil {
			t.Errorf("got %v, wanted nil", got)
		}
	})

	t.Run("single cell results in 4 strikes", func(t *testing.T) {
		set := &game.StrikeSet{}
		set.MakeMove(game.PlayerMove{Cell: geom.Offset{X: 0, Y: 0}, ID: game.P1})

		got := set.Strikes()
		want := []game.Strike{
			StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightUp, ".X."),
			StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRight, ".X."),
			StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightDown, ".X."),
			StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeDown, ".X."),
		}

		td.Cmp(t, got, td.Bag(td.Flatten(want)))
	})

	t.Run("twe cells far apart result in 8 strikes", func(t *testing.T) {
		set := &game.StrikeSet{}
		set.MakeMove(game.PlayerMove{Cell: geom.Offset{X: 0, Y: 0}, ID: game.P2})
		set.MakeMove(game.PlayerMove{Cell: geom.Offset{X: 3, Y: 3}, ID: game.P1})

		got := set.Strikes()
		want := []game.Strike{
			StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightUp, ".O."),
			StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRight, ".O."),
			StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightDown, ".O."),
			StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeDown, ".O."),

			StrikeFromStr(geom.Offset{X: 3, Y: 3}, game.StrikeRightUp, ".X."),
			StrikeFromStr(geom.Offset{X: 3, Y: 3}, game.StrikeRight, ".X."),
			StrikeFromStr(geom.Offset{X: 3, Y: 3}, game.StrikeRightDown, ".X."),
			StrikeFromStr(geom.Offset{X: 3, Y: 3}, game.StrikeDown, ".X."),
		}

		td.Cmp(t, got, td.Bag(td.Flatten(want)))
	})
}
