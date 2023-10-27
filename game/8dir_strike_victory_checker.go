package game

import "github.com/kitsunemikan/six-purrpurrs/geom"

type EightDirStrikeVictoryChecker struct {
	VictoryLength int

	strike []geom.Offset
	player PlayerID
}

func (ch *EightDirStrikeVictoryChecker) StrikeLength() int {
	return ch.VictoryLength
}

func (ch *EightDirStrikeVictoryChecker) CheckAt(strikes *StrikeSet, pos geom.Offset) bool {
	cellStrikes := strikes.StrikesThrough(pos)

	for strikeID := range cellStrikes {
		if cellStrikes[strikeID].Len >= ch.VictoryLength {
			ch.strike = cellStrikes[strikeID].AsCells()
			ch.player = cellStrikes[strikeID].Player
			return true
		}
	}

	return false
}

func (ch *EightDirStrikeVictoryChecker) CandidatesAroundFor(strikes *StrikeSet, pos geom.Offset, player PlayerID) []geom.Offset {
	var candidates []geom.Offset

	for _, dir := range StrikeDirs {
		// Forward direction
		afterCell := pos.Add(dir.Offset())
		afterStrike := strikes.StrikesThrough(afterCell)[dir.FixedID]
		if afterStrike.Player == player {
			cells := afterStrike.AsCells()
			candidates = append(candidates, cells...)
		}

		// Backward direction
		beforeCell := pos.Sub(dir.Offset())
		beforeStrike := strikes.StrikesThrough(beforeCell)[dir.FixedID]
		if beforeStrike.Player == player {
			cells := beforeStrike.AsCells()
			candidates = append(candidates, cells...)
		}
	}

	return candidates
}

func (ch *EightDirStrikeVictoryChecker) Clone() VictoryChecker {
	strikeCopy := make([]geom.Offset, 0, len(ch.strike))
	strikeCopy = append(strikeCopy, ch.strike...)

	return &EightDirStrikeVictoryChecker{
		VictoryLength: ch.VictoryLength,

		strike: strikeCopy,
		player: ch.player,
	}
}

func (ch *EightDirStrikeVictoryChecker) Reset() {
	ch.strike = nil
	ch.player = P1
}

func (ch *EightDirStrikeVictoryChecker) Reached() bool {
	return ch.strike != nil
}

func (ch *EightDirStrikeVictoryChecker) VictoriousStrike() []geom.Offset {
	return ch.strike
}

func (ch *EightDirStrikeVictoryChecker) VictoriousPlayer() PlayerID {
	return ch.player
}
