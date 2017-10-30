package gopool

var (
	checkIndex = 0
)
//watch spare element
func (p *Pool) watch() {
	timeOut := int64(p.HealthSecond)
	for t := range p.watcher.C {
		p.now = t.Unix()
		if p.waitCount == 0 && p.Status == PoolStart {
			if p.length <= p.MinPoolSize {
				p.checkMinPoolClient(timeOut)
				if p.length == 0 {
					p.Status = PoolReStart
					if err := p.init(); err == nil {
						p.Status = PoolStart
					}
				}
			} else {
				p.checkPool(timeOut)
			}
		}
	}
}

//检查最小连接池以外的连接，current以外的连接如果不用就回收
func (p *Pool) checkPool(hs int64) {
	p.lock.Lock()
	defer p.lock.Unlock()
	checkIt := p.length - 1
	if p.length > p.MinPoolSize && checkIt > p.current && p.pooled[checkIt] != nil && !p.pooled[checkIt].isUsed && p.pooled[checkIt].lastTime+hs < p.now {
		p.pooled[checkIt].Client.Close()
		p.length -= 1
	}
}

//检查最小连接池以内的连接，如果不用就ping下，以保持连接一直有数据，如果ping不能，就重启下。重启不成功就关掉。
func (p *Pool) checkMinPoolClient(hs int64) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if checkIndex < 0 || checkIndex <= p.current {
		checkIndex = p.length - 1
	}
	//同一个连接检查要间隔HealthSecond秒
	if p.pooled[checkIndex] != nil && !p.pooled[checkIndex].isUsed && p.pooled[checkIndex].lastTime+hs < p.now {
		p.pooled[checkIndex].lastTime = p.now
		if !p.pooled[checkIndex].Client.Ping() {
			p.pooled[checkIndex].Client.Close()
			if err := p.pooled[checkIndex].Client.Start(); err != nil {
				p.length -= 1
			}
		}
	}
}
