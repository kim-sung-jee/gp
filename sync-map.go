package main

import "sync"

func main() {
	m := sync.Map{}

	wg := new(sync.WaitGroup)
	l := 100000
	results := make(chan int, l)
	for i := 0; i < l; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			actual, load := m.LoadOrStore(l, l)
			if !load {
				results <- actual.(int)
			}
		}()
	}
	wg.Wait()
	close(results)
	for result := range results {
		println(result)
	}
}
