package query

type Query struct {
	Column      string
	ValueHash   string
	CompareType CompareType
}

func NewQuery(col, valHash string) *Query {
	// we currently only support one compareType and that's equal
	return &Query{
		Column:      col,
		ValueHash:   valHash,
		CompareType: Equal,
	}
}

type Statement struct {
	Left     *Query
	Right    *Query
	Operator Op
}

func NewStatement(left, right *Query, operator Op) *Statement {
	return &Statement{
		Left:     left,
		Right:    right,
		Operator: operator,
	}
}
