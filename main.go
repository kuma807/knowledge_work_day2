package main

import (
	"fmt"
	"runtime"
	"strings"
	"time"
	"context"
	"sync"
	"regexp"
)

type goroutineData struct {
	GoroutineId string
	parentGoroutineId string
	GoroutineName string
	createdLine string
	currentLine string
}

func (g *goroutineData) toString() string {
	return g.GoroutineId + " " + g.parentGoroutineId + "\n"
	// return g.GoroutineId + " " + g.parentGoroutineId + " " + g.GoroutineName + " " + g.createdLine + " " + g.currentLine + "/n"
}

func watchGoroutine(ctx context.Context) {
	// 監視したいゴルーチンの処理
	beforeStackTrace := ""
	for {
		select {
		case <-ctx.Done():
			fmt.Println("監視用ゴルーチンが終了します。")
			return
		default:
			stackTraceByte := make([]byte, 8192)
			length := runtime.Stack(stackTraceByte, true)
			stackTrace := string(stackTraceByte[:length])
			// parentGoroutine := extractParentGoroutine(string(stackTrace[:length]))
			if beforeStackTrace != stackTrace {
				beforeStackTrace = stackTrace
				fmt.Printf("start\n")
				// fmt.Printf(stackTrace)
				extractGoroutineData(stackTrace)
				fmt.Printf("end\n")
			}
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go watchGoroutine(ctx)

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

func extractGoroutineData(stackStr string) (gs []goroutineData) {
	lines := strings.Split(stackStr, "\n")
	childReg := regexp.MustCompile(`goroutine\s(\d+)\s\[`)
	parentReg := regexp.MustCompile(`goroutine\s(\d+)`)
	var g goroutineData
	for _, line := range lines {
		matchChild := childReg.FindStringSubmatch(line)
		if len(matchChild) > 1 {
			number := matchChild[1]
			g.GoroutineId = number
		} else if strings.HasPrefix(line, "created by") {
			matchParent := parentReg.FindStringSubmatch(line)
			if len(matchParent) > 1 {
				number := matchParent[1]
				g.parentGoroutineId = number
				gs = append(gs, g)
				fmt.Printf(g.toString())
			}
		}
	}
	return gs
}
