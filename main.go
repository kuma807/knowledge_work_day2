package main

import (
	"context"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go displayGoroutine.watchGoroutine(ctx)

	numParentGoroutines := 1
	var wg sync.WaitGroup

	// ゴルーチンを発生させる
	for i := 0; i < numParentGoroutines; i++ {
		wg.Add(1)
		go parent(&wg)
	}
	wg.Wait()
	cancel()
}

func parent(wg *sync.WaitGroup) {
	defer (*wg).Done()
	wg.Add(2)
	go child1(wg)
	go child2(wg)
	time.Sleep(2 * time.Second)
}

func child1(wg *sync.WaitGroup) {
	defer (*wg).Done()
	time.Sleep(1 * time.Second)
}

func child2(wg *sync.WaitGroup) {
	defer (*wg).Done()
	time.Sleep(2 * time.Second)
}
