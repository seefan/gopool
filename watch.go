package gopool

import (
	"log"
	"time"
)

//watch spare element
func (p *Pool) watch() {
	ticker := time.Tick(time.Second * time.Duration(p.HealthSecond))
	for t := range ticker {
		log.Println("watch", t, p.current, p.avgCurrent, p.waitCount, p.length)
		p.check()
	}
}

//check and close spare element
func (p *Pool) check() {
	p.lock.Lock()
	p.lock.Unlock()
	if p.length > p.MinPoolSize {
		p.avgCurrent += p.current
		p.avgCurrent /= 2
		if p.avgCurrent+p.AcquireIncrement < p.length && !p.pooled[p.length-1].isUsed {
			p.pooled[p.length-1].Client.Close()
			p.length -= 1
		}
	}
}
