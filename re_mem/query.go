package re_mem

import (
	"github.com/SamOrozco/re_mem/data"
)

// the query contract is
// that you must return a slice of zero or more
// keys.
// Once all keys are fetch and operator applied to the sets of keys
// AND and OR for example, we will load the proper rows via their keys
type Query interface {
	get() []string
	Fetch() []data.JsonMap
}

type SingleQuery struct {
	Column      string
	Value       string
	CompareType CompareType
	collection  Collection
}

func (single SingleQuery) Fetch() []data.JsonMap {
	keys := single.get()
	return single.collection.GetRowsForKeys(keys)
}

func (single SingleQuery) get() []string {
	return single.collection.GetRowKeys(single.Column, single.Value)
}

type Clause struct {
	left       Query
	right      Query
	operator   Op
	collection Collection
}

func (cl Clause) get() []string {
	return mergeKeys(cl.left.get(), cl.right.get(), cl.operator)
}

func (cl Clause) Fetch() []data.JsonMap {
	keys := cl.get()
	return cl.collection.GetRowsForKeys(keys)
}

type Statement struct {
	Collection Collection
}

func (stmt Statement) NewQuery(colName, stringValue string) Query {
	return &SingleQuery{
		Column:      colName,
		Value:       stringValue,
		CompareType: Equal,
		collection:  stmt.Collection,
	}
}

func (stmt Statement) NewQueryClause(left, right Query, operator Op) Query {
	return &Clause{
		left:       left,
		right:      right,
		operator:   operator,
		collection: stmt.Collection,
	}
}

func mergeKeys(left, right []string, operator Op) []string {
	if operator == And {
		// merge
		var iter []string
		var mp data.LookupMap
		rightLen := len(right)
		leftLen := len(left)
		if rightLen < leftLen {
			iter = right
			mp = data.StringsToLookupMap(left)
		} else {
			iter = left
			mp = data.StringsToLookupMap(right)
		}

		// iterate list and and check if value to in list and map
		result := make([]string, 0)
		for _, val := range iter {
			_, ok := mp[val]
			// if val is in map and iter
			// add to the result
			if ok {
				result = append(result, val)
			}
		}
		return result
	} else {
		// we need to add unique keys from both sides
		mp := make(data.LookupMap, 0)
		for _, leftVal := range left {
			mp[leftVal] = true
		}

		for _, rightVal := range right {
			mp[rightVal] = true
		}

		// put them into a map to keep them unique
		mpLen := len(mp)
		if mpLen < 1 {
			return make([]string, 0)
		}
		// iterator map keys and put into result
		result := make([]string, mpLen)
		idx := 0
		for k := range mp {
			result[idx] = k
			idx++
		}
		return result
	}
}
