package main

import (
	"fmt"
	"github.com/SamOrozco/re_mem/re"
)

type User struct {
	Name  string
	Age   int
	Email string
}

const storeLocation = "C:\\Users\\samue\\test_location"
const collectionName = "users"

func main() {
	createUsers()
	queryStatement()
}

func queryStatement() {
	store := re.NewLocalStorage(storeLocation)
	usersCollection, err := store.GetCollection(collectionName)
	if err != nil {
		panic(err)
	}
	// get a "Query Statement" from the collection
	// the query statement allows you to build complex queries
	stmt := usersCollection.NewStatement()

	// simple single value query
	nameAbe := stmt.NewQuery("name", "abe")
	// another simple query
	nameSteven := stmt.NewQuery("email", "abe@g.com")
	// combine those two with an and operator
	both := stmt.NewQueryClause(nameAbe, nameSteven, re.And)
	// combine that combo with another query and the OR operator
	combo := stmt.NewQueryClause(both, stmt.NewQuery("name", "steven"), re.Or)
	// fetch docs for all those queries
	docs := combo.Fetch()
	for _, val := range docs {
		println(val.String())
	}

}

func createUsers() {
	store := re.NewLocalStorage(storeLocation)
	usersCollection, err := store.GetCollection(collectionName)
	if err != nil {
		panic(err)
	}

	names := []string{"abe", "lincoln", "steven", "hawking", "richard"}
	ages := []int{34, 231, 121, 12, 77}
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
}
