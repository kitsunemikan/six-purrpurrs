package game

import (
	. "github.com/kitsunemikan/six-purrpurrs/geom"
)

type PlayerAgent interface {
	MakeMove(*GameState) Offset
}
