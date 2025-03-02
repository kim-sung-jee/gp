package signleflight_v2

import (
	"fmt"
	"p/goid"
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
		if c.val == nil {
			c.val = -100
		}
		return c.val, c.err, true
	}

	c.wg.Add(1)

	g.doCall(key, fn)

	return c.val, c.err, c.dups > 0
}

func (g *Group) doCall(key string, fn func() (any, error)) {
	val, _ := g.m.LoadAndDelete(key)
	c := val.(*call)
	if f, err := goid.Goid(); err == nil {
		fmt.Println("goroutine id is ", f)
	}

	c.val, c.err = fn()

	for _, ch := range c.chans {
		ch <- Result{c.val, c.err, c.dups > 0}
	}

	c.wg.Done()
}
