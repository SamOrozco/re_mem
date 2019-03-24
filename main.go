package main

import (
	"os"
	"re-mem/re_mem"
)

type User struct {
	Name    string
	Age     int
	Company string
}

func main() {
	store := re_mem.NewLocalStorage("/Users/samorozco/madison")
	col, err := store.GetCollection("users")
	if err != nil {
		println("failed to get collection")
		os.Exit(1)
	}

	doc, err := col.Get("46d4f45c-2109-45f5-96d9-e13e2de819c0")
	println(doc.String())

	//docs, err := col.Query("company", "amazon")
	//if err != nil {
	//	panic(err)
	//}
	//for _, e := range docs {
	//	println(e.String())
	//}

	//names := []string{"sam", "madison", "steven", "april", "sawyer"}
	//ages := []int{50, 55, 34, 231, 121, 12, 77}
	//company := []string{"apple", "microsoft", "amazon", "dropbox", "this"}
	//
	//for i := 0; i < 1000; i++ {
	//	key, err := col.Create(&User{
	//		Name:    names[rand.Int()%len(names)],
	//		Age:     ages[rand.Int()%len(ages)],
	//		Company: company[rand.Int()%len(company)],
	//	})
	//
	//	if err != nil {
	//		panic(err)
	//	}
	//	println(key)
	//}
}
