package game

import "fmt"

type PlayerID int

const (
	P1 PlayerID = iota
	P2
)

func (p PlayerID) Other() PlayerID {
	if p == P1 {
		return P2
	} else if p == P2 {
		return P1
	}

	panic(fmt.Sprintf("PlayerID: get other player: player is invalid (value=%d)", p))
}

func (p PlayerID) String() string {
	if p == P1 {
		return "P1"
	} else if p == P2 {
		return "P2"
	}

	return fmt.Sprintf("UnknownP%d", int(p))
}
