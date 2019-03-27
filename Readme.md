# re-mem

Re-mem is a local disc document store meant for 
simple Creates, and Fetches 


Every operation is reading from or writing to disc. 

Your "re-mem" database is a directory containing a few files and directories.
Just know that it is more than just one data file.

 


## Getting started


### install
```bash
go get -u github.com/SamOrozco/re_mem
```


### Usage
Start by selecting the directory you would like your database to live in.
```go
store := re_mem.NewLocalStorage("/Users/samorozco/first_db")
```

Next we must initialize a collection of objects that we want to use. 
 
Everything you create, fetch, or delete will be a part of a collection

```go
// if the collection exists it will return the existing collection
// else it will create a new collection
userCollection, err := store.GetCollection("users")
```

It does not matter the structure of the object you create.
In fact you can have any number of totally different object structures in one collection

```go
userKey, err := usersCollection.Create(&User{
		Name:  "re-mem",
		Age:   0,
		Email: "re-mem@gmail.com",
	})
println(userKey)


companyKey, err := usersCollection.Create(&Company{
		Phone: "5556667777"
		Email: "re-mem@gmail.com",
	})
println(companykey)
```

### Complex Queries
The addition of the collection statement has added abilities for composing
different queries together. To form very complex queries that wait to load any data from the disc until they have to. 

```go
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

const storeLocation = "/Users/samorozco/first_db"
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


	
```
  




## Struct Example 

```go
package main

import (
	"fmt"
	"re-mem/re_mem"
)

type User struct {
	Name  string
	Age   int
	Email string
}

func main() {
	store := re_mem.NewLocalStorage("/Users/samorozco/first_db")
	usersCollection, err := store.GetCollection("users")
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

	// query for abe by email
	docs, err := usersCollection.Query("email", "abe@g.com")
	if err != nil {
		panic(err)
	}

	for _, doc := range docs {
		println(doc.String())
	}

	// Get by key
	doc, err := usersCollection.Get("<some_record_key>")
	if err != nil {
		panic(err)
	}
	println(doc)

	// remove doc

	err = usersCollection.Remove("<some_record_key>")
	if err != nil {
		panic(err)
	}

}


``` 


## Json Example

```go
package main

import (
	"encoding/json"
	"re-mem/re_mem"
)

const jsonData = `{
  "data": [
    {
      "CountryName": "United States",
      "Description": "Texas",
      "StateCode": "TX"
    },
    {
      "CountryName": "United States",
      "Description": "United States Minor Outlying Islands (see also separate entry under UM)",
      "StateCode": "UM"
    },
    {
      "CountryName": "United States",
      "Description": "Utah",
      "StateCode": "UT"
    },
    {
      "CountryName": "United States",
      "Description": "Virginia",
      "StateCode": "VA"
    },
    {
      "CountryName": "United States",
      "Description": "Virgin Islands, U.S. (see also separate entry under VI)",
      "StateCode": "VI"
    },
    {
      "CountryName": "United States",
      "Description": "Vermont",
      "StateCode": "VT"
    },
    {
      "CountryName": "United States",
      "Description": "Washington",
      "StateCode": "WA"
    },
    {
      "CountryName": "United States",
      "Description": "Wisconsin",
      "StateCode": "WI"
    },
    {
      "CountryName": "United States",
      "Description": "West Virginia",
      "StateCode": "WV"
    },
    {
      "CountryName": "United States",
      "Description": "Wyoming",
      "StateCode": "WY"
    }
  ]
}`

func main() {
	store := re_mem.NewLocalStorage("/Users/samorozco/first_db")
	stateCollection, err := store.GetCollection("states")
	if err != nil {
		panic(err)
	}

	// parse json into map
	jsonMap := make(map[string]interface{}, 0)
	err = json.Unmarshal([]byte(jsonData), &jsonMap)
	if err != nil {
		panic(err)
	}

	// extact data field
	dataField := jsonMap["data"]

	records := dataField.([]interface{})
	for _, rec := range records {
		key, err := stateCollection.Create(rec)
		if err != nil {
			panic(err)
		}
		println(key)
	}
}

```
