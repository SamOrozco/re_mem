package main

import (
	"fmt"
	"github.com/SamOrozco/re_mem/re"
	"github.com/davidbanham/human_duration"
	"math/rand"
	"testing"
	"time"
)

const storeLocation1 = "/Users/samorozco/performance"
const collectionName1 = "users"

func TestInsertPerformance(test *testing.T) {
	store := re.NewLocalStorage(storeLocation1)
	usersCollection, err := store.GetCollection(collectionName1)
	if err != nil {
		panic(err)
	}

	names := []string{"abe", "lincoln", "steven", "hawking", "richard", "this", "that", "him", "her"}
	ages := []int{34, 231, 77, 12, 77}
	email := []string{"abe@g.com", "l@g.com", "g@g.com", "d@g.com", "t@g.com"}

	start := time.Now()
	// create 5 users
	for i := 0; i < 100000; i++ {
		user := &User{
			Name:  names[rand.Int()%9],
			Age:   ages[rand.Int()%5],
			Email: email[rand.Int()%5],
		}
		_, err := usersCollection.Create(user)
		if err != nil {
			panic(err)
		}
	}

	duration := time.Since(start)
	println(human_duration.String(duration, "second"))
}

func TestQueryPerformance(test *testing.T) {
	store := re.NewLocalStorage(storeLocation1)
	usersCollection, err := store.GetCollection(collectionName1)
	if err != nil {
		panic(err)
	}

	stmt := usersCollection.NewStatement()

	// single query
	start := time.Now()
	query := stmt.NewQuery("name", "abe")
	docs := query.Fetch()
	size := len(docs)
	duration := time.Since(start)
	println(fmt.Sprintf("read %d docs", size))
	println(fmt.Sprintf("duration : %s", duration))

	// two clause query
	start = time.Now()
	nameHawk := stmt.NewQuery("name", "hawking")
	age77 := stmt.NewQuery("age", "77")
	nameAndAge := stmt.NewQueryClause(nameHawk, age77, re.Or)
	docs = nameAndAge.Fetch()
	size = len(docs)
	duration = time.Since(start)
	println(fmt.Sprintf("read %d docs", size))
	println(fmt.Sprintf("duration : %s", duration))

	// complex query
	start = time.Now()
	nameHawk = stmt.NewQuery("name", "hawking")
	age77 = stmt.NewQuery("age", "77")
	nameAndAge = stmt.NewQueryClause(nameHawk, age77, re.Or)
	andEmailAbe := stmt.NewQueryClause(nameAndAge, stmt.NewQuery("email", "abe@g.com"), re.And)
	println(andEmailAbe.Explain())
	docs = andEmailAbe.Fetch()
	size = len(docs)
	duration = time.Since(start)
	println(fmt.Sprintf("read %d docs", size))
	println(fmt.Sprintf("duration : %s", duration))

	for _, doc := range docs {
		println(doc.String())
	}

}

//read 1139 docs
//duration : 72.89515ms
//read 4686 docs
//duration : 221.073691ms
//read 936 docs
//duration : 42.355359ms

// 100000 docs
//read 11214 docs
//duration : 932.862374ms
//read 46669 docs
//duration : 4.463323699s
//read 9337 docs
//duration : 379.915291ms

//read 11214 docs
//duration : 1.580013132s
//read 46669 docs
//duration : 6.726733346s
//read 9337 docs
//duration : 393.000881ms
