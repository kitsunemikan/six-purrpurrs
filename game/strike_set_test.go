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

func TestStrikeSetMakeMove(t *testing.T) {
	t.Run("empty strike set returns nil strike slice", func(t *testing.T) {
		set := game.NewStrikeSet()

		got := set.Strikes()
		if got != nil {
			t.Errorf("got %v, wanted nil", got)
		}
	})

	tests := []struct {
		description string
		moves       []game.PlayerMove
		want        []game.Strike
	}{
		{
			"single cell results in 4 strikes",
			[]game.PlayerMove{
				{Cell: geom.Offset{X: 0, Y: 0}, Player: game.P1},
			},
			[]game.Strike{
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRight, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightDown, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeDown, ".X."),
			},
		},
		{
			"two cells far apart result in 8 strikes",
			[]game.PlayerMove{
				{Cell: geom.Offset{X: 0, Y: 0}, Player: game.P2},
				{Cell: geom.Offset{X: 3, Y: 3}, Player: game.P1},
			},
			[]game.Strike{
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightUp, ".O."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRight, ".O."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightDown, ".O."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeDown, ".O."),

				StrikeFromStr(geom.Offset{X: 3, Y: 3}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 3, Y: 3}, game.StrikeRight, ".X."),
				StrikeFromStr(geom.Offset{X: 3, Y: 3}, game.StrikeRightDown, ".X."),
				StrikeFromStr(geom.Offset{X: 3, Y: 3}, game.StrikeDown, ".X."),
			},
		},
		{
			"two P1 cells making a 2-len strike going right",
			[]game.PlayerMove{
				{Cell: geom.Offset{X: 0, Y: 0}, Player: game.P1},
				{Cell: geom.Offset{X: 1, Y: 0}, Player: game.P1},
			},
			[]game.Strike{
				// Singular strikes for the first move
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightDown, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeDown, ".X."),

				// Singluar strikes for the second move
				StrikeFromStr(geom.Offset{X: 1, Y: 0}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 1, Y: 0}, game.StrikeRightDown, ".X."),
				StrikeFromStr(geom.Offset{X: 1, Y: 0}, game.StrikeDown, ".X."),

				// Mutual 2-len strike
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRight, ".XX."),
			},
		},
		{
			"two P1 cells making a 2-len strike going up",
			[]game.PlayerMove{
				{Cell: geom.Offset{X: 0, Y: 0}, Player: game.P1},
				{Cell: geom.Offset{X: 0, Y: -1}, Player: game.P1},
			},
			[]game.Strike{
				// Singular strikes for the first move
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRight, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightDown, ".X."),

				// Singluar strikes for the second move
				StrikeFromStr(geom.Offset{X: 0, Y: -1}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: -1}, game.StrikeRight, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: -1}, game.StrikeRightDown, ".X."),

				// Mutual 2-len strike
				StrikeFromStr(geom.Offset{X: 0, Y: -1}, game.StrikeDown, ".XX."),
			},
		},
		{
			"three P1 cells making a 3-len strike by joining two single cells in right down direction",
			[]game.PlayerMove{
				{Cell: geom.Offset{X: 0, Y: 0}, Player: game.P1},
				{Cell: geom.Offset{X: 2, Y: 2}, Player: game.P1},
				{Cell: geom.Offset{X: 1, Y: 1}, Player: game.P1},
			},
			[]game.Strike{
				// Singular strikes for the first move
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRight, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeDown, ".X."),

				// Singular strikes for the second move
				StrikeFromStr(geom.Offset{X: 2, Y: 2}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 2, Y: 2}, game.StrikeRight, ".X."),
				StrikeFromStr(geom.Offset{X: 2, Y: 2}, game.StrikeDown, ".X."),

				// Singular strikes for the third move
				StrikeFromStr(geom.Offset{X: 1, Y: 1}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 1, Y: 1}, game.StrikeRight, ".X."),
				StrikeFromStr(geom.Offset{X: 1, Y: 1}, game.StrikeDown, ".X."),

				// Mutual 3-len strike
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightDown, ".XXX."),
			},
		},
		{
			"a 3x3 square of P2 cells",
			[]game.PlayerMove{
				// O.O
				// ...
				// O.O
				{Cell: geom.Offset{X: 0, Y: 0}, Player: game.P2},
				{Cell: geom.Offset{X: 2, Y: 0}, Player: game.P2},
				{Cell: geom.Offset{X: 2, Y: 2}, Player: game.P2},
				{Cell: geom.Offset{X: 0, Y: 2}, Player: game.P2},

				// OOO
				// O.O
				// OOO
				{Cell: geom.Offset{X: 1, Y: 0}, Player: game.P2},
				{Cell: geom.Offset{X: 2, Y: 1}, Player: game.P2},
				{Cell: geom.Offset{X: 1, Y: 2}, Player: game.P2},
				{Cell: geom.Offset{X: 0, Y: 1}, Player: game.P2},

				// OOO
				// OOO
				// OOO
				{Cell: geom.Offset{X: 1, Y: 1}, Player: game.P2},
			},
			[]game.Strike{
				// Vertical strikes
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeDown, ".OOO."),
				StrikeFromStr(geom.Offset{X: 1, Y: 0}, game.StrikeDown, ".OOO."),
				StrikeFromStr(geom.Offset{X: 2, Y: 0}, game.StrikeDown, ".OOO."),

				// Horizontal strikes
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRight, ".OOO."),
				StrikeFromStr(geom.Offset{X: 0, Y: 1}, game.StrikeRight, ".OOO."),
				StrikeFromStr(geom.Offset{X: 0, Y: 2}, game.StrikeRight, ".OOO."),

				// Right up strikes
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightUp, ".O."),
				StrikeFromStr(geom.Offset{X: 0, Y: 1}, game.StrikeRightUp, ".OO."),
				StrikeFromStr(geom.Offset{X: 0, Y: 2}, game.StrikeRightUp, ".OOO."),
				StrikeFromStr(geom.Offset{X: 1, Y: 2}, game.StrikeRightUp, ".OO."),
				StrikeFromStr(geom.Offset{X: 2, Y: 2}, game.StrikeRightUp, ".O."),

				// Right down strikes
				StrikeFromStr(geom.Offset{X: 0, Y: 2}, game.StrikeRightDown, ".O."),
				StrikeFromStr(geom.Offset{X: 0, Y: 1}, game.StrikeRightDown, ".OO."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightDown, ".OOO."),
				StrikeFromStr(geom.Offset{X: 1, Y: 0}, game.StrikeRightDown, ".OO."),
				StrikeFromStr(geom.Offset{X: 2, Y: 0}, game.StrikeRightDown, ".O."),
			},
		},
		{
			"create 2 3-len strikes by combining two singular cells, and then combine those",
			[]game.PlayerMove{
				// X.X.X.X
				{Cell: geom.Offset{X: 0, Y: 0}, Player: game.P1},
				{Cell: geom.Offset{X: 2, Y: 0}, Player: game.P1},
				{Cell: geom.Offset{X: 4, Y: 0}, Player: game.P1},
				{Cell: geom.Offset{X: 6, Y: 0}, Player: game.P1},

				// xxx.xxx
				{Cell: geom.Offset{X: 1, Y: 0}, Player: game.P1},
				{Cell: geom.Offset{X: 5, Y: 0}, Player: game.P1},

				// XXXXXXX
				{Cell: geom.Offset{X: 3, Y: 0}, Player: game.P1},
			},
			[]game.Strike{
				// Vertical strikes
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeDown, ".X."),
				StrikeFromStr(geom.Offset{X: 1, Y: 0}, game.StrikeDown, ".X."),
				StrikeFromStr(geom.Offset{X: 2, Y: 0}, game.StrikeDown, ".X."),
				StrikeFromStr(geom.Offset{X: 3, Y: 0}, game.StrikeDown, ".X."),
				StrikeFromStr(geom.Offset{X: 4, Y: 0}, game.StrikeDown, ".X."),
				StrikeFromStr(geom.Offset{X: 5, Y: 0}, game.StrikeDown, ".X."),
				StrikeFromStr(geom.Offset{X: 6, Y: 0}, game.StrikeDown, ".X."),

				// Horizontal strikes
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRight, ".XXXXXXX."),

				// Right up strikes
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 1, Y: 0}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 2, Y: 0}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 3, Y: 0}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 4, Y: 0}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 5, Y: 0}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 6, Y: 0}, game.StrikeRightUp, ".X."),

				// Right down strikes
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightDown, ".X."),
				StrikeFromStr(geom.Offset{X: 1, Y: 0}, game.StrikeRightDown, ".X."),
				StrikeFromStr(geom.Offset{X: 2, Y: 0}, game.StrikeRightDown, ".X."),
				StrikeFromStr(geom.Offset{X: 3, Y: 0}, game.StrikeRightDown, ".X."),
				StrikeFromStr(geom.Offset{X: 4, Y: 0}, game.StrikeRightDown, ".X."),
				StrikeFromStr(geom.Offset{X: 5, Y: 0}, game.StrikeRightDown, ".X."),
				StrikeFromStr(geom.Offset{X: 6, Y: 0}, game.StrikeRightDown, ".X."),
			},
		},
		{
			"Singular P1 and P2 beside each other",
			[]game.PlayerMove{
				{Cell: geom.Offset{X: 0, Y: 0}, Player: game.P1},
				{Cell: geom.Offset{X: 1, Y: 0}, Player: game.P2},
			},
			[]game.Strike{
				// P1
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRight, ".X"),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightDown, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeDown, ".X."),

				// P2
				StrikeFromStr(geom.Offset{X: 1, Y: 0}, game.StrikeRightUp, ".O."),
				StrikeFromStr(geom.Offset{X: 1, Y: 0}, game.StrikeRight, "O."),
				StrikeFromStr(geom.Offset{X: 1, Y: 0}, game.StrikeRightDown, ".O."),
				StrikeFromStr(geom.Offset{X: 1, Y: 0}, game.StrikeDown, ".O."),
			},
		},
		{
			"Two 2-len P1 and P2 strikes beside each other and extended after restriction",
			[]game.PlayerMove{
				{Cell: geom.Offset{X: 0, Y: 0}, Player: game.P1},
				{Cell: geom.Offset{X: 0, Y: -1}, Player: game.P2},
				{Cell: geom.Offset{X: 0, Y: 1}, Player: game.P1},
				{Cell: geom.Offset{X: 0, Y: -2}, Player: game.P2},
			},
			[]game.Strike{
				// P1
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 1}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRight, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 1}, game.StrikeRight, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightDown, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 1}, game.StrikeRightDown, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeDown, "XX."),

				// P2
				StrikeFromStr(geom.Offset{X: 0, Y: -1}, game.StrikeRightUp, ".O."),
				StrikeFromStr(geom.Offset{X: 0, Y: -2}, game.StrikeRightUp, ".O."),
				StrikeFromStr(geom.Offset{X: 0, Y: -1}, game.StrikeRight, ".O."),
				StrikeFromStr(geom.Offset{X: 0, Y: -2}, game.StrikeRight, ".O."),
				StrikeFromStr(geom.Offset{X: 0, Y: -1}, game.StrikeRightDown, ".O."),
				StrikeFromStr(geom.Offset{X: 0, Y: -2}, game.StrikeRightDown, ".O."),
				StrikeFromStr(geom.Offset{X: 0, Y: -2}, game.StrikeDown, ".OO"),
			},
		},
		{
			"Two 2-len P1 and P2 strikes beside each other and extended before restriction",
			[]game.PlayerMove{
				{Cell: geom.Offset{X: 0, Y: 1}, Player: game.P1},
				{Cell: geom.Offset{X: 0, Y: -2}, Player: game.P2},
				{Cell: geom.Offset{X: 0, Y: 0}, Player: game.P1},
				{Cell: geom.Offset{X: 0, Y: -1}, Player: game.P2},
			},
			[]game.Strike{
				// P1
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 1}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRight, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 1}, game.StrikeRight, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightDown, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 1}, game.StrikeRightDown, ".X."),

				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeDown, "XX."),

				// P2
				StrikeFromStr(geom.Offset{X: 0, Y: -1}, game.StrikeRightUp, ".O."),
				StrikeFromStr(geom.Offset{X: 0, Y: -2}, game.StrikeRightUp, ".O."),
				StrikeFromStr(geom.Offset{X: 0, Y: -1}, game.StrikeRight, ".O."),
				StrikeFromStr(geom.Offset{X: 0, Y: -2}, game.StrikeRight, ".O."),
				StrikeFromStr(geom.Offset{X: 0, Y: -1}, game.StrikeRightDown, ".O."),
				StrikeFromStr(geom.Offset{X: 0, Y: -2}, game.StrikeRightDown, ".O."),

				StrikeFromStr(geom.Offset{X: 0, Y: -2}, game.StrikeDown, ".OO"),
			},
		},
		{
			"Completely restrict a 2-len strike",
			[]game.PlayerMove{
				{Cell: geom.Offset{X: 0, Y: 0}, Player: game.P1},
				{Cell: geom.Offset{X: -1, Y: -1}, Player: game.P2},
				{Cell: geom.Offset{X: 1, Y: 1}, Player: game.P1},
				{Cell: geom.Offset{X: 2, Y: 2}, Player: game.P2},
			},
			[]game.Strike{
				// P1
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRight, ".X."),
				StrikeFromStr(geom.Offset{X: 1, Y: 1}, game.StrikeRight, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 1, Y: 1}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeDown, ".X."),
				StrikeFromStr(geom.Offset{X: 1, Y: 1}, game.StrikeDown, ".X."),

				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightDown, "XX"),

				// P2
				StrikeFromStr(geom.Offset{X: -1, Y: -1}, game.StrikeRight, ".O."),
				StrikeFromStr(geom.Offset{X: 2, Y: 2}, game.StrikeRight, ".O."),
				StrikeFromStr(geom.Offset{X: -1, Y: -1}, game.StrikeRightUp, ".O."),
				StrikeFromStr(geom.Offset{X: 2, Y: 2}, game.StrikeRightUp, ".O."),
				StrikeFromStr(geom.Offset{X: -1, Y: -1}, game.StrikeDown, ".O."),
				StrikeFromStr(geom.Offset{X: 2, Y: 2}, game.StrikeDown, ".O."),

				StrikeFromStr(geom.Offset{X: -1, Y: -1}, game.StrikeRightDown, ".O"),
				StrikeFromStr(geom.Offset{X: 2, Y: 2}, game.StrikeRightDown, "O."),
			},
		},
		{
			"Merge 2 restricted 1-len strikes",
			[]game.PlayerMove{
				{Cell: geom.Offset{X: -1, Y: 0}, Player: game.P1},
				{Cell: geom.Offset{X: -2, Y: 0}, Player: game.P2},
				{Cell: geom.Offset{X: 1, Y: 0}, Player: game.P1},
				{Cell: geom.Offset{X: 2, Y: 0}, Player: game.P2},
				{Cell: geom.Offset{X: 0, Y: 0}, Player: game.P1},
			},
			[]game.Strike{
				// P1
				StrikeFromStr(geom.Offset{X: -1, Y: 0}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 1, Y: 0}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: -1, Y: 0}, game.StrikeRightDown, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightDown, ".X."),
				StrikeFromStr(geom.Offset{X: 1, Y: 0}, game.StrikeRightDown, ".X."),
				StrikeFromStr(geom.Offset{X: -1, Y: 0}, game.StrikeDown, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeDown, ".X."),
				StrikeFromStr(geom.Offset{X: 1, Y: 0}, game.StrikeDown, ".X."),

				StrikeFromStr(geom.Offset{X: -1, Y: 0}, game.StrikeRight, "XXX"),

				// P2
				StrikeFromStr(geom.Offset{X: -2, Y: 0}, game.StrikeRightUp, ".O."),
				StrikeFromStr(geom.Offset{X: 2, Y: 0}, game.StrikeRightUp, ".O."),
				StrikeFromStr(geom.Offset{X: -2, Y: 0}, game.StrikeRightDown, ".O."),
				StrikeFromStr(geom.Offset{X: 2, Y: 0}, game.StrikeRightDown, ".O."),
				StrikeFromStr(geom.Offset{X: -2, Y: 0}, game.StrikeDown, ".O."),
				StrikeFromStr(geom.Offset{X: 2, Y: 0}, game.StrikeDown, ".O."),

				StrikeFromStr(geom.Offset{X: -2, Y: 0}, game.StrikeRight, ".O"),
				StrikeFromStr(geom.Offset{X: 2, Y: 0}, game.StrikeRight, "O."),
			},
		},
		{
			"1-len strike after a merged strike",
			[]game.PlayerMove{
				{Cell: geom.Offset{X: -1, Y: 0}, Player: game.P1},
				{Cell: geom.Offset{X: 1, Y: 0}, Player: game.P1},
				{Cell: geom.Offset{X: 0, Y: 0}, Player: game.P1},
				{Cell: geom.Offset{X: 2, Y: 2}, Player: game.P1},
			},
			[]game.Strike{
				// Right-Up
				StrikeFromStr(geom.Offset{X: -1, Y: 0}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 1, Y: 0}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 2, Y: 2}, game.StrikeRightUp, ".X."),

				// Right
				StrikeFromStr(geom.Offset{X: -1, Y: 0}, game.StrikeRight, ".XXX."),
				StrikeFromStr(geom.Offset{X: 2, Y: 2}, game.StrikeRight, ".X."),

				// Right-Down
				StrikeFromStr(geom.Offset{X: -1, Y: 0}, game.StrikeRightDown, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightDown, ".X."),
				StrikeFromStr(geom.Offset{X: 1, Y: 0}, game.StrikeRightDown, ".X."),
				StrikeFromStr(geom.Offset{X: 2, Y: 2}, game.StrikeRightDown, ".X."),

				// Down
				StrikeFromStr(geom.Offset{X: -1, Y: 0}, game.StrikeDown, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeDown, ".X."),
				StrikeFromStr(geom.Offset{X: 1, Y: 0}, game.StrikeDown, ".X."),
				StrikeFromStr(geom.Offset{X: 2, Y: 2}, game.StrikeDown, ".X."),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			set := game.NewStrikeSet()

			for _, move := range test.moves {
				set.MakeMove(move.Cell, move.Player)
			}

			got := set.Strikes()

			td.Cmp(t, got, td.Bag(td.Flatten(test.want)))
		})
	}
}

