package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"patterns/cache"
	"patterns/fastsearcher"
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

	fmt.Println("Rate counter run")

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

	fmt.Println("Cache run")

	c := cache.NewCache()
	c.Set("num1", 12, 3*time.Second)
	fmt.Println("cache set check")
	fmt.Println(*c.Get("num1") == 12)

	time.Sleep(4 * time.Second)
	fmt.Println("expired cache value check")
	fmt.Println(c.Get("num1") == nil)

	c.Set("num1", 12, 3*time.Second)
	time.Sleep(2 * time.Second)

	c.Set("num1", 34, 4*time.Second)
	time.Sleep(2 * time.Second)
	fmt.Println(*c.Get("num1") == 34)

	fastSearcher := fastsearcher.NewFastSearcher([]string{
		"https://www.google.com",
		"https://ya.ru",
		"https://www.pik.ru",
		"https://www.avito.ru",
		"https://www.ozon.ru",
		"https://vk.com",
		"https://booking.com",
		"https://ctc.ru",
		"https://www.1tv.ru",
		"https://www.okko.tv",
	})

	fmt.Println("Fastest site")
	s := fastSearcher.Search()
	fmt.Printf("Site %s", s.GetUrl())
	fmt.Printf("Duration %s", s.GetDuration().String())
}
