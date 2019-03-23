package re_mem

import (
	"re-mem/data"
	"re-mem/files"
)

const rowDirName = ".row"
const colDirName = ".col"
const keyFileName = ".key"

type LocalCollection struct {
	rootDir               string
	collectionInitialized bool
}

func NewCollection(rootDir string) Collection {
	return &LocalCollection{rootDir: rootDir}
}

func (*LocalCollection) Get(key string) (Document, error) {
	panic("implement me")
}

func (col *LocalCollection) Create(document interface{}) (string, error) {
	err := col.initIfNeeded()
	if err != nil {
		return "", err
	}
	_, err = data.ParseToJsonMap(document)
	if err != nil {
		return "", err
	}
	//parseData(dataMap)
	return "", nil
}

func (*LocalCollection) Update(key string, document interface{}) (Document, error) {
	panic("implement me")
}

func (*LocalCollection) Query(column, value string) ([]Document, error) {
	panic("implement me")
}

func (*LocalCollection) Remove(key string) error {
	panic("implement me")
}

func (col LocalCollection) initIfNeeded() error {
	if !col.collectionInitialized {
		err := col.initCollection()
		if err != nil {
			return err
		}
	}
	return nil
}

func (col LocalCollection) initCollection() error {
	// init col dir
	if err := files.CreateDirIfNotExists(col.getColDir()); err != nil {
		return err
	}

	// init rowDir
	if err := files.CreateDirIfNotExists(col.getRowDir()); err != nil {
		return err
	}
	return nil
}

func (col LocalCollection) getRowDir() string {
	return col.rootDir + files.FileSep() + rowDirName

}
func (col LocalCollection) getColDir() string {
	return col.rootDir + files.FileSep() + colDirName
}
