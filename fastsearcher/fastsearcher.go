package fastsearcher

import (
	"net/http"
	"sync"
	"time"
)

const workersNum = 6

type Resp struct {
	url string
	d   time.Duration
}

func (r *Resp) GetUrl() string {
	if r == nil {
		return ""
	}
	return r.url
}

func (r *Resp) GetDuration() time.Duration {
	if r == nil {
		var v time.Duration
		return v
	}
	return r.d
}

type FastSearcher struct {
	urls     []string
	urlChan  chan string
	respChan chan Resp
}

func NewFastSearcher(urls []string) *FastSearcher {
	return &FastSearcher{
		urls:     urls,
		urlChan:  make(chan string),
		respChan: make(chan Resp),
	}
}

func (fs *FastSearcher) Search() *Resp {
	go func() {
		for _, url := range fs.urls {
			fs.urlChan <- url
		}
		close(fs.urlChan)
	}()

	var wg sync.WaitGroup
	wg.Add(workersNum)
	for range workersNum {
		go func() {
			defer wg.Done()
			for v := range fs.urlChan {
				t := time.Now()
				resp, err := http.Get(v)
				if err != nil || resp.StatusCode != 200 {
					continue
				}
				fs.respChan <- Resp{
					url: v,
					d:   time.Since(t),
				}
			}

		}()
	}

	go func() {
		wg.Wait()
		close(fs.respChan)
	}()

	var minResp *Resp
	for v := range fs.respChan {
		if minResp == nil {
			minResp = &v
			continue
		}

		if v.d < minResp.d {
			minResp = &v
		}
	}

	return minResp
}
