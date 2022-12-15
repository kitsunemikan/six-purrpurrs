package game

import "github.com/kitsunemikan/six-purrpurrs/geom"

type StrikeDir geom.Offset

var (
	StrikeUpRight   = StrikeDir{X: 1, Y: -1}
	StrikeRight     = StrikeDir{X: 1, Y: 0}
	StrikeDownRight = StrikeDir{X: 1, Y: 1}
	StrikeDown      = StrikeDir{X: 0, Y: 1}
)

var StrikeDirs = []StrikeDir{StrikeUpRight, StrikeRight, StrikeDownRight, StrikeDown}

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
			Dir:              StrikeUpRight,
			Len:              1,
			ExtendableBefore: true,
			ExtendableAfter:  true,
		},
		{
			Player:           P1,
			Start:            geom.Offset{X: 0, Y: 0},
			Dir:              StrikeDownRight,
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
