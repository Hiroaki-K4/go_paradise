package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
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
}