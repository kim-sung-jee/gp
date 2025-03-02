package main

import (
	"fmt"
	"math/rand"
	singleflight "p/signleflight-v2"
	"sync"
	"time"
)

func main() {

	wait := new(sync.WaitGroup)

	g := singleflight.Group{}
	key := "key"
	len := 1000
	var arr []int
	results := make(chan int, len)
	for i := 0; i < len; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			result, _, _ := g.Do(key, func() (interface{}, error) {
				return someFunction(key), nil
			})
			results <- result.(int)
		}()
	}
	wait.Wait()
	close(results)
	for result := range results {
		arr = append(arr, result)
	}

	fmt.Println("result", allElementsSame(arr))
}

func someFunction(name string) int {
	time.Sleep(time.Second)
	//fmt.Println("called processingRequest function, name is ", name)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(100)
}

func allElementsSame(arr []int) bool {
	if len(arr) == 0 {
		return true
	}

	first := arr[0]
	for _, num := range arr {
		if num != first {
			return false
		}
	}
	return true
}
