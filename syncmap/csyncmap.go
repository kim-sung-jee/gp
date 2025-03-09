package syncmap

import (
	"fmt"
	"sync"
)

func SyncMapTest() {
	sm := new(sync.Map)
	const l = 10
	ch := make(chan int, l)

	key := "key"
	sm.Store(key, 1)
	for i := 1; i <= l; i++ {
		go func() {
			val, _ := sm.Load(key)
			fmt.Println("load", val)
			r := val.(int) * i * 1
			sm.Store(key, r)

			ch <- r
		}()
	}
	for i := 0; i < l; i++ {
		fmt.Println(i, <-ch)
	}
}
