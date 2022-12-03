package game

import (
	. "github.com/kitsunemikan/ttt-cli/geom"
)

type PlayerAgent interface {
	MakeMove(*BoardState) Offset
}
