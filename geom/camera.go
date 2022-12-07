package geom

import "fmt"

type Camera struct {
	View       Rect
	TrackDepth int
}

func (c Camera) InnerView() Rect {
	inner := Rect{
		X: c.View.X + c.TrackDepth,
		Y: c.View.Y + c.TrackDepth,
		W: c.View.W - 2*c.TrackDepth,
		H: c.View.H - 2*c.TrackDepth,
	}

	if c.View.W <= 2*c.TrackDepth {
		inner.X = c.View.X + (c.View.W-1)/2
		inner.W = 1
	}

	if c.View.H <= 2*c.TrackDepth {
		inner.Y = c.View.Y + (c.View.H-1)/2
		inner.H = 1
	}

	return inner
}

func (c Camera) NudgeTo(pos Offset) Camera {
	// inner := c.InnerView()

	return c
}

func (c Camera) String() string {
	return fmt.Sprintf("(camera %v; depth: %v)", c.View, c.TrackDepth)
}
