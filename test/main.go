package main

import (
	"github.com/seefan/gopool"
	"log"
//	"math/rand"
	"strconv"
	"sync"
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
	for i := 0; i < 100; i++ {
		e := &TestValue{"name" + strconv.Itoa(i)}
		p.Append(e)
	}
	now := time.Now()
	wait := new(sync.WaitGroup)
	for i := 0; i < 83; i++ {
		go run(p, wait)
	}
	wait.Wait()
	println(time.Since(now).String())
}
func run(p *gopool.Pool, wait *sync.WaitGroup) {
	for i := 0; i < 1000; i++ {
		wait.Add(1)
		go func(index int) {
			defer wait.Done()
			c, e := p.Get()
			if e != nil {
				log.Print(e)
				return
			}
			time.Sleep(time.Millisecond)
			p.Set(c)
		}(i)
		time.Sleep(time.Millisecond)
	}
}
