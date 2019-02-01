package main

import (
	"github.com/seefan/gopool"
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
		c,err:=pool.Get()
		c1,err:=pool.Get()
		if err==nil{
			c.Client.Ping()
			c1.Client.Close()
			pool.Set(c1)
			pool.Set(c)
		}
	}
}


