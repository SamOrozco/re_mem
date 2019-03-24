package query

import "github.com/SamOrozco/re_mem/hash"

type Query struct {
	Column      string
	ValueHash   string
	CompareType CompareType
}

func NewQuery(col, rawValue string) *Query {
	// we currently only support one compareType and that's equal
	return &Query{
		Column:      col,
		ValueHash:   hash.NewHashString(rawValue),
		CompareType: Equal,
	}
}

type Predicate struct {
	Left     *Query
	Right    *Query
	Operator Op
}

func NewPredicate(left, right *Query, operator Op) *Predicate {
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
