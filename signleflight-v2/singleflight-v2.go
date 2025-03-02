package signleflight_v2

import (
	"fmt"
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
	fmt.Println(loaded)
	if loaded {
		//fmt.Print("key already exists, value is ", c.val, "\n")
		c.dups++
		c.wg.Wait()

		return c.val, c.err, true
	}

	c.wg.Add(1)

	g.doCall(key, fn)

	return c.val, c.err, c.dups > 0
}

func (g *Group) doCall(key string, fn func() (any, error)) {
	val, _ := g.m.Load(key)
	c := val.(*call)

	c.val, c.err = fn()
	g.m.Delete(key)
	c.wg.Done()

	for _, ch := range c.chans {
		ch <- Result{c.val, c.err, c.dups > 0}
	}
}
