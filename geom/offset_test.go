package geom_test

import (
	"testing"

	"github.com/kitsunemikan/six-purrpurrs/geom"
)

func TestOffsetSnapIntoRect(t *testing.T) {
	cases := []struct {
		Desc  string
		Point geom.Offset
		Rect  geom.Rect
		Want  geom.Offset
	}{
		{
			"Point inside rect doesn't change",
			geom.Offset{X: 1, Y: 2},
			geom.Rect{X: 0, Y: 0, W: 2, H: 3},
			geom.Offset{X: 1, Y: 2},
		},
		{
			"Point to the left",
			geom.Offset{X: -1, Y: 2},
			geom.Rect{X: 0, Y: 0, W: 2, H: 3},
			geom.Offset{X: 0, Y: 2},
		},
		{
			"Point to the right of inner view",
			geom.Offset{X: 2, Y: 2},
			geom.Rect{X: 0, Y: 0, W: 2, H: 3},
			geom.Offset{X: 1, Y: 2},
		},
		{
			"Point to the top of inner view",
			geom.Offset{X: 1, Y: -1},
			geom.Rect{X: 0, Y: 0, W: 2, H: 3},
			geom.Offset{X: 1, Y: 0},
		},
		{
			"Point to the bottom of inner view",
			geom.Offset{X: 1, Y: 3},
			geom.Rect{X: 0, Y: 0, W: 2, H: 3},
			geom.Offset{X: 1, Y: 2},
		},
		{
			"Point to the top-left of inner view",
			geom.Offset{X: -1, Y: -1},
			geom.Rect{X: 0, Y: 0, W: 2, H: 3},
			geom.Offset{X: 0, Y: 0},
		},
		{
			"Point to the top-right of inner view",
			geom.Offset{X: 2, Y: -1},
			geom.Rect{X: 0, Y: 0, W: 2, H: 3},
			geom.Offset{X: 1, Y: 0},
		},
		{
			"Point to the bottom-right of inner view",
			geom.Offset{X: 2, Y: 3},
			geom.Rect{X: 0, Y: 0, W: 2, H: 3},
			geom.Offset{X: 1, Y: 2},
		},
		{
			"Point to the bottom-right of inner view",
			geom.Offset{X: -1, Y: 3},
			geom.Rect{X: 0, Y: 0, W: 2, H: 3},
			geom.Offset{X: 0, Y: 2},
		},
	}

	for _, test := range cases {
		t.Run(test.Desc, func(t *testing.T) {
			got := test.Point.SnapIntoRect(test.Rect)

			if !got.IsEqual(test.Want) {
				t.Errorf("got %v, want %v, when snapping %v into %v", got, test.Want, test.Point, test.Rect)
			}
		})
	}
}
