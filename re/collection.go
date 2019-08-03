package re

import (
	"fmt"
	"github.com/SamOrozco/re_mem/data"
	"github.com/SamOrozco/re_mem/files"
	"github.com/SamOrozco/re_mem/hash"
	"github.com/SamOrozco/re_mem/util"
	"strings"
)

const rowDirName = ".row"
const colDirName = ".col"
const keyFileName = ".key"
const id = "_id"
const indexDirName = ".idx"
const indexFileName = ".indexes"

type LocalCollection struct {
	rootDir               string
	collectionInitialized bool
}

func NewCollection(rootDir string) Collection {
	return &LocalCollection{rootDir: rootDir}
}

// add a column to the indexed column directory
func (col LocalCollection) IndexColumn(columnName string) error {
	// lowercase column name
	columnName = strings.ToLower(columnName)
	// init index dir if needed
	if err := col.initIndexIfNeeded(); err != nil {
		return err
	}

	// write column to index file so we know that this column indexed
	if err := col.writeIndexColToFile(columnName); err != nil {
		return err
	}

	// add column to index dir for index data to be written to
	if err := col.addColIndexDir(columnName); err != nil {
		return err
	}
	return nil
}

// add a directory for this columns index
func (col LocalCollection) addColIndexDir(columnName string) error {
	colIndexDir := col.getColIndexDir(columnName)
	if err := files.CreateDirIfNotExists(colIndexDir); err != nil {
		return err
	}
	return nil
}

func (col *LocalCollection) Get(key string) (data.JsonMap, error) {
	return files.ReadJsonMapFromFile(col.getRowValueLocation(key))
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

func (col *LocalCollection) Query(column, value string) ([]data.JsonMap, error) {
	colLocation := col.getColValueLocation(column, hash.NewHashString(value))
	keys, err := col.readKeysFromLocation(colLocation)
	if err != nil {
		return nil, err
	}

	// fetch values for keys
	docs, keysString, err := col.readDocumentsFromRowKeys(keys)
	if err != nil {
		return nil, err
	}

	// we have to write new keys to file incase a record was deleted
	err = files.WriteNewData(colLocation, keysString)
	if err != nil {
		return nil, err
	}

	return docs, nil

}

func (col *LocalCollection) NewStatement() Statement {
	return Statement{Collection: col}
}

// this method turns around and calls the SingleQuery(col, val) method
func (col *LocalCollection) ExecuteQuery(query *SingleQuery) ([]data.JsonMap, error) {
	return col.Query(query.Column, query.Value)
}

func (col *LocalCollection) Remove(key string) error {
	err := files.DeleteFile(col.getRowValueLocation(key))
	if err != nil {
		return err
	}
	return col.removeKey(key)
}

// this method gets all rows keys with the given column name and string value
// in this collection
func (col *LocalCollection) GetRowKeys(columnName, stringValue string) []string {
	rows, err := col.readKeysFromColVal(columnName, stringValue)
	if err != nil {
		panic(err)
	}
	return rows
}

// this method reads all docs from the collection for the given row keys
func (col *LocalCollection) GetRowsForKeys(keys []string) []data.JsonMap {
	dat, err := col.fetchDocsForKeys(keys)
	if err != nil {
		panic(err)
	}
	return dat
}

func (col *LocalCollection) fetchDocsForKeys(keys []string) ([]data.JsonMap, error) {
	// fetch values for keys
	docs, _, err := col.readDocumentsFromRowKeys(keys)
	if err != nil {
		return nil, err
	}
	return docs, nil
}

func (col LocalCollection) removeKey(key string) error {
	keys, err := files.ReadLinesFromFile(col.getKeyFileLocation())
	if err != nil {
		return err
	}

	bldr := strings.Builder{}
	for _, curKey := range keys {
		if curKey != key {
			bldr.WriteString(fmt.Sprintf("%s \n", curKey))
		}
	}

	return files.WriteNewData(col.getKeyFileLocation(), bldr.String())
}

func (col LocalCollection) readKeysFromQuery(query *SingleQuery) ([]string, error) {
	if len(query.Column) < 1 {
		return make([]string, 0), nil
	}
	loc := col.getColValueLocation(query.Column, query.Value)
	return files.ReadLinesFromFile(loc)
}

func (col LocalCollection) readKeysFromColVal(colName, val string) ([]string, error) {
	loc := col.getColValueLocation(colName, hash.NewHashString(val))
	return files.ReadLinesFromFile(loc)
}

func (col LocalCollection) readKeysFromLocation(loca string) ([]string, error) {
	return files.ReadLinesFromFile(loca)
}

func (col *LocalCollection) readDocumentsFromRowKeys(rows []string) ([]data.JsonMap, string, error) {
	result := make([]data.JsonMap, 0)
	keyBldr := strings.Builder{}
	for _, rowKey := range rows {
		rowLoc := col.getRowValueLocation(rowKey)
		jsonMap, err := files.ReadJsonMapFromFile(rowLoc)
		if err != nil {
			// THIS WILL CHANGE
			// if we get an error we assume it is only because the file does not exists
			// TODO fix this
			continue
		}
		keyBldr.WriteString(fmt.Sprintf("%s \n", rowKey))
		result = append(result, jsonMap)
	}

	return result, keyBldr.String(), nil
}

// this method will iterate every key in the data map
// it will then write a col file for each col if it doesn't exist and
// and then try to hash each col by value
func (col LocalCollection) createColumnData(data data.JsonMap) (string, error) {
	recordKey := hash.NewRandomKey()

	// add id to record
	data[id] = recordKey
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

func (col LocalCollection) initIndexIfNeeded() error {
	if err := files.CreateDirIfNotExists(col.getIndexDir()); err != nil {
		return err
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

func (col LocalCollection) writeIndexColToFile(columnName string) error {
	return files.WriteLine(col.getIndexFileLocation(), columnName)
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

func (col LocalCollection) getIndexDir() string {
	return col.rootDir + files.FileSep() + colDirName + files.FileSep() + indexDirName
}

func (col LocalCollection) getColIndexDir(columnName string) string {
	return col.getIndexDir() + files.FileSep() + columnName
}

func (col LocalCollection) getIndexFileLocation() string {
	return col.getIndexDir() + files.FileSep() + indexFileName
}
