package geom_test

import (
	"testing"

	"github.com/kitsunemikan/ttt-cli/geom"
)

func TestCameraInnerView(t *testing.T) {
	cases := []struct {
		Desc   string
		Camera geom.Camera
		Want   geom.Rect
	}{
		{
			"inner view of depth 0 is the same as camera view",
			geom.Camera{
				View:       geom.Rect{X: -1, Y: -2, W: 3, H: 3},
				TrackDepth: 0,
			},
			geom.Rect{X: -1, Y: -2, W: 3, H: 3},
		},
		{
			"depth of one shrinks size from all sides by one",
			geom.Camera{
				View:       geom.Rect{X: -1, Y: -2, W: 5, H: 5},
				TrackDepth: 1,
			},
			geom.Rect{X: 0, Y: -1, W: 3, H: 3},
		},
		{
			"depth that makes odd square camera into a point",
			geom.Camera{
				View:       geom.Rect{X: 0, Y: 0, W: 5, H: 5},
				TrackDepth: 2,
			},
			geom.Rect{X: 2, Y: 2, W: 1, H: 1},
		},
		{
			"overly big depth still makes odd square camera into a point",
			geom.Camera{
				View:       geom.Rect{X: 0, Y: 0, W: 5, H: 5},
				TrackDepth: 3,
			},
			geom.Rect{X: 2, Y: 2, W: 1, H: 1},
		},
		{
			"overly big depth still makes odd rect camera into a point",
			geom.Camera{
				View:       geom.Rect{X: 0, Y: 0, W: 7, H: 5},
				TrackDepth: 3,
			},
			geom.Rect{X: 3, Y: 2, W: 1, H: 1},
		},
		{
			"overly big depth clamps dimensions independently",
			geom.Camera{
				View:       geom.Rect{X: 0, Y: 0, W: 3, H: 7},
				TrackDepth: 2,
			},
			geom.Rect{X: 1, Y: 2, W: 1, H: 3},
		},
		{
			"overly big depth makes even square into upper-left point",
			geom.Camera{
				View:       geom.Rect{X: 0, Y: 0, W: 2, H: 2},
				TrackDepth: 1,
			},
			geom.Rect{X: 0, Y: 0, W: 1, H: 1},
		},
		{
			"overly big depth independently makes dimensions prefer negative side",
			geom.Camera{
				View:       geom.Rect{X: 0, Y: 0, W: 6, H: 2},
				TrackDepth: 2,
			},
			geom.Rect{X: 2, Y: 0, W: 2, H: 1},
		},
	}

	for _, test := range cases {
		t.Run(test.Desc, func(t *testing.T) {
			got := test.Camera.InnerView()

			if !got.IsEqual(test.Want) {
				t.Errorf("for camera %v: got %v, want %v", test.Camera, got, test.Want)
			}
		})
	}
}
