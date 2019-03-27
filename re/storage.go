package re

import (
	"github.com/SamOrozco/re_mem/files"
	"github.com/SamOrozco/re_mem/hash"
)

type LocalStorage struct {
	rootDir string
}

// root dir is the directory re-mem will write and read it's objects
func NewLocalStorage(rootDir string) Storage {
	// confirm users disired root dir exists
	if err := files.CreateDirIfNotExists(rootDir); err != nil {
		panic(err)
	}

	// at this point we know we have initialized properly
	return &LocalStorage{rootDir: rootDir}
}

// the get collection method tries to find a collection with the given name
// if the collection exists it will return that else it will create a new collection
func (store LocalStorage) GetCollection(name string) (Collection, error) {
	// collection hash
	collectionHash := hash.NewHashString(name)
	collectionDir := store.getCollectionDir(collectionHash)
	// create collection dir if it doesn't exists
	if err := files.CreateDirIfNotExists(collectionDir); err != nil {
		return nil, err
	}
	return NewCollection(collectionDir), nil
}

func (LocalStorage) RemoveCollection(name string) error {
	panic("implement me")
}
func (store LocalStorage) getCollectionDir(hash string) string {
	return store.rootDir + files.FileSep() + hash
}