func TestStrikeSetMarkUnoccupied(t *testing.T) {
	t.Run("undoing a single move results in no strikes", func(t *testing.T) {
		set := game.NewStrikeSet()

		set.MakeMove(geom.Offset{X: 0, Y: 0}, game.P1)
		err := set.MarkUnoccupied(geom.Offset{X: 0, Y: 0})
		if err != nil {
			t.Errorf("got error [%v], want none", err)
		}

		got := set.Strikes()
		if got != nil {
			t.Errorf("got %v, wanted nil", got)
		}
	})

	tests := []struct {
		description string
		moves       []game.PlayerMove
		toRevert    geom.Offset
		want        []game.Strike
	}{
		{
			"remove a lone 1-len strike",
			[]game.PlayerMove{
				{Cell: geom.Offset{X: 0, Y: 0}, Player: game.P1},
				{Cell: geom.Offset{X: 2, Y: 2}, Player: game.P1},
			},
			geom.Offset{X: 2, Y: 2},
			[]game.Strike{
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRight, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightDown, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeDown, ".X."),
			},
		},
		{
			"cut a strike formed by appending moves",
			[]game.PlayerMove{
				{Cell: geom.Offset{X: 0, Y: 0}, Player: game.P1},
				{Cell: geom.Offset{X: 1, Y: 0}, Player: game.P1},
				{Cell: geom.Offset{X: 2, Y: 0}, Player: game.P1},
			},
			geom.Offset{X: 1, Y: 0},
			[]game.Strike{
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRight, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightDown, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeDown, ".X."),

				StrikeFromStr(geom.Offset{X: 2, Y: 0}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 2, Y: 0}, game.StrikeRight, ".X."),
				StrikeFromStr(geom.Offset{X: 2, Y: 0}, game.StrikeRightDown, ".X."),
				StrikeFromStr(geom.Offset{X: 2, Y: 0}, game.StrikeDown, ".X."),
			},
		},
		{
			"cut a strike formed by a merge",
			[]game.PlayerMove{
				{Cell: geom.Offset{X: 0, Y: 0}, Player: game.P1},
				{Cell: geom.Offset{X: 2, Y: 0}, Player: game.P1},
				{Cell: geom.Offset{X: 1, Y: 0}, Player: game.P1},
			},
			geom.Offset{X: 1, Y: 0},
			[]game.Strike{
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRight, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightDown, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeDown, ".X."),

				StrikeFromStr(geom.Offset{X: 2, Y: 0}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 2, Y: 0}, game.StrikeRight, ".X."),
				StrikeFromStr(geom.Offset{X: 2, Y: 0}, game.StrikeRightDown, ".X."),
				StrikeFromStr(geom.Offset{X: 2, Y: 0}, game.StrikeDown, ".X."),
			},
		},
		{
			"remove cell from a strike's beginning",
			[]game.PlayerMove{
				{Cell: geom.Offset{X: 0, Y: 0}, Player: game.P1},
				{Cell: geom.Offset{X: 1, Y: 1}, Player: game.P1},
				{Cell: geom.Offset{X: 2, Y: 2}, Player: game.P1},
			},
			geom.Offset{X: 0, Y: 0},
			[]game.Strike{
				StrikeFromStr(geom.Offset{X: 1, Y: 1}, game.StrikeRightDown, ".XX."),

				StrikeFromStr(geom.Offset{X: 1, Y: 1}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 1, Y: 1}, game.StrikeRight, ".X."),
				StrikeFromStr(geom.Offset{X: 1, Y: 1}, game.StrikeDown, ".X."),

				StrikeFromStr(geom.Offset{X: 2, Y: 2}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 2, Y: 2}, game.StrikeRight, ".X."),
				StrikeFromStr(geom.Offset{X: 2, Y: 2}, game.StrikeDown, ".X."),
			},
		},
		{
			"remove cell from a strike's end",
			[]game.PlayerMove{
				{Cell: geom.Offset{X: 0, Y: 0}, Player: game.P1},
				{Cell: geom.Offset{X: 1, Y: -1}, Player: game.P1},
				{Cell: geom.Offset{X: 2, Y: -2}, Player: game.P1},
			},
			geom.Offset{X: 2, Y: -2},
			[]game.Strike{
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightUp, ".XX."),

				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRight, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightDown, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeDown, ".X."),

				StrikeFromStr(geom.Offset{X: 1, Y: -1}, game.StrikeRight, ".X."),
				StrikeFromStr(geom.Offset{X: 1, Y: -1}, game.StrikeRightDown, ".X."),
				StrikeFromStr(geom.Offset{X: 1, Y: -1}, game.StrikeDown, ".X."),
			},
		},
		{
			"derestrict a 1-len strike from a single opponent move from the beginning",
			[]game.PlayerMove{
				{Cell: geom.Offset{X: 0, Y: 0}, Player: game.P1},
				{Cell: geom.Offset{X: -1, Y: 0}, Player: game.P2},
			},
			geom.Offset{X: -1, Y: 0},
			[]game.Strike{
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRight, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightDown, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeDown, ".X."),
			},
		},
		{
			"derestrict a 1-len strike from a single opponent move from the end",
			[]game.PlayerMove{
				{Cell: geom.Offset{X: 0, Y: 0}, Player: game.P1},
				{Cell: geom.Offset{X: 1, Y: 0}, Player: game.P2},
			},
			geom.Offset{X: 1, Y: 0},
			[]game.Strike{
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightUp, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRight, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeRightDown, ".X."),
				StrikeFromStr(geom.Offset{X: 0, Y: 0}, game.StrikeDown, ".X."),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			set := game.NewStrikeSet()

			for _, move := range test.moves {
				set.MakeMove(move.Cell, move.Player)
			}

			set.MarkUnoccupied(test.toRevert)

			got := set.Strikes()

			td.Cmp(t, got, td.Bag(td.Flatten(test.want)))
		})
	}
}

