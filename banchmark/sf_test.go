package benchmark

import (
	signleflight_v2 "p/signleflight-v2"
	"p/singleflight"
	"sync"
	"testing"
)

func dummyFunc() (interface{}, error) {
	return 42, nil
}

func BenchmarkSingleflight(b *testing.B) {
	var g singleflight.Group
	key := "benchmark_key"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.Do(key, dummyFunc)
	}
}

func BenchmarkSingleflightV2(b *testing.B) {
	var g signleflight_v2.Group
	key := "benchmark_key"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.Do(key, dummyFunc)
	}
}

func BenchmarkSingleflightConcurrent(b *testing.B) {
	var g singleflight.Group
	key := "benchmark_key"
	var wg sync.WaitGroup

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			g.Do(key, dummyFunc)
		}()
	}
	wg.Wait()
}

func BenchmarkSingleflightV2Concurrent(b *testing.B) {
	var g signleflight_v2.Group
	key := "benchmark_key"
	var wg sync.WaitGroup

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			g.Do(key, dummyFunc)
		}()
	}
	wg.Wait()
}
