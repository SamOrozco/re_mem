package main

import (
	"fmt"
	"github.com/SamOrozco/re_mem/re"
	"testing"
)

var testLocation = "C:\\Users\\samue\\re_mem"

func TestInitData(t *testing.T) {
	store := re.NewLocalStorage("/Users/samorozco/first_db")
	usersCollection, err := store.GetCollection("users")
	if err != nil {
		panic(err)
	}

	names := []string{"abe", "lincoln", "steven", "hawking", "richard"}
	ages := []int{34, 231, 77, 12, 77}
	email := []string{"abe@g.com", "l@g.com", "g@g.com", "d@g.com", "t@g.com"}

	// create 5 users
	for i := 0; i < 5; i++ {
		user := &User{
			Name:  names[i],
			Age:   ages[i],
			Email: email[i],
		}
		recordKey, err := usersCollection.Create(user)
		if err != nil {
			panic(err)
		}
		println(fmt.Sprintf("record %s", recordKey))
	}

	// query for abe by email
	docs, err := usersCollection.Query("email", "abe@g.com")
	if err != nil {
		panic(err)
	}

	for _, doc := range docs {
		println(doc.String())
	}
}

func TestComplexQuery(test *testing.T) {
	store := re.NewLocalStorage("/Users/samorozco/first_db")
	col, err := store.GetCollection("users")
	if err != nil {
		panic(err)
	}

	stmt := col.NewStatement()
	query := stmt.NewQuery("name", "sam")
	docs := query.Fetch()

	for _, v := range docs {
		println(v.String())
	}
}

func TestWriteIndexCol(test *testing.T) {
	store := re.NewLocalStorage(testLocation)
	collection, err := store.GetCollection("users")
	if err != nil {
		panic(err)
	}

	err = collection.IndexColumn("name")
	if err != nil {
		panic(err)
	}
}
