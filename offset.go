package main

type Offset struct {
	X, Y int
}

func (a Offset) Add(b Offset) Offset {
	return Offset{a.X + b.X, a.Y + b.Y}
}

func (a Offset) Scale(c int) Offset {
	return Offset{c * a.X, c * a.Y}
}
