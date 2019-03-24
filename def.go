package main

import "github.com/SamOrozco/re_mem/data"

type Storage interface {
	GetCollection(name string) (Collection, error)
	RemoveCollection(name string) error
}

type Collection interface {
	Get(key string) (data.JsonMap, error)
	Create(doc interface{}) (string, error)
	Query(column, value string) ([]data.JsonMap, error)
	Remove(key string) error
}
