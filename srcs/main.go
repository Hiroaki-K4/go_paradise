package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"encoding/json"
	"encoding/csv"
	"io"

	"github.com/gen2brain/beeep"
)

type Warning interface {
	Show(message string)
}

type ConsoleWarning struct{}

func (c ConsoleWarning) Show(message string) {
	fmt.Fprintf(os.Stderr, "[%s]: %s\n", os.Args[0], message)
}

type DesktopWarning struct{}

func (d DesktopWarning) Show(message string) {
	beeep.Alert(os.Args[0], message, "")
}

func interfaceTest() {
	var warn Warning

	warn = &ConsoleWarning{}
	warn.Show("Hello World to console")

	warn = &DesktopWarning{}
	warn.Show("Hello World to desktop")
}

func castTest() {
	ctx := context.WithValue(context.Background(), "favorite", "zenigata")

	if s, ok := ctx.Value("favorite").(string); ok {
		log.Printf("I love %s\n", s)
	}

	switch v := ctx.Value("favorite").(type) {
	case string:
		log.Printf("I love %s\n", v)
	case int:
		log.Printf("I love %d\n", v)
	case complex128:
		log.Printf("I love %f\n", v)
	default:
		log.Printf("I love %v\n", v)
	}
}

func validate(length int) error {
	if length <= 0 {
		return fmt.Errorf("length must be greater than 0, length = %d", length)
	}

	return nil
}

func errorTest() {
	len := -10
	error := validate(len)
	fmt.Println(error)
}

type HTTPError struct {
	StatusCode int
	URL string
}

type ip struct {
	Origin string `json:"origin"`
	URL string `json:"url"`
}

func jsonTest() {
	f, err := os.Open("ip.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var resp ip
	if err := json.NewDecoder(f).Decode(&resp); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", resp)
}

func jsonSliceTest() {
	type user struct {
		UserID string `json:"user_id"`
		UserName string `json:"user_name"`
		Languages []string `json:"languages"`
	}

	{
		u := user{
			UserID: "001",
			UserName: "gopher",
		}
		b, _ := json.Marshal(u)
		fmt.Println(string(b))
	}

	{
		u := user{
			UserID: "001",
			UserName: "gopher",
			Languages: []string{},
		}
		b, _ := json.Marshal(u)
		fmt.Println(string(b))
	}
}

func omitEmptyTest() {
	type FormInput struct {
		Name string `json:"name"`
		CompanyName string `json:"company_name,omitempty"`
	}

	in := FormInput{Name: "yamada"}

	b, _ := json.Marshal(in)
	fmt.Println(string(b))
}

func csvReaderTest () {
	f, err := os.Open("country.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(record)
	}
}

func csvWriterTest() {
	records := [][]string{
		{"Name", "year", "page"},
		{"Go lang web dev", "2016", "280"},
		{"Go lang thread", "2018", "256"},
		{"Go lang interpreter", "2018", "316"},
	}

	f, err := os.OpenFile("oreilly.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatal(err)
		}
	}

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	interfaceTest()
	castTest()
	errorTest()
	jsonTest()
	jsonSliceTest()
	omitEmptyTest()
	csvReaderTest()
	csvWriterTest()
}
