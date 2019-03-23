package re_mem

type Document map[string]interface{}

type Storage interface {
	GetCollection(name string) (Collection, error)
	RemoveCollection(name string) error
}

type Collection interface {
	Get(key string) (Document, error)
	Create(doc interface{}) (string, error)
	Update(key string, doc interface{}) (Document, error)
	Query(column, value string) ([]Document, error)
	Remove(key string) error
}
