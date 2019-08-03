package re

import (
	"github.com/SamOrozco/re_mem/data"
)

type Storage interface {
	// returns an existing or initializes a collection with the given name
	GetCollection(name string) (Collection, error)
	// removes a collection with the given name
	RemoveCollection(name string) error
}

type Collection interface {
	// get a json object for the given row key
	Get(key string) (data.JsonMap, error)
	// creates an entry in the collection from the interface passed to this method
	Create(doc interface{}) (string, error)
	// removes the object with the given row key
	Remove(key string) error
	// executes a simple query based on the column and the given value
	Query(column, value string) ([]data.JsonMap, error)
	// executes a query based on the given SingleQuery
	ExecuteQuery(query *SingleQuery) ([]data.JsonMap, error)
	// returns a new query statement for the given location
	NewStatement() Statement
	// gets all row keys where the records has the given value for the given column
	GetRowKeys(colName, stringValue string) []string
	// gets all rows for the given keys
	GetRowsForKeys(keys []string) []data.JsonMap
}

