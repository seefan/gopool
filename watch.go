package gopool

import (
	"time"
)

//watch spare element
func (p *Pool) watch() {
	ticker := time.Tick(time.Second * time.Duration(p.HealthSecond))
	for range ticker {
		p.avgCurrent += p.current
		p.avgCurrent /= 2
		if p.waitCount == 0 {
			p.check()
		}
	}
}

//check and close spare element
func (p *Pool) check() {
	p.lock.Lock()
	p.lock.Unlock()
	println(p.avgCurrent, p.length)
	if p.length > p.MinPoolSize {
		if p.avgCurrent+p.AcquireIncrement < p.length && !p.pooled[p.length-1].isUsed {
			p.pooled[p.length-1].Client.Close()
			p.length -= 1
		}
	}
}
