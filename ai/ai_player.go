package ai

import (
	//"log"
	"math/rand"
	"time"

	"github.com/kitsunemikan/six-purrpurrs/game"
	. "github.com/kitsunemikan/six-purrpurrs/geom"
)

// TODO: test if accounting for free-standing 1-lengths does anything
type RankMetric struct {
	Len        int
	Extensions int
}

type IsRankBetterForPlayerFunc func(p game.PlayerID, canMoveNext bool, old *BoardRank, candidate *BoardRank) bool

// defaultMetricBasis represents an ordered set of basis vectors in terms of which
// the player metrics are represented. Each basis represents a length of a strike
// and the number sides from which it can be extended. The order of the metric basis
// allows us to establesh a less-than ordering on the player metric vectors so that
// we can compare them.
var defaultMetricBasis = []RankMetric{
	{8, 2},
	{8, 1},
	{8, 0},

	{7, 2},
	{7, 1},
	{7, 0},

	{6, 2},
	{6, 1},
	{6, 0},

	{5, 2},
	{5, 1},

	{4, 2},
	{3, 2},

	{4, 1},
	{3, 1},

	{2, 2},
	{2, 1},

	{1, 2},
	{1, 1},
}

// BoardRank is represented by a pair of player metric vectors in the default metric basis
type BoardRank struct {
	P1, P2 playerMetrics
}

type moveOutcome struct {
	Cell Offset
	Rank BoardRank
}

func computeRank(s *game.GameState) BoardRank {
	rank := BoardRank{
		P1: newPlayerMetrics(defaultMetricBasis),
		P2: newPlayerMetrics(defaultMetricBasis),
	}

	strikes := s.StrikeStat.Strikes()
	for _, strike := range strikes {
		var metric RankMetric
		metric.Len = strike.Len

		if strike.ExtendableBefore {
			metric.Extensions++
		}

		if strike.ExtendableAfter {
			metric.Extensions++
		}

		switch strike.Player {
		case game.P1:
			rank.P1.Add(metric, 1)

		case game.P2:
			rank.P2.Add(metric, 1)

		default:
			panic("unknown player")
		}
	}

	return rank
}

// Represents the strike statistics for a player expressed as a vector
// in the span of specified basis.
type playerMetrics struct {
	basis []RankMetric
	count []int
}

func newPlayerMetrics(basis []RankMetric) playerMetrics {
	return playerMetrics{
		basis: basis,
		count: make([]int, len(basis)),
	}
}

func (pm playerMetrics) Add(rank RankMetric, count int) {
	for i := range pm.basis {
		if pm.basis[i] != rank {
			continue
		}

		pm.count[i] += count
		return
	}
}

func (a playerMetrics) subtract(b playerMetrics) {
	if &a.basis[0] != &b.basis[0] {
		panic("ai: subtract metric vectors: different bases")
	}

	for i := range a.count {
		a.count[i] -= b.count[i]
	}
}

func (a playerMetrics) lessThan(b playerMetrics) bool {
	if &a.basis[0] != &b.basis[0] {
		panic("ai: subtract metric vectors: dimensionality differs")
	}

	for i := range a.count {
		if a.count[i] < b.count[i] {
			return true
		} else if a.count[i] > b.count[i] {
			return false
		}
	}

	return false
}

// The more longer strikes a player has the better,
// and the more two-side extensible strikes a player has the better.
// How we balance between these to heuristics is defined by the order of basis vectors in the
// defaultMetricBasis variable. Additionally, depending of whether the current player can move or not
// different metrics may be considered. For example if a player has 4-len 2-ext strike and can make a move
// then it's a guaranteed victory, but if cannot, the other player blocks the strike and turns it into 1-side
// extensible, robbing it of victory.
func MetricTwoSideExtensible(player game.PlayerID, canMoveNext bool, old *BoardRank, candidate *BoardRank) bool {
	oldUs := old.P1
	oldThem := old.P2
	candUs := candidate.P1
	candThem := candidate.P2

	if player == game.P2 {
		oldUs = old.P2
		oldThem = old.P1
		candUs = candidate.P2
		candThem = candidate.P1
	}

	// TODO: where check victoriousness? If AI can't look far into the future, it
	// can use heuristics to recognize some board configurations that will lead to
	// the opponents victory 100% if certain moves aren't made. For example,
	// the 4-len 2-ext example, or 5-len 1-ext examples....

	// NOTE: we ignore canMoveNext, as the AI performs well enough

	oldUs.subtract(oldThem)
	candUs.subtract(candThem)

	return oldUs.lessThan(candUs)
}

type AIPlayer struct {
	id   game.PlayerID
	rand *rand.Rand
	cmp  IsRankBetterForPlayerFunc

	recdepth int

	SearchDepth int

	gameCopy *game.GameState
}

func NewDefaultAIPlayer(id game.PlayerID) *AIPlayer {
	return &AIPlayer{
		id:   id,
		rand: rand.New(rand.NewSource(time.Now().Unix())),
		cmp:  MetricTwoSideExtensible,
	}
}

func (p *AIPlayer) minimax(state *game.GameState, player game.PlayerID, depth int) (BoardRank, Offset) {
	outcomes := make([]moveOutcome, 0, len(state.Board.UnoccupiedCells()))
	for move := range state.Board.UnoccupiedCells() {
		state.MarkCell(move, player)

		p.recdepth++

		var rank BoardRank
		if depth == 1 {
			// Calculate rank with move
			rank = computeRank(state)
		} else {
			rank, _ = p.minimax(state, player.Other(), depth-1)
		}

		outcomes = append(outcomes, moveOutcome{move, rank})

		state.UndoLastMove()
	}

	// CanMoveNext tells us whether our current player can make a move
	canMoveNext := depth%2 != 0

	bestOutcome := outcomes[0]
	for i := range outcomes {
		if p.cmp(player, canMoveNext, &bestOutcome.Rank, &outcomes[i].Rank) {
			bestOutcome = outcomes[i]
		}
	}

	return bestOutcome.Rank, bestOutcome.Cell
}

func (p *AIPlayer) MakeMove(g *game.GameState) Offset {
	if p.gameCopy == nil {
		// NOTE: border radius of 2 is the smallest playeble width, hence by using it
		// we can skip checking if AI moves are within or not the board.
		// Additionally, this game copy with all of its unoccupied cells will represent
		// all the cells that the AI should consider. By limiting it to 2, we greatly limit
		// the size of the search space for the performance's sake.
		// Alas, if the border radius is increased, the AI player will be able to play more
		// optimally, althogh it's up for a debate whether it's a good idea, as the game
		// may as well never end if players play optimally (mathematicians couldn't prove it)
		p.gameCopy = game.NewGame(game.GameOptions{
			Border:  2,
			Victory: g.VictoryChecker().Clone(),
		})
	}

	history := g.MoveHistoryCopy()
	for p.gameCopy.MoveNumber() < g.MoveNumber() {
		move := history[p.gameCopy.MoveNumber()-1]
		p.gameCopy.MarkCell(move.Cell, move.Player)
	}

	p.recdepth = 0

	// log.Printf("%v: level 1 cell count: %v", p.id, len(p.gameCopy.Board.UnoccupiedCells()))

	_, bestCell := p.minimax(p.gameCopy, p.id, p.SearchDepth)

	// log.Printf("%v: chose move %v\n", p.id, bestCell)
	// log.Printf("%v: rec depth  %v\n", p.id, p.recdepth)

	// log.Println()

	return bestCell
}
