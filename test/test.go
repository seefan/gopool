package main

import (
	"github.com/seefan/gopool"
	"sync"
	"time"
)

type TestClient struct {
	isOpen bool
}

//打开连接
func (s *TestClient) Start() error {
	s.isOpen = true
	return nil
}
func (s *TestClient) Close() error {
	s.isOpen = false
	return nil
}
func (s *TestClient) IsOpen() bool {
	return s.isOpen
}
func (s *TestClient) Ping() bool {
	return s.Start() == nil
}
func main() {
	pool := gopool.NewPool()
	pool.NewClient = func() gopool.IClient {
		return &TestClient{}
	}
	pool.MinPoolSize = 50
	pool.MaxPoolSize = 1000
	pool.MaxWaitSize = 10000
	pool.GetClientTimeout = 10
	pool.HealthSecond = 10
	if err := pool.Start(); err == nil {
		test(pool, 10, 10)
		test(pool, 50, 10)
		test(pool, 100, 10)
		test(pool, 200, 10)
		test(pool, 500, 10)
		test(pool, 800, 10)
		test(pool, 1000, 10)
		test(pool, 3000, 10)
		test(pool, 5000, 10)
	}
}

func test(pool *gopool.Pool, threadCount, callCount int) {
	now := time.Now()
	wait := new(sync.WaitGroup)
	failed := 0
	for i := 0; i < threadCount; i++ {
		wait.Add(1)
		go func(p *gopool.Pool, w *sync.WaitGroup, idx int) {
			for j := 0; j < callCount; j++ {
				if c, e := p.Get(); e != nil {
					failed += 1
				} else {
					p.Set(c)
				}
			}
			w.Done()
		}(pool, wait, i)

	}
	wait.Wait()
	println("thread=", threadCount, "call=", callCount, "failed=", failed, "time=", time.Since(now).String())
}
