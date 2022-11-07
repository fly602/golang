package main

import "fmt"

const (
	STACK_SIZE = (8 * 1024 * 1024)
)

func main() {
	var a [2 * STACK_SIZE]int
	a[0] = 1
	a[2*STACK_SIZE-1] = 2
	fmt.Println(a[2*STACK_SIZE-1])
}
