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
	len := 10
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

	fmt.Println("result", arr)
}

func someFunction(name string) int {
	time.Sleep(time.Second * 3)
	//fmt.Println("called processingRequest function, name is ", name)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(100)
}
func distinctElements(arr []int) []int {
	elements := make(map[int]bool)
	result := make([]int, len(arr))
	for i := 0; i < len(arr); i++ {
		if _, value := elements[arr[i]]; !value {
			elements[arr[i]] = true
			result = append(result, arr[i])
		}
	}
	return result
}
