package ai

import (
	"math/rand"
	//"sync"
	"log"
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

type BoardRank struct {
	P1, P2 playerMetrics
}

/*
func (a BoardRank) Victorious(player game.PlayerID, strikeLen int, canMoveNext bool) bool {
	ranks := a.P1
	if player == game.P2 {
		ranks = a.P2
	}

	return ranks.victorious(strikeLen, canMoveNext)
}
*/

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

/*
func (pm playerMetrics) victorious(strikeLen int, canMoveNext bool) bool {
	for i, rank := range pm.basis {
		if pm.count[i] == 0 {
			continue
		}

		if rank.Len >= strikeLen {
			return true
		}

		if !canMoveNext && rank.Len == strikeLen-1 && rank.Extensions == 1 {
			return true
		}

		if !canMoveNext && rank.Len == strikeLen-2 && rank.Extensions == 2 {
			return true
		}

		if canMoveNext && rank.Len == strikeLen-1 && rank.Extensions == 2 {
			return true
		}
	}

	return false
}
*/

func MetricTwoSideExtensible(player game.PlayerID, canMoveNext bool, old *BoardRank, candidate *BoardRank) bool {
	oldUs := old.P1
	oldThem := old.P2
	candUs := candidate.P1
	candThem := candidate.P2

	if player == game.P2 {
		/*
			oldUs, oldThem = oldThem, oldUs
			candUs, candThem = candThem, candUs
		*/
		oldUs = old.P2
		oldThem = old.P1
		candUs = candidate.P2
		candThem = candidate.P1
	}

	// TODO: where check victoriousness?

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

		/*
			if state.Over() {
				rank := computeRank(state)
				state.UndoLastMove()
				return rank, move
			}
		*/

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

	canMoveNext := depth%2 != 0

	bestOutcome := outcomes[0]
	// foundValid := false
	for i := range outcomes {
		/*
				if outcomes[i].Rank.Victorious(player.Other(), state.Conf.StrikeLength, depth%2 != 0) {
					continue
				}

			if !foundValid {
				bestOutcome = outcomes[i]
				foundValid = true
				continue
			}
		*/

		if p.cmp(player, canMoveNext, &bestOutcome.Rank, &outcomes[i].Rank) {
			bestOutcome = outcomes[i]
		}
	}

	return bestOutcome.Rank, bestOutcome.Cell
}

/*
func (p *AIPlayer) parallelMinimax(state *game.GameState, player game.PlayerID, depth int, threadCount int) (BoardRank, Offset) {
	moveCh := make(chan Offset)
	outcomeCh := make(chan moveOutcome)

	var wg sync.WaitGroup
	wg.Add(threadCount)
	for i := 0; i < threadCount; i++ {
		stateCopy := state.Clone()

		go func(state *game.GameState) {
			for move := range moveCh {
				state.MarkCell(move, player)

				if state.Over() {
					rank := computeRank(state)
					state.UndoLastMove()
					outcomeCh <- moveOutcome{move, rank}
					continue
				}

				var rank BoardRank
				if depth == 1 {
					// Calculate rank with move
					rank = computeRank(state)
				} else {
					rank, _ = p.minimax(state, player.Other(), depth-1)
				}

				outcomeCh <- moveOutcome{move, rank}

				state.UndoLastMove()
			}

			wg.Done()
		}(stateCopy)
	}

	go func() {
		for move := range state.Board.UnoccupiedCells() {
			moveCh <- move
		}
		close(moveCh)
	}()

	go func() {
		wg.Wait()
		close(outcomeCh)
	}()

	bestOutcome := <-outcomeCh
	foundValid := !bestOutcome.Rank.Victorious(player.Other(), state.Conf.StrikeLength, depth%2 != 0)
	for outcome := range outcomeCh {
		if outcome.Rank.Victorious(player.Other(), state.Conf.StrikeLength, depth%2 != 0) {
			continue
		}

		if !foundValid {
			bestOutcome = outcome
			foundValid = true
			continue
		}

		if p.cmp(player, &bestOutcome.Rank, &outcome.Rank) {
			bestOutcome = outcome
		}
	}

	return bestOutcome.Rank, bestOutcome.Cell
}
*/

func (p *AIPlayer) MakeMove(g *game.GameState) Offset {
	if p.gameCopy == nil {
		// TODO: illegal moves, border width 5 allows cells that 2 doesn't
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

	log.Printf("%v: level 1 cell count: %v", p.id, len(p.gameCopy.Board.UnoccupiedCells()))

	_, bestCell := p.minimax(p.gameCopy, p.id, p.SearchDepth)

	log.Printf("%v: chose move %v\n", p.id, bestCell)
	log.Printf("%v: rec depth  %v\n", p.id, p.recdepth)

	log.Println()

	return bestCell
}
