package main

import (
	"fmt"
	"github.com/seefan/gopool"
	"github.com/ssdb/gossdb/ssdb"
	"log"
	"sync"
	"time"
)

type SSDBClient struct {
	isOpen   bool
	password string
	host     string
	port     int
	conn     *ssdb.Client
}

//打开连接
func (s *SSDBClient) Start() error {
	conn, err := ssdb.Connect(s.host, s.port)
	if err != nil {
		return err
	}
	s.isOpen = true
	s.conn = conn
	return nil
}
func (s *SSDBClient) Close() error {
	s.isOpen = false
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}
func (s *SSDBClient) IsOpen() bool {
	return s.isOpen
}
func (s *SSDBClient) Ping() bool {
	s.conn.Close()
	return s.Start() == nil
}

type Success struct {
	count   int
	success int
	fail    int
	lock    sync.Mutex
	wait    sync.WaitGroup
}

func (s *Success) Add() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.wait.Add(1)
	s.count += 1
}
func (s *Success) Ok() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.wait.Done()
	s.success += 1
}
func (s *Success) Fail() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.wait.Done()
	s.fail += 1
}

func (s *Success) Show() string {
	s.wait.Wait()
	return fmt.Sprintf("count:%d,success:%d,fail %d", s.count, s.success, s.fail)
}
func main() {

	p := gopool.NewPool()
	p.NewClient = func() gopool.IClient {
		return &SSDBClient{
			host: "192.168.56.101",
			port: 8888,
		}
	}
	p.MinPoolSize = 5
	p.MaxPoolSize = 100
	p.MaxWaitSize = 1000
	p.GetClientTimeout = 5
	p.HealthSecond = 10
	err := p.Start()
	if err != nil {
		log.Println(err)
		return
	}

	now := time.Now()
	wait := new(Success)
	for i := 0; i < 1; i++ {
		run(p, wait)
	}
	time.Sleep(time.Millisecond * 10)
	println(wait.Show())
	println(time.Since(now).String())
	time.Sleep(time.Minute * 5)
}
func run(p *gopool.Pool, wait *Success) {
	for i := 0; i < 100; i++ {
		wait.Add()
		go func(index int) {
			c, e := p.Get()
			if e != nil {
				log.Println(e.Error())
				wait.Fail()
				return
			}
			//time.Sleep(time.Millisecond * 15)
			p.Set(c)
			wait.Ok()
		}(i)
		//time.Sleep(time.Millisecond )
	}
}
