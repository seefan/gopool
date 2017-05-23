package main

import (
	"fmt"
	"github.com/seefan/gopool"
	"github.com/ssdb/gossdb/ssdb"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
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
	f, _ := os.Create("profile_file")
	pprof.StartCPUProfile(f)     // 开始cpu profile，结果写到文件f中
	defer pprof.StopCPUProfile() // 结束profile
	p := gopool.NewPool()
	p.NewClient = func() gopool.IClient {
		return &SSDBClient{
			host: "192.168.56.101",
			port: 8888,
		}
	}
	p.MinPoolSize = 10
	p.MaxPoolSize = 10
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
	for i := 0; i < 500; i++ {
		go run(p, wait)
	}
	time.Sleep(time.Second * 5)
	println(wait.Show())
	println(time.Since(now).String())
}
func run(p *gopool.Pool, wait *Success) {
	for i := 0; i < 1000; i++ {
		wait.Add()
		go func(index int) {
			c, e := p.Get()
			if e != nil {
				log.Println(e.Error())
				wait.Fail()
				return
			}
			rnd := rand.Float32()
			time.Sleep(time.Duration(rnd) * time.Second)
			p.Set(c)
			wait.Ok()
		}(i)
	}
}
