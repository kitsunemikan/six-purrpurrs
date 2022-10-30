package main

type Offset struct {
	X, Y int
}

func (a Offset) Add(b Offset) Offset {
	return Offset{a.X + b.X, a.Y + b.Y}
}

func (a Offset) Sub(b Offset) Offset {
	return Offset{a.X - b.X, a.Y - b.Y}
}

func (a Offset) ScaleUp(c int) Offset {
	return Offset{c * a.X, c * a.Y}
}

func (a Offset) ScaleDown(c int) Offset {
	return Offset{a.X / c, a.Y / c}
}

// Not counting the exact bottomRight
func (a Offset) IsInsideRect(topLeft, bottomRight Offset) bool {
	return topLeft.X <= a.X && a.X < bottomRight.X && topLeft.Y <= a.Y && a.Y < bottomRight.Y
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
