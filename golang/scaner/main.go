package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func task() {
	fmt.Println("task start...")
	start := time.Now().UnixNano()
	var i int
	for {
		i++
		if i >= 1000000000 {
			break
		}
	}
	end := time.Now().UnixNano()
	fmt.Println("task end, coast =", end-start)
}

func main() {
	scaner := bufio.NewScanner(os.Stdin)
	for scaner.Scan() {
		go task()
	}
}
