package gopool

import (
	"github.com/seefan/goerr"
	"time"
)

//element back,Current -1
func (p *Pool) setPoolClient(element *PooledClient) {
	element.isUsed = false
	if element.index != p.current-1 {
		p.pooled[p.current-1], p.pooled[element.index] = p.pooled[element.index], p.pooled[p.current-1]
		p.pooled[p.current-1].index, element.index = element.index, p.pooled[p.current-1].index
	}
	p.current -= 1
}

//归还连接到连接池
//
//  cc 连接
func (p *Pool) Set(element *PooledClient) {
	lastTime = time.Now()
	if element == nil {
		return
	}
	if p.Status == PoolStart {
		if element.Client.IsOpen() {
			p.lock.Lock()
			defer p.lock.Unlock()
			if p.waitCount > 0 { //有等待的连接
				p.poolWait <- element
				p.waitCount -= 1
			} else {
				p.setPoolClient(element)
			}
		} else {
			p.closeClient(element)
		}
	} else {
		if element.Client.IsOpen() {
			element.Client.Close()
		}
	}
}

//element closed,move to the end , length -1
func (p *Pool) closeClient(element *PooledClient) {
	if element.index != p.length-1 {
		p.pooled[p.length-1], p.pooled[element.index] = p.pooled[element.index], p.pooled[p.length-1]
		p.pooled[p.length-1].index, element.index = element.index, p.pooled[p.length-1].index
	}
	p.current -= 1
	p.length -= 1
}

//get a element,Current +1
func (p *Pool) getPoolClient() (*PooledClient, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.current < p.length {
		element := p.pooled[p.current]
		element.isUsed = true
		p.current += 1
		return element, nil
	}
	return nil, goerr.New("poolWait is empty")
}

//在连接池取一个新连接
//
//  返回 client，一个新的连接
//  返回 err，可能的错误，操作成功返回 nil
func (p *Pool) Get() (client *PooledClient, err error) {
	switch p.Status {
	case PoolStop:
		return nil, goerr.New("the Connectors is Closed, can not get new client.")
	case PoolInit:
		return nil, goerr.New("the Connectors is not inited, can not get new client.")
	}
	//检查是否有缓存的连接
	client, err = p.getPoolClient()
	if err == nil {
		return
	}
	//检查是否可以扩展
	if err = p.poolAppend(); err == nil {
		client, err = p.getPoolClient()
		if err == nil {
			return
		}
	}

	p.lock.Lock()
	if p.waitCount >= p.MaxWaitSize {
		p.lock.Unlock()
		return nil, goerr.New("poolWait is busy,Wait for connection creation has reached %d", p.waitCount)
	}
	p.waitCount += 1
	p.lock.Unlock()
	//enter slow poolWait
	timeout := time.After(time.Duration(p.GetClientTimeout) * time.Second)
	select {
	case <-timeout:
		p.lock.Lock()
		p.waitCount -= 1
		p.lock.Unlock()
		return nil, goerr.New("ssdb poolWait is busy,can not get new client in %d seconds", p.GetClientTimeout, p.current, p.length, p.waitCount)
	case cc := <-p.poolWait:
		if cc == nil {
			return nil, goerr.New("the Connectors is Closed, can not get new client.")
		}
		return cc, nil
	}

}

func (p *Pool) newPooledClient() *PooledClient {
	return &PooledClient{
		Client: p.NewClient(),
		index:  len(p.pooled),
		pool:   p,
	}
}
