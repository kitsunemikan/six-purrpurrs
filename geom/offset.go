package geom

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

func (a Offset) IsEqual(b Offset) bool {
	return a.X == b.X && a.Y == b.Y
}

func (a Offset) IsZero() bool {
	return a.X == 0 && a.Y == 0
}
