package game_test

import (
	"reflect"
	"testing"

	"github.com/kitsunemikan/six-purrpurrs/game"
	"github.com/kitsunemikan/six-purrpurrs/geom"
)

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
			{
				Player:           game.P1,
				Start:            geom.Offset{X: 0, Y: 0},
				Dir:              game.StrikeRight,
				Len:              1,
				ExtendableBefore: true,
				ExtendableAfter:  true,
			},
			{
				Player:           game.P1,
				Start:            geom.Offset{X: 0, Y: 0},
				Dir:              game.StrikeUpRight,
				Len:              1,
				ExtendableBefore: true,
				ExtendableAfter:  true,
			},
			{
				Player:           game.P1,
				Start:            geom.Offset{X: 0, Y: 0},
				Dir:              game.StrikeDownRight,
				Len:              1,
				ExtendableBefore: true,
				ExtendableAfter:  true,
			},
			{
				Player:           game.P1,
				Start:            geom.Offset{X: 0, Y: 0},
				Dir:              game.StrikeDown,
				Len:              1,
				ExtendableBefore: true,
				ExtendableAfter:  true,
			},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
