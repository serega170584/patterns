package main

import "fmt"

const cnt = 10000
const outCnt = 5

func main() {
	in := make(chan int)
	out := make(chan int, outCnt)

	go func() {
		for v := range in {
			select {
			case out <- v:
				continue
			default:
				<-out
				out <- v
			}
		}
		close(out)
	}()

	for i := 0; i < cnt; i++ {
		in <- i
	}
	close(in)

	for v := range out {
		fmt.Println(v)
	}
}
