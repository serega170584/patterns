package ratecounter

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type RateCounter struct {
	intervalCnt    uint64
	periodDuration time.Duration
	threshold      uint64
	mu             sync.RWMutex
	currPeriod     uint64
	cnt            []uint64
	alertCh        chan uint64
}

func NewRateCounter(
	intervalCnt uint64,
	periodDuration time.Duration,
	threshold uint64,
) *RateCounter {
	return &RateCounter{
		intervalCnt:    intervalCnt,
		periodDuration: periodDuration,
		threshold:      threshold,
		cnt:            make([]uint64, 1, intervalCnt),
		alertCh:        make(chan uint64),
	}
}

func (rc *RateCounter) ExecCalc(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(rc.periodDuration)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				func() {
					fmt.Println("tick")
					rc.mu.Lock()
					defer rc.mu.Unlock()
					sum := rc.calcSum()

					if sum > rc.threshold {
						fmt.Println("Alerts write start")
						rc.alertCh <- sum
						fmt.Println("Alerts write end")
					}

					if uint64(len(rc.cnt)) < rc.intervalCnt {
						rc.cnt = append(rc.cnt, 0)
						return
					}

					for i := 1; i < len(rc.cnt); i++ {
						rc.cnt[i-1] = rc.cnt[i]
					}
					rc.cnt[len(rc.cnt)-1] = 0
					fmt.Println("Tick end")
				}()
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		wg.Wait()
		close(rc.alertCh)
	}()
}

func (rc *RateCounter) calcSum() uint64 {
	var s uint64
	for _, v := range rc.cnt {
		s += v
	}
	return s
}

func (rc *RateCounter) Add() {
	fmt.Println("Add")
	rc.mu.Lock()
	atomic.AddUint64(&rc.cnt[len(rc.cnt)-1], 1)
	rc.mu.Unlock()
}

func (rc *RateCounter) GetAlertsChannel() chan uint64 {
	return rc.alertCh
}
