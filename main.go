package main

import (
	"context"
	"sync"
	"time"

	"github.com/kuma807/knowledge_work_day2/displayGoroutine"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go displayGoroutine.Watch(ctx, "testGoroutine")
	time.Sleep(time.Second * 1)
	numParentGoroutines := 1
	var wg sync.WaitGroup

	// ゴルーチンを発生させる
	for i := 0; i < numParentGoroutines; i++ {
		wg.Add(1)
		go parent(&wg)
	}
	wg.Wait()
	time.Sleep(time.Second * 2)
	cancel()
	displayGoroutine.Show("testGoroutine")
	time.Sleep(time.Second * 5)
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
	time.Sleep(5 * time.Second)
}

func child2(wg *sync.WaitGroup) {
	defer (*wg).Done()
	time.Sleep(6 * time.Second)
}
