package game

import "github.com/kitsunemikan/six-purrpurrs/geom"

type PlayerMove struct {
	Cell   geom.Offset
	Player PlayerID
}
