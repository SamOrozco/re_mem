package main

import (
	"os"
	"re-mem/re_mem"
)

type User struct {
	Name string
	Age  int
}

func main() {
	store := re_mem.NewLocalStorage("/Users/samorozco/testdb")
	col, err := store.GetCollection("users")
	if err != nil {
		println("failed to get collection")
		os.Exit(1)
	}
	key, err := col.Create(&User{
		Name: "sam",
		Age:  100,
	})

	if err != nil {
		panic(err)
	}
	print(key)
}
