package gopool

//element interface
type IClient interface {
	Start() error
	//close me
	Close() error
	IsOpen() bool
}

// pooled element
type PooledClient struct {
	//pos
	index int
	// The poolWait to which this element belongs.
	pool *Pool
	//value
	Client IClient
	//used status
	isUsed bool
}
