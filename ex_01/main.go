package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func iotaTest() {
	type CarType int
	const (
		Sedan CarType = iota + 1
		SUV
		Crossover
	)
	var t CarType
	t = SUV
	fmt.Println("Car type: ", t)

	type CarOption uint64
	const (
		GPS CarOption = 1 << iota
		AWD
		SunRoof
		HeatedSeat
	)
	var o CarOption
	o = SunRoof | HeatedSeat
	if o&SunRoof != 0 {
		fmt.Println("with SunRoof")
	}
}

func New(text string) error {
	return &errorString{text}
}

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

func errorTest() {
	var EOF = errors.New(("EOF"))
	fmt.Println(EOF)
}

func errorHandlingTest() {
	f, err := os.Open(("man.go"))
	if err != nil {
		fmt.Println("Open error")
	}

	r := bufio.NewReader((f))
	l, err := r.ReadString(('\n'))
	if err != nil {
		fmt.Println("Read error")
	}
	fmt.Println("Read: ", l)
}

type Portion int

const (
	Regular Portion = iota
	Small
	Large
)

type Udon struct {
	men Portion
	aburaage bool
	ebiten uint
}

func NewUdon(p Portion, aburaage bool, ebiten uint) *Udon {
	return &Udon{
		men: p,
		aburaage: aburaage,
		ebiten: ebiten,
	}
}

func optionTest() {
	var temuraUdon = NewUdon(Large, false, 2)
	fmt.Println("temuraUdon aburaage: ", temuraUdon.aburaage)
}

func main() {
	iotaTest()
	errorTest()
	errorHandlingTest()
	optionTest()
}