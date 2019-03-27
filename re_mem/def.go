package re_mem

import (
	"github.com/SamOrozco/re_mem/data"
	"github.com/SamOrozco/re_mem/query"
)

type Storage interface {
	GetCollection(name string) (Collection, error)
	RemoveCollection(name string) error
}

type Collection interface {
	Get(key string) (data.JsonMap, error)
	Create(doc interface{}) (string, error)
	Remove(key string) error
	Query(column, value string) ([]data.JsonMap, error)
	ExecuteQuery(query *query.SingleQuery) ([]data.JsonMap, error)
	NewStatement() query.Statement
	GetRowKeys(colName, stringValue string) []string
	GetRowsForKeys(keys []string) []data.JsonMap
}

type ColStore interface {
}
