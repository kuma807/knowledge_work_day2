package displayGoroutine

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
