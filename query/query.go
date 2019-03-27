package query

import "github.com/SamOrozco/re_mem/hash"

// the query contract is
// that you must return a slice of zero or more
// keys.
// Once all keys are fetch and operator applied to the sets of keys
// AND and OR for example, we will load the proper rows via their keys
type Query interface {
	Get() []string
}

type Q struct {
	Column      string
	ValueHash   string
	CompareType CompareType
}

func NewQuery(col, rawValue string) *Q {
	// we currently only support one compareType and that's equal
	return &Q{
		Column:      col,
		ValueHash:   hash.NewHashString(rawValue),
		CompareType: Equal,
	}
}

type Predicate struct {
	Left     *Q
	Right    *Q
	Operator Op
}

func NewPredicate(left, right *Q, operator Op) *Predicate {
	return &Predicate{
		Left:     left,
		Right:    right,
		Operator: operator,
	}
}

type Statement struct {
	Left     *Predicate
	Right    *Predicate
	Operator Op
}
