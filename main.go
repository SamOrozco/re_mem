package main

type User struct {
	Name  string
	Age   int
	Email string
}

func main() {
	store := NewLocalStorage("/Users/samorozco/first_db")
	objects, err := store.GetCollection("objects")
	if err != nil {
		panic(err)
	}

}
