package game

// CellState has a property of, when positive, being equal
// to a correct PlayerID that occupies it. Special meanings
// are negative. This way, CellState can be used for indexing
// slices related to players without any additional manipulation
type CellState int

const (
	CellUnavailable CellState = iota - 2
	CellUnoccupied
	CellP1
	CellP2
)

func (cs CellState) IsOccupiedBy(player PlayerID) bool {
	return cs >= 0 && PlayerID(cs) == player
}
