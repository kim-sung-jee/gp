package signleflight_v2

import (
	"sync"
)

type call struct {
	wg sync.WaitGroup

	val   any
	err   error
	dups  int
	chans []chan<- Result
}

type Group struct {
	m sync.Map
}

type Result struct {
	Val    any
	Err    error
	Shared bool
}

func (g *Group) Do(key string, fn func() (any, error)) (v any, err error, shared bool) {
	actual, loaded := g.m.LoadOrStore(key, new(call))
	c := actual.(*call)
	if loaded {
		c.dups++
		c.wg.Wait()

		return c.val, c.err, true
	}

	c.wg.Add(1)

	g.doCall(c, key, fn)

	return c.val, c.err, c.dups > 0
}

func (g *Group) DoChan(key string, fn func() (any, error)) <-chan Result {
	ch := make(chan Result, 1)
	actual, loaded := g.m.LoadOrStore(key, &call{chans: []chan<- Result{ch}})
	c := actual.(*call)
	if loaded {
		c.dups++
		c.chans = append(c.chans, ch)
		return ch
	}
	c.wg.Add(1)
	go g.doCall(c, key, fn)
	return ch
}

func (g *Group) doCall(c *call, key string, fn func() (any, error)) {
	//val, _ := g.m.Load(key)
	//c := val.(*call)

	c.val, c.err = fn()
	// 안됨
	c.wg.Done()
	if val, _ := g.m.Load(key); val == c {
		g.m.Delete(key)
	}

	for _, ch := range c.chans {
		ch <- Result{c.val, c.err, c.dups > 0}
	}
}

func (g *Group) ForgetUnshared(key string) bool {
	val, ok := g.m.Load(key)
	if !ok {
		return true
	}
	c := val.(*call)
	if c.dups == 0 {
		g.m.Delete(key)
		return true
	}
	return false
}