func TestStrikeSetStrikesThrough(t *testing.T) {
	set := game.NewStrikeSet()

	// 21CB <- Strike merge
	// 3..A
	// 4..9
	// 5678
	set.MakeMove(geom.Offset{X: 2, Y: 0}, game.P1)
	set.MakeMove(geom.Offset{X: 3, Y: 0}, game.P1)
	set.MakeMove(geom.Offset{X: 3, Y: 1}, game.P1)
	set.MakeMove(geom.Offset{X: 3, Y: 2}, game.P1)
	set.MakeMove(geom.Offset{X: 3, Y: 3}, game.P1)
	set.MakeMove(geom.Offset{X: 2, Y: 3}, game.P1)
	set.MakeMove(geom.Offset{X: 1, Y: 3}, game.P1)
	set.MakeMove(geom.Offset{X: 0, Y: 3}, game.P1)
	set.MakeMove(geom.Offset{X: 0, Y: 2}, game.P1)
	set.MakeMove(geom.Offset{X: 0, Y: 1}, game.P1)
	set.MakeMove(geom.Offset{X: 0, Y: 0}, game.P1)
	set.MakeMove(geom.Offset{X: 1, Y: 0}, game.P1)

	strikes := set.StrikesThrough(geom.Offset{X: 3, Y: 0})
	gotStrike := strikes[game.StrikeRight.FixedID]
	wantStrike := game.Strike{
		Player:           game.P1,
		Dir:              game.StrikeRight,
		Start:            geom.Offset{X: 0, Y: 0},
		Len:              4,
		ExtendableBefore: true,
		ExtendableAfter:  true,
	}

	td.Cmp(t, gotStrike, wantStrike)
}
