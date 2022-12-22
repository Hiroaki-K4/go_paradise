package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type Book struct {
	Title string
	Author string
	Publisher string
	ReleasedAt time.Time
	ISBN string
}

func main() {
	f, err := os.Open("book.json")
	if err != nil {
		log.Fatal("file open error: ", err)
	}
	d := json.NewDecoder((f))
	var b Book
	d.Decode(&b)
	fmt.Println(b)
}