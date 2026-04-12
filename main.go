package main

import (
	"context"
	"fmt"
	"math/rand"
	"patterns/ratecounter"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rc := ratecounter.NewRateCounter(
		8,
		time.Second,
		1000,
	)
	rc.ExecCalc(ctx)

	go func() {
		for range 29 {
			for range rand.Intn(400) {
				rc.Add()
			}
			time.Sleep(time.Second)
		}
	}()

	for v := range rc.GetAlertsChannel() {
		fmt.Printf("Alert: %d errors", v)
	}
}
