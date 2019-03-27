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


## Performance


Inserts Users: 
```go
type User struct {
	Name  string
	Age   int
	Email string
}
```
```
Time of inserts with random data
1 insert -> 0.02
1 insert -> 0.03

1000 inserts -> .71s
1000 inserts -> .71s
1000 inserts -> .73s

10000 inserts -> 6.75s
10000 inserts -> 6.57s
100000 inserts -> 6.57s
100000 inserts -> 7.35s

100000 inserts -> 1 minute 2 seconds


```

```
Time for querying a collection of 10,000 records
Query : (name="abe")
10,000 docs -> read 1139 docs
10,000 docs -> duration : 72.89515ms
10,000 docs -> read 1139 docs
10,000 docs -> duration : 60.839699ms
100,000 -> read 11214 docs
100,000 -> duration : 932.862374ms
100,000 -> read 11214 docs
100,000 -> duration : 1.580013132s

Query : (name="hawking" or age=77)
10,000 docs -> read 4686 docs
10,000 docs -> duration : 208.764212ms
10,000 docs -> read 4686 docs
10,000 docs -> duration : 221.073691ms
100,000 -> read 46669 docs
100,000 -> duration : 4.463323699s
100,000 -> read 46669 docs
100,000 -> duration : 6.726733346s

Query : (name="hawking" or age=77) and email="abe@g.com")
10,000 docs -> read 936 docs
10,000 docs -> duration : 42.242082ms
10,000 docs -> read 936 docs
10,000 docs -> duration : 42.355359ms
100,000 -> read 9337 docs
100,000 -> duration : 379.915291ms
100,000 -> read 9337 docs
100,000 -> duration : 393.000881ms

```

**Note** the test used for performance numbers here : re-mem/performance_test.go 

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
````



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
