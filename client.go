package gopool

//element interface
type IClient interface {
	//start connection
	Start() error
	//close connection
	Close() error
	//connection status
	IsOpen() bool
	//connection ping
	Ping() bool
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
