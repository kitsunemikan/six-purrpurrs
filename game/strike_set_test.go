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
	t.Run("empty strike set returns nil strike slice", func(t *testing.T) {
		set := &game.StrikeSet{}

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
				{Cell: geom.Offset{X: 0, Y: 0}, ID: game.P1},
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
				{Cell: geom.Offset{X: 0, Y: 0}, ID: game.P2},
				{Cell: geom.Offset{X: 3, Y: 3}, ID: game.P1},
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
				{Cell: geom.Offset{X: 0, Y: 0}, ID: game.P1},
				{Cell: geom.Offset{X: 1, Y: 0}, ID: game.P1},
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
				{Cell: geom.Offset{X: 0, Y: 0}, ID: game.P1},
				{Cell: geom.Offset{X: 0, Y: -1}, ID: game.P1},
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
				{Cell: geom.Offset{X: 0, Y: 0}, ID: game.P1},
				{Cell: geom.Offset{X: 2, Y: 2}, ID: game.P1},
				{Cell: geom.Offset{X: 1, Y: 1}, ID: game.P1},
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
				{Cell: geom.Offset{X: 0, Y: 0}, ID: game.P2},
				{Cell: geom.Offset{X: 2, Y: 0}, ID: game.P2},
				{Cell: geom.Offset{X: 2, Y: 2}, ID: game.P2},
				{Cell: geom.Offset{X: 0, Y: 2}, ID: game.P2},

				// OOO
				// O.O
				// OOO
				{Cell: geom.Offset{X: 1, Y: 0}, ID: game.P2},
				{Cell: geom.Offset{X: 2, Y: 1}, ID: game.P2},
				{Cell: geom.Offset{X: 1, Y: 2}, ID: game.P2},
				{Cell: geom.Offset{X: 0, Y: 1}, ID: game.P2},

				// OOO
				// OOO
				// OOO
				{Cell: geom.Offset{X: 1, Y: 1}, ID: game.P2},
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
				{Cell: geom.Offset{X: 0, Y: 0}, ID: game.P1},
				{Cell: geom.Offset{X: 2, Y: 0}, ID: game.P1},
				{Cell: geom.Offset{X: 4, Y: 0}, ID: game.P1},
				{Cell: geom.Offset{X: 6, Y: 0}, ID: game.P1},

				// xxx.xxx
				{Cell: geom.Offset{X: 1, Y: 0}, ID: game.P1},
				{Cell: geom.Offset{X: 5, Y: 0}, ID: game.P1},

				// XXXXXXX
				{Cell: geom.Offset{X: 3, Y: 0}, ID: game.P1},
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
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			set := game.NewStrikeSet()

			for _, move := range test.moves {
				set.MakeMove(move)
			}

			got := set.Strikes()

			td.Cmp(t, got, td.Bag(td.Flatten(test.want)))
		})
	}
}
