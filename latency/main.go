package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// 1 2 3 4 5 1 2 3 4 5 6 7
// 1 2 3 4 5
// 2 3 4 5   1
// 3 4 5       1
// 4 5           1
// 5               1
//
// 2 3 4 5 1
func main() {
	var maxLatency *time.Duration
	tt := time.Now()
	var cnt uint64
	var mu sync.Mutex
	defer func() {
		fmt.Println("Total time:" + time.Since(tt).String())
	}()

	ch := make(chan struct{})
	outCh := make(chan struct{})

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				mu.Lock()
				fmt.Println("Count " + strconv.Itoa(int(cnt)))
				cnt = 0
				mu.Unlock()
			case <-ch:
				outCh <- struct{}{}
				close(outCh)
				return
			}
		}
	}()

	var wg sync.WaitGroup
	for range 20 {
		wg.Add(500)
		it := time.Now()
		for i := 0; i < 500; i++ {
			go func() {
				defer wg.Done()
				t := time.Now()
				defer func() {
					if maxLatency == nil {
						v := time.Since(t)
						maxLatency = &v
						return
					}
					if time.Since(t) > *maxLatency {
						v := time.Since(t)
						maxLatency = &v
					}
				}()
				for range 1 {
					mu.Lock()
					atomic.AddUint64(&cnt, 1)
					mu.Unlock()
					_, _ = http.Get("http://localhost:8080/list/")
				}
			}()
		}
		wg.Wait()
		func() {
			fmt.Println("Iteration time: " + time.Since(it).String())
		}()
	}

	ch <- struct{}{}
	close(ch)

	<-outCh

	fmt.Println("Max goroutine latency: " + (*maxLatency).String())
}
