package main

import (
	"github.com/seefan/gopool"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type TestValue struct {
	Name string
}

func (t *TestValue) Close() error {
	return nil
}
func main() {
	p := gopool.NewPool()
	p.AcquireIncrement = 3
	for i := 0; i < 20; i++ {
		e := &TestValue{"name" + strconv.Itoa(i)}
		p.Append(e)
	}

	go run(p, 0)
	time.Sleep(time.Microsecond * 13)
	go run(p, 10)
	time.Sleep(time.Minute)
}
func run(p *gopool.Pool, start int) {
	for i := start; i < start+10; i++ {
		go func(index int) {
			c, e := p.Get()
			if e != nil {
				log.Print(e)
			}
			mss := rand.Intn(2000)
			time.Sleep(time.Microsecond * time.Duration(mss))
			p.Set(c)
		}(i)
		ms := rand.Intn(500)
		time.Sleep(time.Microsecond * time.Duration(ms))
	}
}
