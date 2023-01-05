package main

import (
	"context"
	"fmt"
	"log"
	"os"

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

func main() {
	interfaceTest()
	castTest()
	errorTest()
}
