package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"patterns/adapter"
	"patterns/cache"
	"patterns/fastsearcher"
	"patterns/ratecounter"
	"strconv"
	"sync"
	"time"
)

type Test struct {
	data []byte
}

func (t Test) SetData(data []byte) {
	t.data = data
}

const workersCnt = 30

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
	fmt.Println()

	fmt.Println("Adapter start")
	ctx, cancel = context.WithTimeout(context.Background(), 200*time.Second)
	adapterCh := make(chan struct{})
	adapterOut := make(chan int)
	go func() {
		for i := 0; i < 300000; i++ {
			adapterCh <- struct{}{}
		}
		close(adapterCh)
	}()
	var wg sync.WaitGroup
	wg.Add(workersCnt)
	for i := 0; i < workersCnt; i++ {
		go func() {
			defer wg.Done()
			for range adapterCh {
				v, err := adapter.Adapter(ctx)
				if err != nil {
					fmt.Printf("Adapter error %v", err)
					fmt.Println()
					continue
				}
				adapterOut <- *v
			}
		}()
	}

	go func() {
		wg.Wait()
		close(adapterOut)
	}()

	for v := range adapterOut {
		fmt.Println("Adapter value " + strconv.Itoa(v))
	}

	cancel()

	nextTime := time.Now().Add(2 * 60 * time.Second)
	ticker := time.NewTicker(3 * time.Second)
	var t Test
	for time.Now().Before(nextTime) {
		select {
		case <-ticker.C:
			t.SetData(make[])
		}
	}
}
