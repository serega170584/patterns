package adapter

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func Adapter(ctx context.Context) (*int, error) {
	ch := make(chan int)

	go func() {
		ch <- longProcess()
		close(ch)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case v := <-ch:
		return &v, nil
	case <-time.After(6 * time.Second):
		return nil, fmt.Errorf("3 seconds is out")
	}
}

func longProcess() int {
	time.Sleep(time.Duration(rand.Intn(40)) * time.Second)

	return rand.Intn(10000)
}
