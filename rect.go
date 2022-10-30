package main

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

func (r Rect) IsOffsetInside(pos Offset) bool {
	return r.X <= pos.X && pos.X < r.X+r.W && r.Y <= pos.Y && pos.Y < r.Y+r.H
}
