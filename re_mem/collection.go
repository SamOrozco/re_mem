package re_mem

import (
	"re-mem/data"
	"re-mem/files"
	"re-mem/hash"
	"re-mem/util"
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
	jsonMap, err := data.ParseToJsonMap(document)
	if err != nil {
		return "", err
	}

	return col.createRecord(jsonMap)
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

// this method will iterate every key in the data map
// it will then write a col file for each col if it doesn't exist and
// and then try to hash each col by value
func (col LocalCollection) createRecord(data data.JsonMap) (string, error) {
	// this is currently only supporting flat objects
	recordKey := hash.NewRandomKey()
	for key, value := range data {
		if stringValue, ok := value.(string); ok {
			colName := util.CleanseName(key)
			hashedValue := hash.NewHashString(stringValue)
			err := col.writeRecord(colName, hashedValue, recordKey)
			if err != nil {
				return "", err
			}
		}
	}
	return recordKey, nil
}

func (col LocalCollection) writeRecord(colName, value, key string) error {
	// create dir if not exists
	if err := files.CreateDirIfNotExists(col.getColLocation(colName)); err != nil {
		return err
	}

	// does the file with this value NOT exist?
	if !files.Exists(col.getColValueLocation(colName, value)) {
		// create new file
		return col.writeNewColValue(colName, value, key)
	} else {
		// append existing file
		return col.appendNewColValue(colName, value, key)
	}
}

func (col LocalCollection) writeNewColValue(colName, value, key string) error {
	return files.WriteData(col.getColValueLocation(colName, value), key)
}

func (col LocalCollection) appendNewColValue(colName, value, key string) error {
	return files.AppendData(col.getColValueLocation(colName, value), key)
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

func (col LocalCollection) getColValueLocation(colName, hashedValue string) string {
	return col.getColDir() + files.FileSep() + colName + files.FileSep() + hashedValue
}

func (col LocalCollection) getColLocation(colName string) string {
	return col.getColDir() + files.FileSep() + colName
}

func (col LocalCollection) getRowDir() string {
	return col.rootDir + files.FileSep() + rowDirName

}
func (col LocalCollection) getColDir() string {
	return col.rootDir + files.FileSep() + colDirName
}
