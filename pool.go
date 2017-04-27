package gopool

import (
	"github.com/seefan/goerr"
	"sync"
)

const (
	//连接池状态：创建
	PoolInit = 0
	//连接池状态：运行
	PoolStart = 1
	//连接池状态：关闭
	PoolStop = -1
)

// poolWait
type Pool struct {
	//poolWait element count
	length int
	//used index
	current   int
	waitCount int
	//element list
	pooled []*PooledClient
	//等待池
	poolWait chan *PooledClient //连接池
	//lock
	lock sync.Mutex
	//get count
	avgCurrent int
	//create new Closed
	NewClient func() IClient
	//状态
	Status int
	//config
	//获取连接超时时间，单位为秒。默认值: 5
	GetClientTimeout int
	//最大连接池个数。默认值: 20
	MaxPoolSize int
	//最小连接池数。默认值: 5
	MinPoolSize int
	//当连接池中的连接耗尽的时候一次同时获取的连接数。默认值: 5
	AcquireIncrement int
	//最大等待数目，当连接池满后，新建连接将等待池中连接释放后才可以继续，本值限制最大等待的数量，超过本值后将抛出异常。默认值: 1000
	MaxWaitSize int
	//连接池内缓存的连接状态检查时间隔，单位为秒。默认值: 5
	HealthSecond int
}

func (p *Pool) init() *Pool {
	p.defaultConfig()
	p.pooled = []*PooledClient{}
	p.poolWait = make(chan *PooledClient, p.MaxWaitSize)

	return p
}
func (p *Pool) defaultConfig() {
	//默认值处理
	if p.MaxPoolSize < 1 {
		p.MaxPoolSize = 20
	}
	if p.MinPoolSize < 1 {
		p.MinPoolSize = 5
	}
	if p.GetClientTimeout < 1 {
		p.GetClientTimeout = 5
	}
	if p.AcquireIncrement < 1 {
		p.AcquireIncrement = 5
	}
	if p.MaxWaitSize < 1 {
		p.MaxWaitSize = 1000
	}
	if p.HealthSecond < 1 {
		p.HealthSecond = 1
	}

	if p.MinPoolSize > p.MaxPoolSize {
		p.MinPoolSize = p.MaxPoolSize
	}
}

//启动连接池
//
//  返回 err，可能的错误，操作成功返回 nil
func (p *Pool) Start() error {
	p.waitCount = 0
	p.current = 0
	p.length = 0
	for i := 0; i < p.MinPoolSize; i++ {
		client := p.newPooledClient()
		client.pool = p
		client.index = p.length
		if err := client.Client.Start(); err != nil {
			return goerr.NewError(err, "start poolWait error")
		}
		p.pooled = append(p.pooled, client)
		p.length += 1
	}
	p.Status = PoolStart
	go p.watch()
	return nil
}

//new  poolWait and init
func NewPool() *Pool {
	return new(Pool).init()
}

//close all
func (p *Pool) Close() {
	p.Status = PoolStop
	for _, e := range p.pooled {
		e.Client.Close()
	}
	p.current, p.length = 0, 0
}

