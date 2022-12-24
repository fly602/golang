package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type atomic_lock struct {
	data int32
	old  int32
	new  int32
}

func (al *atomic_lock) Lock() {
	for {
		if atomic.CompareAndSwapInt32(&al.data, al.old, al.new) {
			return
		}
	}
}

func (al *atomic_lock) Unlock() {
	al.data = al.old
}

func (al *atomic_lock) init() {
	al.old = rand.New(rand.NewSource(time.Now().Unix())).Int31()
	al.data = al.old
	for {
		al.new = rand.New(rand.NewSource(time.Now().Unix())).Int31()
		if al.old != al.new {
			return
		}
	}
}

var g_val int64

func add() {
	defer wg.Done()
	for {
		al.Lock()
		if g_val >= MAX_VAL {
			al.Unlock()
			return
		}
		g_val++
		al.Unlock()
	}

}

func add2() {
	defer wg.Done()
	for {
		mu.Lock()
		if g_val >= MAX_VAL {
			mu.Unlock()
			return
		}
		g_val++
		mu.Unlock()
	}

}

var al atomic_lock
var wg sync.WaitGroup
var mu sync.Mutex

const (
	GO_NUM  = 10
	MAX_VAL = GO_NUM * 10000000
)

func main() {
	al.init()
	wg.Add(GO_NUM)
	fmt.Println("atomic_lock=", al)
	for i := 0; i < GO_NUM; i++ {
		go add()
	}
	wg.Wait()
	fmt.Println("g_val=", g_val)
}
