package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
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

type Option struct {
	men Portion
	aburaage bool
	ebiten uint
}

func NewUdon2(opt Option) *Udon {
	if opt.ebiten == 0 && time.Now().Hour() < 10 {
		opt.ebiten = 1
	}
	return &Udon{
		men: opt.men,
		aburaage: opt.aburaage,
		ebiten:  opt.ebiten,
	}
}

func noOptionTest() {
	var opt Option
	var udon = NewUdon2(opt)
	fmt.Println("ebiten: ", udon.ebiten)
}

type fluentOpt struct {
	men Portion
	aburaage bool
	ebiten uint
}

func NewUdon3(p Portion) *fluentOpt {
	return &fluentOpt{
		men: p,
		aburaage: false,
		ebiten: 1,
	}
}

func (o *fluentOpt) Aburaage() *fluentOpt {
	o.aburaage = true
	return o
}

func (o *fluentOpt) Ebiten(n uint) *fluentOpt {
	o.ebiten = n
	return o
}

func (o *fluentOpt) Order() *Udon {
	return &Udon{
		men: o.men,
		aburaage: o.aburaage,
		ebiten: o.ebiten,
	}
}

func optionBuilderTest() {
	var udon = NewUdon3(Large).Aburaage().Order()
	fmt.Println("Aburaage: ", udon.aburaage)
}

type OptFunc func (r *Udon)

func NewUdon4(opts ...OptFunc) *Udon {
	r := &Udon{}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func OptMen(p Portion) OptFunc {
	return func(r *Udon) { r.men = p }
}

func OptAburaage() OptFunc {
	return func(r *Udon) { r.aburaage = true }
}

func OptEbiten(n uint) OptFunc {
	return func(r *Udon) { r.ebiten = n }
}

func functionalOptionTest() {
	var tokuseiUdon = NewUdon4(OptAburaage(), OptEbiten(3))
	fmt.Println("tokuseiUdon ebiten: ", tokuseiUdon.ebiten)
}

func commandlineTest() {
	var (
		FlagStr = flag.String("string", "default", "文字列")
		FlagInt = flag.Int("int", 1, "数値")
	)

	flag.Parse()
	log.Println(*FlagStr)
	log.Println(*FlagInt)
	log.Println(flag.Args())
}

type Config struct {
	Port uint16 `envconfig:"PORT" default:"3000"`
	Host string `envconfig:"HOST" required:"true"`
	AdminPort uint16 `envconfig:"ADMIN_PORT" default:"3001"`
}

func envTest() {
	var c Config
	err := envconfig.Process("", &c)
	if err != nil {
		fmt.Println(err)
	}
	log.Println(c)
}

func memoryTest() {
	s1 := make([]int, 1000)
	fmt.Println(len(s1))
	fmt.Println(cap(s1))

	s2 := make([]int, 0, 1000)
	fmt.Println(len(s2))
	fmt.Println(cap(s2))

	m := make(map[string]string, 1000)
	fmt.Println(len(m))
}

func stringConnectionTest() {
	src := []string{"Back", "To", "The", "Future", "Part", "III"}
	var builder strings.Builder
	builder.Grow(100)
	for i, word := range src {
		if i != 0 {
			builder.WriteByte(' ')
		}
		builder.WriteString(word)
	}
	log.Println(builder.String())
}

func timeTest() {
	now := time.Now()
	tz, _ := time.LoadLocation("America/Los_Angeles")
	future := time.Date(2015, time.October, 21, 7, 28, 0, 0, tz)
	fmt.Println(now.String())
	fmt.Println(future.Format(time.RFC3339Nano))
}

func timeDurationTest() {
	fiveMinute := 5 * time.Minute
	fmt.Println("fiveMinute: ", fiveMinute)
	var seconds int = 10
	tenSeconds := time.Duration(seconds) * time.Second
	fmt.Println("tenSeconds: ", tenSeconds)

	past := time.Date(1955, time.November, 12, 6, 38, 0, 0, time.UTC)
	dur := time.Now().Sub(past)
	fmt.Println("dur: ", dur)

	fiveMinuteAfter := time.Now().Add(fiveMinute)
	fiveMinuteBefore := time.Now().Add(-fiveMinute)
	fmt.Println("fiveMinuteAfter: ", fiveMinuteAfter)
	fmt.Println("fiveMinuteBefore: ", fiveMinuteBefore)

	// fmt.Println("3 seconds start")
	// time.Sleep(3 * time.Second)
	// fmt.Println("3 seconds end")

	jst, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Date(2021, 6, 8, 20, 56, 00, 000, jst)
	nextMonth := now.AddDate(0, 1, 0)
	fmt.Println(nextMonth)
}

type Book struct {
	Title string
	Author string
	Publisher string
	ReleasedAt time.Time
	ISBN string
}

type Person struct {
	FirstName string
	LastName string
}

func NewPerson(first, last string) *Person {
	return &Person{
		FirstName: first,
		LastName: last,
	}
}

type Parent struct{}

func (p Parent) m1() {
	p.m2()
}

func (p Parent) m2() {
	fmt.Println("Parent")
}

type Child struct {
	Parent
}

func (c Child) m2() {
	fmt.Println("Child")
}

func structTest() {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	book := Book{
		Title: "Real world",
		Author: "Shibukawa",
		Publisher: "olily",
		ISBN: "48",
		ReleasedAt: time.Date(2017, time.June, 14, 0, 0, 0, 0, jst),
	}
	fmt.Println(book.Title)

	person := NewPerson("aaa", "bbb")
	fmt.Println("First: ", person.FirstName)
	fmt.Println("Last: ", person.LastName)

	c := Child{}
	c.m1()
	c.m2()
}

func main() {
	iotaTest()
	errorTest()
	errorHandlingTest()
	optionTest()
	noOptionTest()
	optionBuilderTest()
	functionalOptionTest()
	commandlineTest()
	envTest()
	memoryTest()
	stringConnectionTest()
	timeTest()
	timeDurationTest()
	structTest()
}