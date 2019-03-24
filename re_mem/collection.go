package re_mem

import (
	"fmt"
	"re-mem/data"
	"re-mem/files"
	"re-mem/hash"
	"re-mem/util"
)

const rowDirName = ".row"
const colDirName = ".col"
const keyFileName = ".key"
const id = "_id"

type LocalCollection struct {
	rootDir               string
	collectionInitialized bool
}

func NewCollection(rootDir string) Collection {
	return &LocalCollection{rootDir: rootDir}
}

func (*LocalCollection) Get(key string) (data.JsonMap, error) {
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

	rowkey, err := col.createColumnData(jsonMap)
	return rowkey, col.insertRowData(rowkey, jsonMap)
}

func (*LocalCollection) Update(key string, document interface{}) (data.JsonMap, error) {
	panic("implement me")
}

func (col *LocalCollection) Query(column, value string) ([]data.JsonMap, error) {
	// does file for given value exist
	colLocation := col.getColValueLocation(column, hash.NewHashString(value))
	if !files.Exists(colLocation) {
		return nil, nil
	}

	rowKeys, err := files.ReadLinesFromFile(colLocation)
	if err != nil {
		return nil, err
	}
	return col.readDocumentsFromRowKeys(rowKeys)
}

func (*LocalCollection) Remove(key string) error {
	panic("implement me")
}

func (col *LocalCollection) readDocumentsFromRowKeys(rows []string) ([]data.JsonMap, error) {
	result := make([]data.JsonMap, len(rows))
	resultIndex := 0
	for _, rowKey := range rows {
		contents, err := files.ReadDataFromFile(col.getRowValueLocation(rowKey))
		if err != nil {
			return nil, err
		}

		jsonMap, err := data.ParseJsonBytesToMap(contents)
		if err != nil {
			return nil, err
		}

		// add the row key to the jsonMap
		jsonMap[id] = rowKey
		result[resultIndex] = jsonMap
		resultIndex++
	}

	return result, nil
}

// this method will iterate every key in the data map
// it will then write a col file for each col if it doesn't exist and
// and then try to hash each col by value
func (col LocalCollection) createColumnData(data data.JsonMap) (string, error) {
	recordKey := hash.NewRandomKey()

	// we want to range each column and we are going to store only values that
	// can be asserted as a string
	// we will then hash the value are store the record key for each column
	// saying I this record have a value with this column
	for key, value := range data {
		stringValue := fmt.Sprintf("%v", value)
		colName := util.CleanseName(key)
		hashedValue := hash.NewHashString(stringValue)
		err := col.writeRecord(colName, hashedValue, recordKey)
		if err != nil {
			return "", err
		}
	}
	return recordKey, nil
}

func (col *LocalCollection) insertRowData(key string, data data.JsonMap) error {
	// add key to our keys file
	// the keys file allows us to iterate all rows in file
	// we read all keys in the file
	if err := col.appendKeysFile(key); err != nil {
		return err
	}

	// write row file data
	if err := files.WriteNewData(col.getRowValueLocation(key), data.String()); err != nil {
		return err
	}
	return nil

}

func (col *LocalCollection) appendKeysFile(key string) error {
	return files.WriteData(col.getKeyFileLocation(), key)
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
	return files.WriteNewData(col.getColValueLocation(colName, value), key)
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

func (col LocalCollection) getRowValueLocation(rowKey string) string {
	return col.getRowDir() + files.FileSep() + rowKey
}

func (col LocalCollection) getKeyFileLocation() string {
	return col.getRowDir() + files.FileSep() + keyFileName
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
