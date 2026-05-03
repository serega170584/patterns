package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const workersCnt = 10000

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	in := make(chan int)
	go func(ctx context.Context) {
		for range 10000000 {
			select {
			case <-ctx.Done():
				break
			default:
				in <- rand.Intn(100000)
			}
		}
		close(in)
	}(ctx)

	out := make(chan int)
	var wg sync.WaitGroup
	wg.Add(workersCnt)
	for i := 0; i < workersCnt; i++ {
		go func() {
			defer wg.Done()
			for v := range in {
				out <- v
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	for v := range out {
		fmt.Println(v)
	}
}
