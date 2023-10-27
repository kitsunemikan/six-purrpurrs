package game

import (
	"github.com/kitsunemikan/six-purrpurrs/geom"
)

type VictoryChecker interface {
	Reset()
	Clone() VictoryChecker
	StrikeLength() int
	CheckAt(strikes *StrikeSet, pos geom.Offset) bool
	CandidatesAroundFor(strikes *StrikeSet, pos geom.Offset, player PlayerID) []geom.Offset
	Reached() bool
	VictoriousStrike() []geom.Offset
	VictoriousPlayer() PlayerID
}
