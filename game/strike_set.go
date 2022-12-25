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
	strikes []Strike

	board map[geom.Offset][]int
}

func NewStrikeSet() *StrikeSet {
	return &StrikeSet{
		strikes: nil,
		board:   make(map[geom.Offset][]int),
	}
}

// It is assumed that the board is filled only with unoccupied cells, and invalid cells don't exist
func (s *StrikeSet) MakeMove(move PlayerMove) error {
	for _, dir := range StrikeDirs {
		if strikes, ok := s.board[move.Cell.Sub(dir.Offset())]; ok {
			for _, strikeID := range strikes {
				if s.strikes[strikeID].Dir.IsEqual(dir) {
					// Found
					s.board[move.Cell] = append(s.board[move.Cell], strikeID)
					s.strikes[strikeID].Len++
				}
			}
		} else if strikes, ok := s.board[move.Cell.Add(dir.Offset())]; ok {
			for _, strikeID := range strikes {
				if s.strikes[strikeID].Dir.IsEqual(dir) {
					// Found
					s.board[move.Cell] = append(s.board[move.Cell], strikeID)
					s.strikes[strikeID].Start = move.Cell
					s.strikes[strikeID].Len++
				}
			}
		} else {
			s.strikes = append(s.strikes, Strike{
				Player:           move.ID,
				Start:            move.Cell,
				Len:              1,
				Dir:              dir,
				ExtendableBefore: true,
				ExtendableAfter:  true,
			})

			s.board[move.Cell] = append(s.board[move.Cell], len(s.strikes)-1)
		}
	}

	return nil
}

func (s *StrikeSet) Strikes() []Strike {
	return s.strikes
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

func (dir StrikeDir) Offset() geom.Offset {
	return geom.Offset(dir)
}

func (dir StrikeDir) IsEqual(other StrikeDir) bool {
	return dir.Offset().IsEqual(other.Offset())
}
