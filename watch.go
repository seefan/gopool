package gopool

import (
	"time"
)

var (
	lastTime time.Time
)

//watch spare element
func (p *Pool) watch() {
	for p.Status > 0 {
		p.avgCurrent += p.current
		p.avgCurrent /= 2
		if p.waitCount == 0 {
			p.check()
			if time.Since(lastTime).Seconds() > float64(p.HealthSecond) {
				p.checkClient()
				if p.length == 0 {
					p.Status = PoolReStart
					if err := p.init(); err == nil {
						p.Status = PoolStart
					}
				}
			}
		}
		time.Sleep(time.Second)
	}
}

//check and close spare element
func (p *Pool) check() {
	p.lock.Lock()
	p.lock.Unlock()
	if p.length > p.MinPoolSize {
		if p.length-1 > p.current && p.avgCurrent+p.AcquireIncrement < p.length && !p.pooled[p.length-1].isUsed {
			p.pooled[p.length-1].Client.Close()
			p.length -= 1
		}
	}
}
func (p *Pool) checkClient() bool {
	if c, err := p.getPoolClient(); err == nil {
		if c.Client.Ping() == false {
			p.closeClient(c)
			return p.checkClient()
		} else {
			p.setPoolClient(c)
		}
	}
	return true
}
