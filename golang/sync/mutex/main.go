package main

import "sync"

func main() {
	var mu sync.Mutex
	mu.Lock()
	// sync的mutex是不可重入锁
	// mu.Lock()
	// mu.Unlock()
	mu.Unlock()
}
