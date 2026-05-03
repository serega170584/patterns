package main

import (
	"fmt"
	"sync"
)

const producersCnt = 3

func main() {
	list := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	in := make(chan int)

	go func() {
		for _, v := range list {
			in <- v
		}
		close(in)
	}()

	var wg sync.WaitGroup

	producers := make([]chan int, producersCnt)
	for i := 0; i < producersCnt; i++ {
		producers[i] = makeProducer(in, &wg)
	}

	consumer := makeConsumer(producers, &wg)

	for v := range consumer {
		fmt.Println(v)
	}

}

func makeProducer(in chan int, wg *sync.WaitGroup) chan int {
	wg.Add(1)
	res := make(chan int)
	go func() {
		defer wg.Done()
		for v := range in {
			res <- v
		}
		close(res)
	}()

	return res
}

func makeConsumer(producers []chan int, wg *sync.WaitGroup) chan int {
	wg.Add(len(producers))
	consumer := make(chan int)
	for _, producer := range producers {
		go func(producer chan int) {
			defer wg.Done()
			for v := range producer {
				consumer <- v
			}
		}(producer)
	}

	go func() {
		wg.Wait()
		close(consumer)
	}()

	return consumer
}
