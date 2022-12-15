package game

import (
	"fmt"
	"strings"

	"github.com/kitsunemikan/six-purrpurrs/geom"
)

type StrikeDir geom.Offset

var (
	StrikeRightUp   = StrikeDir{X: 1, Y: -1}
	StrikeRight     = StrikeDir{X: 1, Y: 0}
	StrikeRightDown = StrikeDir{X: 1, Y: 1}
	StrikeDown      = StrikeDir{X: 0, Y: 1}
)

var StrikeDirs = []StrikeDir{StrikeRightUp, StrikeRight, StrikeRightDown, StrikeDown}

type Strike struct {
	Player           PlayerID
	Start            geom.Offset
	Dir              StrikeDir
	Len              int
	ExtendableBefore bool
	ExtendableAfter  bool
}

type StrikeSet struct {
	moveCount int
}

// It is assumed that the board is filled only with unoccupied cells, and invalid cells don't exist
func (s *StrikeSet) MakeMove(move PlayerMove) error {
	s.moveCount++
	return nil
}

func (s *StrikeSet) Strikes() []Strike {
	if s.moveCount == 0 {
		return nil
	}

	return []Strike{
		{
			Player:           P1,
			Start:            geom.Offset{X: 0, Y: 0},
			Dir:              StrikeRight,
			Len:              1,
			ExtendableBefore: true,
			ExtendableAfter:  true,
		},
		{
			Player:           P1,
			Start:            geom.Offset{X: 0, Y: 0},
			Dir:              StrikeRightUp,
			Len:              1,
			ExtendableBefore: true,
			ExtendableAfter:  true,
		},
		{
			Player:           P1,
			Start:            geom.Offset{X: 0, Y: 0},
			Dir:              StrikeRightDown,
			Len:              1,
			ExtendableBefore: true,
			ExtendableAfter:  true,
		},
		{
			Player:           P1,
			Start:            geom.Offset{X: 0, Y: 0},
			Dir:              StrikeDown,
			Len:              1,
			ExtendableBefore: true,
			ExtendableAfter:  true,
		},
	}
}

func (dir StrikeDir) String() string {
	var str strings.Builder

	if dir.X >= 1 {
		str.WriteString("Right")
	} else if dir.X <= -1 {
		str.WriteString("Left")
	}

	if dir.Y >= 1 {
		str.WriteString("Down")
	} else if dir.Y <= -1 {
		str.WriteString("Up")
	}

	if dir.X < -1 || dir.X > 1 || dir.Y < -1 || dir.Y > 1 {
		str.WriteString(fmt.Sprint(geom.Offset{X: dir.X, Y: dir.Y}))
	}

	return str.String()
}
