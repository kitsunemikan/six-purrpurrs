package geom

import "fmt"

type Rect struct {
	X, Y, W, H int
}

func NewRectFromOffsets(pos, size Offset) Rect {
	return Rect{
		X: pos.X,
		Y: pos.Y,
		W: size.X,
		H: size.Y,
	}
}

func (r Rect) Move(ds Offset) Rect {
	r.X += ds.X
	r.Y += ds.Y
	return r
}

func (r Rect) TopLeft() Offset {
	return Offset{r.X, r.Y}
}

func (r Rect) Dimensions() Offset {
	return Offset{r.W, r.H}
}

func (r Rect) IsEqual(other Rect) bool {
	return r.X == other.X && r.Y == other.Y && r.W == other.W && r.H == other.H
}

func (r Rect) ToWorld(local Offset) Offset {
	return Offset{r.X + local.X, r.Y + local.Y}
}

func (r Rect) ToWorldXY(x, y int) Offset {
	return Offset{r.X + x, r.Y + y}
}

func (r Rect) ToLocal(world Offset) Offset {
	return Offset{world.X - r.X, world.Y - r.Y}
}

func (r Rect) Area() int {
	return r.W * r.H
}

func (r Rect) CenterOn(pos Offset) Rect {
	return Rect{
		X: pos.X - r.W/2,
		Y: pos.Y - r.H/2,
		W: r.W,
		H: r.H,
	}
}

func (r Rect) Center() Offset {
	return Offset{r.X + r.W/2, r.Y + r.H/2}
}

func (r Rect) IsInsideRect(other Rect) bool {
	return other.X <= r.X && r.X+r.W <= other.X+other.W &&
		other.Y <= r.Y && r.Y+r.H <= other.Y+other.H
}

// SnapInto will snap rectangle to the closest boundary of the
// bound rectangle. If rectangle is inside bound rectangle, the
// same rectangle is returned. If bound rectangle is smaller than
// the rectangle, the rectangle is centered on the bound rectangle.
func (r Rect) SnapInto(bound Rect) Rect {
	if r.W >= bound.W {
		r.X = bound.Center().X - r.W/2
	} else if bound.X-r.X > 0 && r.X+r.W-bound.X-bound.W < 0 {
		// It's poking out of bound rect from the left
		r.X = bound.X
	} else if bound.X-r.X < 0 && r.X+r.W-bound.X-bound.W > 0 {
		// It's poking out of bound rect from the right
		r.X = bound.X + bound.W - r.W
	} else {
		// It's already inside
	}

	if r.H >= bound.H {
		r.Y = bound.Center().Y - r.H/2
	} else if bound.Y-r.Y > 0 && r.Y+r.H-bound.Y-bound.H < 0 {
		// It's poking out of bound rect from the left
		r.Y = bound.Y
	} else if bound.Y-r.Y < 0 && r.Y+r.H-bound.Y-bound.H > 0 {
		// It's poking out of bound rect from the right
		r.Y = bound.Y + bound.H - r.H
	} else {
		// It's already inside
	}

	return r
}

func (r Rect) GrowToContainOffset(pos Offset) Rect {
	if pos.X < r.X {
		r.W += r.X - pos.X
		r.X = pos.X
	} else if pos.X >= r.X+r.W {
		r.W += pos.X - r.X - r.W + 1
	}

	if pos.Y < r.Y {
		r.H += r.Y - pos.Y
		r.Y = pos.Y
	} else if pos.Y >= r.Y+r.H {
		r.H += pos.Y - r.Y - r.H + 1
	}

	return r
}

func (r Rect) GrowToContainRect(other Rect) Rect {
	if other.X < r.X {
		r.W += r.X - other.X
		r.X = other.X
	}

	if other.X+other.W > r.X+r.W {
		r.W = other.X + other.W - r.X
	}

	if other.Y < r.Y {
		r.H += r.Y - other.Y
		r.Y = other.Y
	}

	if other.Y+other.H > r.Y+r.H {
		r.H = other.Y + other.H - r.Y
	}

	return r
}

func (r Rect) String() string {
	return fmt.Sprintf("(%d;%d) - (%d;%d)", r.X, r.Y, r.X+r.W, r.Y+r.H)
}
