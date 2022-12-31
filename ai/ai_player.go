package ai

import (
	"math/rand"
	"time"

	"github.com/kitsunemikan/six-purrpurrs/game"
	. "github.com/kitsunemikan/six-purrpurrs/geom"
)

// TODO: test if accounting for free-standing 1-lengths does anything
type RankMetric struct {
	Length     int
	Extensions int
}

type RankLessThanForP1Func func(a *BoardRank, b *BoardRank) bool

type BoardRank struct {
	P1 map[RankMetric]int
	P2 map[RankMetric]int
}

func (a BoardRank) IsWorseThan(b BoardRank, player game.PlayerID) bool {
	return true
}

type moveOutcome struct {
	Cell Offset
	Rank BoardRank
}

type AIPlayer struct {
	id   game.PlayerID
	rand *rand.Rand
	cmp  RankLessThanForP1Func
}

func NewDefaultAIPlayer(id game.PlayerID) game.PlayerAgent {
	return &AIPlayer{
		id:   id,
		rand: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (p *AIPlayer) minimax(g *game.BoardState, player game.PlayerID, depth int) (BoardRank, Offset) {
	outcomes := make([]moveOutcome, 0, len(g.UnoccupiedCells()))
	for move := range g.UnoccupiedCells() {
		g.MarkCell(move, player)

		var rank BoardRank
		if depth == 1 {
			// Calculate rank with move
			rank = computeRank(g)
		} else {
			rank, _ = p.minimax(g, player.Other(), depth-1)
		}

		outcomes = append(outcomes, moveOutcome{move, rank})

		g.UndoLastMove()
	}

	bestOutcome := outcomes[0]
	for i := 1; i < len(outcomes); i++ {
		if p.cmp(&bestOutcome.Rank, &outcomes[i].Rank) {
			bestOutcome = outcomes[i]
		}
	}

	return bestOutcome.Rank, bestOutcome.Cell
}

func (p *AIPlayer) MakeMove(b *game.BoardState) Offset {
	board := make(map[Offset]game.CellState, len(b.AllCells()))
	for k, v := range b.AllCells() {
		board[k] = v
	}

	boardCopy := b.Clone()

	_, bestCell := p.minimax(boardCopy, p.id, 4)

	return bestCell
}

func computeRank(g *game.BoardState) BoardRank {
	return BoardRank{}
}
