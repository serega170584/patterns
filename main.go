package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"patterns/ratecounter"
	"time"
)

func main() {
	go func() {
		fmt.Println("Pprof server started on :6060")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatalf("Pprof server failed: %v", err)
		}
	}()

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
			for range rand.Intn(200) {
				rc.Add()
			}
			time.Sleep(time.Second)
		}
	}()

	for v := range rc.GetAlertsChannel() {
		fmt.Printf("Alert: %d errors", v)
		fmt.Println()
	}
}
