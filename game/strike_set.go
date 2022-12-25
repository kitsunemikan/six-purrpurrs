package game

import (
	"fmt"
	"strings"

	"github.com/kitsunemikan/six-purrpurrs/geom"
)

type StrikeDir struct {
	X, Y    int
	fixedID int
}

var (
	StrikeRightUp   = StrikeDir{X: 1, Y: -1, fixedID: 0}
	StrikeRight     = StrikeDir{X: 1, Y: 0, fixedID: 1}
	StrikeRightDown = StrikeDir{X: 1, Y: 1, fixedID: 2}
	StrikeDown      = StrikeDir{X: 0, Y: 1, fixedID: 3}
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
		// Create reference arary, if it's a new cell
		if _, ok := s.board[move.Cell]; !ok {
			strikeRef := make([]int, len(StrikeDirs))
			for i := range strikeRef {
				strikeRef[i] = -1
			}

			s.board[move.Cell] = strikeRef
		}

		if beforeStrikes, ok := s.board[move.Cell.Sub(dir.Offset())]; ok {
			// Try to find a strike in the opposite direction
			// In this case, the strike we find will have its length extended by one

			if strikeID := beforeStrikes[dir.fixedID]; strikeID != -1 {
				// Found
				s.board[move.Cell][dir.fixedID] = strikeID
				s.strikes[strikeID].Len++

			}
		} else if afterStrikes, ok := s.board[move.Cell.Add(dir.Offset())]; ok {
			// Try to find a strike in the strike direction
			// In this case, we will become a new starting cell for the strike
			// and the length will be extended

			if strikeID := afterStrikes[dir.fixedID]; strikeID != -1 {
				// Found
				s.board[move.Cell][dir.fixedID] = strikeID
				s.strikes[strikeID].Start = move.Cell
				s.strikes[strikeID].Len++

			}
		} else {
			// In case there's no already existing strikes nearby, create a new strike

			s.strikes = append(s.strikes, Strike{
				Player:           move.ID,
				Start:            move.Cell,
				Len:              1,
				Dir:              dir,
				ExtendableBefore: true,
				ExtendableAfter:  true,
			})

			s.board[move.Cell][dir.fixedID] = len(s.strikes) - 1
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
	return geom.Offset{X: dir.X, Y: dir.Y}
}

func (dir StrikeDir) IsEqual(other StrikeDir) bool {
	return dir.X == other.X && dir.Y == other.Y && dir.fixedID == other.fixedID
}
