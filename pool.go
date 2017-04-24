package gopool

import (
	"github.com/seefan/goerr"
	"sync"
	"time"
)

//element interface
type Closed interface {
	//close me
	Close() error
}

// pooled element
type Element struct {
	//pos
	index int
	// The pool to which this element belongs.
	pool *Pool
	//value
	Value Closed
	//used status
	isUsed bool
}

// pool
type Pool struct {
	//pool element count
	Length int
	//used index
	Current int
	//element list
	elements []*Element
	//lock
	lock sync.Mutex
	//increment count
	AcquireIncrement int
	//time now
	now time.Time
	//get count
	avgCurrent int
	//watch delay seconds
	WatchTime int
}

func (p *Pool) init() *Pool {
	p.elements = []*Element{}
	p.now = time.Now()
	if p.WatchTime < 1 {
		p.WatchTime = 1
	}
	go p.watch()
	return p
}

//new  pool and init
func NewPool() *Pool {
	return new(Pool).init()
}

//append new element
func (p *Pool) Append(c Closed) {
	p.lock.Lock()
	defer p.lock.Unlock()
	var element *Element
	if len(p.elements) > p.Length {
		element = p.elements[p.Length]
	} else {
		element = new(Element)
		element.pool = p
		p.elements = append(p.elements, element)
	}
	element.index = p.Length
	element.Value = c
	element.isUsed = true
	p.Length += 1
}

//watch spare element
func (p *Pool) watch() {
	ticker := time.Tick(time.Second * time.Duration(p.WatchTime))
	for t := range ticker {
		p.now = t
		p.check()
	}
}

//check and close spare element
func (p *Pool) check() {
	p.lock.Lock()
	p.lock.Unlock()
	if p.avgCurrent+p.AcquireIncrement < p.Length && !p.elements[p.Length-1].isUsed {
		p.elements[p.Length-1].Value.Close()
		p.Length -= 1
	}
}

//element back,current -1
func (p *Pool) Set(element *Element) {
	p.lock.Lock()
	defer p.lock.Unlock()
	element.isUsed = false
	if element.index != p.Current-1 {
		p.elements[p.Current-1], p.elements[element.index] = p.elements[element.index], p.elements[p.Current-1]
		p.elements[p.Current-1].index, element.index = element.index, p.elements[p.Current-1].index
	}
	p.Current -= 1
}

//element closed,move to the end , length -1
func (p *Pool) CloseClient(element *Element) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if element.index != p.Length-1 {
		p.elements[p.Length-1], p.elements[element.index] = p.elements[element.index], p.elements[p.Length-1]
		p.elements[p.Length-1].index, element.index = element.index, p.elements[p.Length-1].index
	}
	p.Current -= 1
	p.Length -= 1
}

//get a element,current +1
func (p *Pool) Get() (*Element, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.Current < p.Length {
		element := p.elements[p.Current]
		element.isUsed = true
		p.avgCurrent += p.Current
		p.avgCurrent /= 2
		p.Current += 1
		return element, nil
	}
	return nil, goerr.New("pool is empty")
}

//close all
func (p *Pool) Close() {
	for _, e := range p.elements {
		e.Value.Close()
	}
	p.Current, p.Length = 0, 0
}
