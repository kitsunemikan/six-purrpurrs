package game

import (
	"errors"
	"fmt"
	"strings"

	"github.com/kitsunemikan/six-purrpurrs/geom"
)

type StrikeDir struct {
	X, Y    int
	FixedID int
}

var (
	StrikeRightUp   = StrikeDir{X: 1, Y: -1, FixedID: 0}
	StrikeRight     = StrikeDir{X: 1, Y: 0, FixedID: 1}
	StrikeRightDown = StrikeDir{X: 1, Y: 1, FixedID: 2}
	StrikeDown      = StrikeDir{X: 0, Y: 1, FixedID: 3}
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

func (s *Strike) AsCells() []geom.Offset {
	cells := make([]geom.Offset, s.Len)

	for i, cell := 0, s.Start; i < s.Len; i, cell = i+1, cell.Add(s.Dir.Offset()) {
		cells[i] = cell
	}

	return cells
}

type StrikeSet struct {
	strikes        []Strike
	deletedStrikes []int

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
func (s *StrikeSet) MakeMove(atCell geom.Offset, as PlayerID) error {
	move := PlayerMove{Cell: atCell, Player: as}

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
				beforeStrikeID = s.board[beforeCell][dir.FixedID]
			} else if p == move.Player.Other() {
				enemyBeforeStrikeID = s.board[beforeCell][dir.FixedID]
			}
		}

		enemyAfterStrikeID := -1
		afterStrikeID := -1
		afterCell := move.Cell.Add(dir.Offset())
		if p, ok := s.players[afterCell]; ok {
			if p == move.Player {
				afterStrikeID = s.board[afterCell][dir.FixedID]
			} else if p == move.Player.Other() {
				enemyAfterStrikeID = s.board[afterCell][dir.FixedID]
			}
		}

		switch {
		case beforeStrikeID != -1 && afterStrikeID == -1:
			// There's only one strike in the opposite direction
			// In this case, the strike we find will have its length extended by one

			s.board[move.Cell][dir.FixedID] = beforeStrikeID
			s.strikes[beforeStrikeID].Len++

		case beforeStrikeID == -1 && afterStrikeID != -1:
			// There's only one strike in the strike direction
			// In this case, we will become a new starting cell for the strike
			// and the length will be extended

			s.board[move.Cell][dir.FixedID] = afterStrikeID
			s.strikes[afterStrikeID].Start = move.Cell
			s.strikes[afterStrikeID].Len++

		case beforeStrikeID != -1 && afterStrikeID != -1:
			s.strikes[beforeStrikeID].Len += s.strikes[afterStrikeID].Len + 1
			s.strikes[beforeStrikeID].ExtendableAfter = s.strikes[afterStrikeID].ExtendableAfter

			s.board[move.Cell][dir.FixedID] = beforeStrikeID

			// Route after strike cells to the new extended beforeStrike
			// Note that there will be no references to the after strike after this
			// in the board map
			afterStrike := s.strikes[afterStrikeID]
			for cell, i := afterStrike.Start, 0; i < afterStrike.Len; i, cell = i+1, cell.Add(dir.Offset()) {
				// TODO: buggg!!! cell is not updated!
				s.board[cell][dir.FixedID] = beforeStrikeID
			}

			// We shouldn't literally remove strike from the strikes array,
			// since we'll need to update all map strike references, which is expensive.
			// Instead we'll make its length 0, meaning it's an invalid strike
			s.strikes[afterStrikeID].Len = 0
			s.deletedStrikes = append(s.deletedStrikes, afterStrikeID)

		case beforeStrikeID == -1 && afterStrikeID == -1:
			// In case there's no already existing strikes nearby, create a new strike

			newStrikeID := len(s.strikes)
			if len(s.deletedStrikes) == 0 {
				s.strikes = append(s.strikes, Strike{})
			} else {
				newStrikeID = s.deletedStrikes[len(s.deletedStrikes)-1]
				s.deletedStrikes = s.deletedStrikes[:len(s.deletedStrikes)-1]
			}

			s.strikes[newStrikeID] = Strike{
				Player:           move.Player,
				Start:            move.Cell,
				Len:              1,
				Dir:              dir,
				ExtendableBefore: true,
				ExtendableAfter:  true,
			}

			s.board[move.Cell][dir.FixedID] = newStrikeID
		}

		assignedStrikeID := s.board[move.Cell][dir.FixedID]
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

func (s *StrikeSet) MarkUnoccupied(cell geom.Offset) error {
	if _, occupied := s.players[cell]; !occupied {
		// TODO: proper error
		return errors.New("strike set: mark unoccupied: cell is already unoccupied")
	}

	for _, dir := range StrikeDirs {
		strikeID := s.board[cell][dir.FixedID]
		s.board[cell][dir.FixedID] = -1

		// Derestrict oponent strikes if any
		if afterPlayer, moveExists := s.players[cell.Add(dir.Offset())]; moveExists {
			if afterPlayer.Other() == s.players[cell] {
				afterCell := cell.Add(dir.Offset())
				enemyAfterStrikeID := s.board[afterCell][dir.FixedID]

				s.strikes[enemyAfterStrikeID].ExtendableBefore = true
			}
		}

		if beforePlayer, moveExists := s.players[cell.Sub(dir.Offset())]; moveExists {
			if beforePlayer.Other() == s.players[cell] {
				beforeCell := cell.Sub(dir.Offset())
				enemyBeforeStrikeID := s.board[beforeCell][dir.FixedID]

				s.strikes[enemyBeforeStrikeID].ExtendableAfter = true
			}
		}

		// Determine the index of the cell to be removed inside the strike
		// We will handle different cases depending whether it's located on the sides
		// or somewhere in the middle
		ds := cell.Sub(s.strikes[strikeID].Start)
		shift := 0
		if shift < ds.X {
			shift = ds.X
		}

		if shift < ds.Y {
			shift = ds.Y
		}

		atStart := shift == 0
		atEnd := shift == s.strikes[strikeID].Len-1

		switch {
		case atStart && atEnd:
			// Removing a 1-len strike
			s.strikes[strikeID].Len = 0
			s.deletedStrikes = append(s.deletedStrikes, strikeID)

		case atStart && !atEnd:
			s.strikes[strikeID].Start = cell.Add(dir.Offset())
			s.strikes[strikeID].Len--
			s.strikes[strikeID].ExtendableBefore = true

		case !atStart && atEnd:
			s.strikes[strikeID].Len--
			s.strikes[strikeID].ExtendableAfter = true

		case !atStart && !atEnd:
			// NOTE: since this function can mark any cell as unoccupied,
			// it's not necessarily the case that we will have some spare deleted strikes.
			// E.g., we have a strike that was constructed by appending cells one after another
			// and then we "cut it in half" by calling this function
			newStrikeID := len(s.strikes)
			if len(s.deletedStrikes) > 0 {
				newStrikeID = s.deletedStrikes[len(s.deletedStrikes)-1]
				s.deletedStrikes = s.deletedStrikes[:len(s.deletedStrikes)-1]
			} else {
				s.strikes = append(s.strikes, Strike{})
			}

			// Route second half of the strike to the newly created strike
			for i := 1; i < s.strikes[strikeID].Len-shift; i++ {
				rerouteCell := cell.Add(dir.Offset().ScaleUp(i))
				s.board[rerouteCell][dir.FixedID] = newStrikeID
			}

			s.strikes[newStrikeID].Start = cell.Add(dir.Offset())
			s.strikes[newStrikeID].Len = s.strikes[strikeID].Len - shift - 1
			s.strikes[newStrikeID].Dir = dir
			s.strikes[newStrikeID].ExtendableBefore = true
			s.strikes[newStrikeID].ExtendableAfter = s.strikes[strikeID].ExtendableAfter

			s.strikes[strikeID].Len = shift
			s.strikes[strikeID].ExtendableAfter = true
		}
	}

	// Make unoccupied
	delete(s.players, cell)

	return nil
}

func (s *StrikeSet) StrikesThrough(cell geom.Offset) [4]Strike {
	var strikes [4]Strike

	strikeRefs := s.board[cell]
	for i, strikeID := range strikeRefs {
		if strikeID == -1 {
			// strikes[i].Len will be 0
			continue
		}

		strikes[i] = s.strikes[strikeID]
	}

	return strikes
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
	return dir.X == other.X && dir.Y == other.Y && dir.FixedID == other.FixedID
}
