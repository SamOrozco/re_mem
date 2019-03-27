package main

import (
	"fmt"
	"github.com/SamOrozco/re_mem/re_mem"
)

type User struct {
	Name  string
	Age   int
	Email string
}

const storeLocation = "/Users/samorozco/first_db"
const collectionName = "users"

func main() {
	queryStatement()
	//createUsers()
}

func queryStatement() {
	store := re_mem.NewLocalStorage(storeLocation)
	usersCollection, err := store.GetCollection(collectionName)
	if err != nil {
		panic(err)
	}
	stmt := usersCollection.NewStatement()

	nameAbe := stmt.NewQuery("name", "abe")
	nameSteven := stmt.NewQuery("email", "abe@g.com")
	both := stmt.NewQueryClause(nameAbe, nameSteven, re_mem.And)
	combo := stmt.NewQueryClause(both, stmt.NewQuery("name", "steven"), re_mem.Or)

	docs := combo.Fetch()
	for _, val := range docs {
		println(val.String())
	}

}

func createUsers() {
	store := re_mem.NewLocalStorage(storeLocation)
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
