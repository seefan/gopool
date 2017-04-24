package gopool

import (
	"strconv"
	"testing"
)

type TestValue struct {
	Name string
}

func (t *TestValue) Close() error {
	return nil
}

func TestPool_Append(t *testing.T) {
	p := NewPool()
	e := new(TestValue)
	p.Append(e)
}

func TestPool_Get(t *testing.T) {
	p := NewPool()
	for i := 0; i < 10; i++ {
		e := &TestValue{"name" + strconv.Itoa(i)}
		p.Append(e)
	}

	c, e := p.Get()
	if e != nil {
		t.Errorf("error is ", e)
	}
	t.Log(c.index)
}
func TestPool_Set(t *testing.T) {
	p := NewPool()
	for i := 0; i < 10; i++ {
		e := &TestValue{"name" + strconv.Itoa(i)}
		p.Append(e)
	}

	c, e := p.Get()
	if e != nil {
		t.Errorf("error is ", e)
	}
	t.Log(c.index)
	t.Log("Current is ", p.Current)
	c, e = p.Get()
	if e != nil {
		t.Errorf("error is ", e)
	}
	t.Log(c.index)
	t.Log("Current is ", p.Current)
	p.Set(c)
	t.Log("Current is ", p.Current)
}
func TestNewPool(t *testing.T) {
	NewPool()
}
