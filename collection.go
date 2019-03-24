package main

import (
	"fmt"
	"github.com/SamOrozco/re_mem/data"
	"github.com/SamOrozco/re_mem/files"
	"github.com/SamOrozco/re_mem/hash"
	"github.com/SamOrozco/re_mem/query"
	"github.com/SamOrozco/re_mem/util"
	"strings"
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
	colLocation := col.getColValueLocation(column, value)
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

// this method turns around and calls the Query(col, val) method
func (col *LocalCollection) ExecuteQuery(query *query.Query) ([]data.JsonMap, error) {
	return col.Query(query.Column, query.ValueHash)
}

// The trick behind the execute statement method is to only deal with the keys of a result until you
// are ready to query data
func (col *LocalCollection) ExecuteStatement(stmt *query.Statement) ([]data.JsonMap, error) {
	keys, err := col.queryKeysForStatement(stmt)
	if err != nil {
		return nil, err
	}
	return col.fetchDocsForKeys(keys)
}

func (col *LocalCollection) Remove(key string) error {
	err := files.DeleteFile(col.getRowValueLocation(key))
	if err != nil {
		return err
	}
	return col.removeKey(key)
}

func (col *LocalCollection) fetchDocsForKeys(keys []string) ([]data.JsonMap, error) {
	// fetch values for keys
	docs, _, err := col.readDocumentsFromRowKeys(keys)
	if err != nil {
		return nil, err
	}
	return docs, nil
}

func (col *LocalCollection) queryKeysForStatement(stmt *query.Statement) ([]string, error) {
	left := stmt.Left
	right := stmt.Right
	// no queries return an empty return
	if right == nil && left == nil {
		return make([]string, 0), nil
	}
	// execute one if the other is nil
	if right == nil {
		return col.readKeysFromQuery(left)
	}
	if left == nil {
		return col.readKeysFromQuery(right)
	}

	leftKeys, err := col.readKeysFromQuery(left)
	if err != nil {
		return nil, err
	}

	rightKeys, err := col.readKeysFromQuery(right)

	return col.mergeKeys(leftKeys, rightKeys, stmt.Operator), nil
}

func (col *LocalCollection) mergeKeys(left, right []string, operator query.Op) []string {
	if operator == query.And {
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

func (col LocalCollection) readKeysFromQuery(query *query.Query) ([]string, error) {
	if len(query.Column) < 1 {
		return make([]string, 0), nil
	}
	loc := col.getColValueLocation(query.Column, query.ValueHash)
	return files.ReadLinesFromFile(loc)
}

func (col LocalCollection) readKeysFromColVal(colName, val string) ([]string, error) {
	loc := col.getColValueLocation(colName, val)
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
