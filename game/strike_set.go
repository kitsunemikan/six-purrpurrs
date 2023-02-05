package game

import (
	"errors"
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

	board   map[geom.Offset][]int
	players map[geom.Offset]PlayerID
}

func NewStrikeSet() *StrikeSet {
	return &StrikeSet{
		strikes: nil,
		board:   make(map[geom.Offset][]int),
		players: make(map[geom.Offset]PlayerID),
	}
}

// It is assumed that the board is filled only with unoccupied cells, and invalid cells don't exist
// TODO: add error handling
func (s *StrikeSet) MakeMove(move PlayerMove) error {
	if _, exists := s.players[move.Cell]; exists {
		// TODO: add test for the error + make sentinel + wrap error
		return errors.New("strike set: make move: move already done")
	}

	s.players[move.Cell] = move.Player

	for _, dir := range StrikeDirs {
		// Create reference arary, if it's a new cell
		if _, ok := s.board[move.Cell]; !ok {
			strikeRef := make([]int, len(StrikeDirs))
			for i := range strikeRef {
				strikeRef[i] = -1
			}

			s.board[move.Cell] = strikeRef
		}

		enemyBeforeStrikeID := -1
		beforeStrikeID := -1
		beforeCell := move.Cell.Sub(dir.Offset())
		if p, ok := s.players[beforeCell]; ok {
			if p == move.Player {
				beforeStrikeID = s.board[beforeCell][dir.fixedID]
			} else if p == move.Player.Other() {
				enemyBeforeStrikeID = s.board[beforeCell][dir.fixedID]
			}
		}

		enemyAfterStrikeID := -1
		afterStrikeID := -1
		afterCell := move.Cell.Add(dir.Offset())
		if p, ok := s.players[afterCell]; ok {
			if p == move.Player {
				afterStrikeID = s.board[afterCell][dir.fixedID]
			} else if p == move.Player.Other() {
				enemyAfterStrikeID = s.board[afterCell][dir.fixedID]
			}
		}

		switch {
		case beforeStrikeID != -1 && afterStrikeID == -1:
			// There's only one strike in the opposite direction
			// In this case, the strike we find will have its length extended by one

			s.board[move.Cell][dir.fixedID] = beforeStrikeID
			s.strikes[beforeStrikeID].Len++

		case beforeStrikeID == -1 && afterStrikeID != -1:
			// There's only one strike in the strike direction
			// In this case, we will become a new starting cell for the strike
			// and the length will be extended

			s.board[move.Cell][dir.fixedID] = afterStrikeID
			s.strikes[afterStrikeID].Start = move.Cell
			s.strikes[afterStrikeID].Len++

		case beforeStrikeID != -1 && afterStrikeID != -1:
			s.strikes[beforeStrikeID].Len += s.strikes[afterStrikeID].Len + 1
			s.strikes[beforeStrikeID].ExtendableAfter = s.strikes[afterStrikeID].ExtendableAfter

			s.board[move.Cell][dir.fixedID] = beforeStrikeID

			// Route after strike cells to the new extended beforeStrike
			// Note that there will be no references to the after strike after this
			// in the board map
			afterStrike := s.strikes[afterStrikeID]
			for cell, i := afterStrike.Start, 0; i < afterStrike.Len; i++ {
				s.board[cell][dir.fixedID] = beforeStrikeID
			}

			// We shouldn't literally remove strike from the strikes array,
			// since we'll need to update all map strike references, which is expensive.
			// Instead we'll make its length 0, meaning it's an invalid strike
			s.strikes[afterStrikeID].Len = 0

		case beforeStrikeID == -1 && afterStrikeID == -1:
			// In case there's no already existing strikes nearby, create a new strike

			s.strikes = append(s.strikes, Strike{
				Player:           move.Player,
				Start:            move.Cell,
				Len:              1,
				Dir:              dir,
				ExtendableBefore: true,
				ExtendableAfter:  true,
			})

			s.board[move.Cell][dir.fixedID] = len(s.strikes) - 1
		}

		assignedStrikeID := s.board[move.Cell][dir.fixedID]
		if enemyBeforeStrikeID != -1 {
			s.strikes[assignedStrikeID].ExtendableBefore = false

			s.strikes[enemyBeforeStrikeID].ExtendableAfter = false
		}

		if enemyAfterStrikeID != -1 {
			s.strikes[assignedStrikeID].ExtendableAfter = false

			s.strikes[enemyAfterStrikeID].ExtendableBefore = false
		}
	}

	return nil
}

func (s *StrikeSet) UndoLastMove() error {
	return nil
}

func (s *StrikeSet) StrikesUnfiltered() []Strike {
	return s.strikes
}

func (s *StrikeSet) Strikes() []Strike {
	strikes := make([]Strike, 0, len(s.strikes))
	for i := range s.strikes {
		if s.strikes[i].Len == 0 {
			continue
		}

		strikes = append(strikes, s.strikes[i])
	}

	if len(strikes) == 0 {
		return nil
	}

	return strikes
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
