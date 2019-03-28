package main

import (
	"github.com/SamOrozco/re_mem/files"
	"github.com/SamOrozco/re_mem/re"
	"os"
	"testing"
)

func getHomeDir() string {
	return os.Getenv("HOME")
}

func TestCreateDirOnInit(test *testing.T) {
	homeDir := getHomeDir()
	testDir := homeDir + "/test_db"
	_ = re.NewLocalStorage(testDir)
	if !files.Exists(testDir) {
		test.Fail()
	}
	// remove dir
	err := files.DeleteDir(testDir)
	if err != nil {
		test.Fatal(err)
	}

	//confirm removed
	if files.Exists(testDir) {
		test.Fail()
	}
}

func TestCreateRecordAndFetchRecord(test *testing.T) {
	homeDir := getHomeDir()
	storeDir := homeDir + "/test_db"
	store := re.NewLocalStorage(storeDir)
	carColl, err := store.GetCollection("cars")
	if err != nil {
		test.Fatal(err)
	}

	// we have our collection now lets create and read

	origDoc := map[string]interface{}{
		"brand": "chevy",
		"year":  "2001",
		"color": "red",
	}
	key, err := carColl.Create(origDoc)
	doc, err := carColl.Get(key)
	delete(doc, "_id")
	if err != nil {
		test.Fatal(err)
	}

	if doc["brand"] != origDoc["brand"] {
		test.Fail()
	}

	if doc["year"].(string) != origDoc["year"].(string) {
		test.Fail()
	}

	if doc["color"] != origDoc["color"] {
		test.Fail()
	}
}
