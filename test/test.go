package main

import "fmt"

type Test struct {
	data []int
}

func (t Test) SetData(data []int) *Test {
	t.data = data
	return &t
}

func (t Test) Data() []int {
	return t.data
}

func (t Test) SetFirst(v int) {
	t.data[0] = 123
}

func (t *Test) ExtendData() {
	fmt.Println(cap(t.data))
	t.data = append(t.data, 123)
	fmt.Println(cap(t.data))
}

func main() {
	var t Test
	b := make([]int, 1024)
	t1 := t.SetData(b)
	t2 := t.SetData(b)
	t1.SetFirst(123)
	t1.ExtendData()
	fmt.Printf("%p", t2.Data())
	fmt.Println()
	t1Data := t1.Data()
	fmt.Println(t1Data[1024])
	fmt.Printf("%p", t1.Data())
}
