package geom

import "fmt"

type Offset struct {
	X, Y int
}

func (a Offset) Add(b Offset) Offset {
	return Offset{a.X + b.X, a.Y + b.Y}
}

func (a Offset) AddXY(dx, dy int) Offset {
	return Offset{a.X + dx, a.Y + dy}
}

func (a Offset) Sub(b Offset) Offset {
	return Offset{a.X - b.X, a.Y - b.Y}
}

func (a Offset) SubXY(dx, dy int) Offset {
	return Offset{a.X - dx, a.Y - dy}
}

func (a Offset) ScaleUp(c int) Offset {
	return Offset{c * a.X, c * a.Y}
}

func (a Offset) ScaleDown(c int) Offset {
	return Offset{a.X / c, a.Y / c}
}

func (a Offset) Area() int {
	return a.X * a.Y
}

func (a Offset) IsInsideCircle(radius int) bool {
	radius += 1
	return a.X*a.X+a.Y*a.Y <= radius*radius
}

func (a Offset) IsInsideRect(r Rect) bool {
	return r.X <= a.X && a.X < r.X+r.W && r.Y <= a.Y && a.Y < r.Y+r.H
}

func (a Offset) IsEqual(b Offset) bool {
	return a.X == b.X && a.Y == b.Y
}

func (a Offset) IsZero() bool {
	return a.X == 0 && a.Y == 0
}

func (a Offset) SnapIntoRect(r Rect) Offset {
	if a.X < r.X {
		a.X = r.X
	} else if a.X >= r.X+r.W {
		a.X = r.X + r.W - 1
	}

	if a.Y < r.Y {
		a.Y = r.Y
	} else if a.Y >= r.Y+r.H {
		a.Y = r.Y + r.H - 1
	}

	return a
}

func (a Offset) String() string {
	return fmt.Sprintf("(%v;%v)", a.X, a.Y)
}
