package displayGoroutine

import (
	"context"
	"fmt"
	"regexp"
	"runtime"
	"strings"
  "path/filepath"
  "os"
)

type goroutineData struct {
	GoroutineId       string
	parentGoroutineId string
	GoroutineName     string
	createdLine       string
	currentLine       string
}

func (g *goroutineData) toString() string {
	return g.GoroutineId + " " + g.parentGoroutineId + " " + g.GoroutineName + "\n"
	// return g.GoroutineId + " " + g.parentGoroutineId + " " + g.GoroutineName + " " + g.createdLine + " " + g.currentLine + "/n"
}

//大文字でpub
func Watch(ctx context.Context, goroutineName string) {
	// 監視したいゴルーチンの処理
  // runtime.GOMAXPROCS(1)
	beforeStackTrace := ""
  fileName := getFileName()
  folderName := fileName + "_" + goroutineName
  os.RemoveAll(folderName)
  creatFolder(folderName)
  f, _ := os.Create(folderName + "/tree_data.txt")
  defer f.Close()
	for {
		select {
		case <-ctx.Done():
      // show(folderName + "/tree_data.txt")
			return
		default:
      runtime.LockOSThread()
			stackTraceByte := make([]byte, 1092)
			length := runtime.Stack(stackTraceByte, true)
			stackTrace := string(stackTraceByte[:length])
      // fmt.Printf(stackTrace)
			// parentGoroutine := extractParentGoroutine(string(stackTrace[:length]))
			if beforeStackTrace != stackTrace {
				beforeStackTrace = stackTrace
				f.WriteString("start\n")
				gs := extractGoroutineData(stackTrace)
        for _, g := range gs {
          f.WriteString(g.toString())
        }
				f.WriteString("end\n")
			}
      runtime.UnlockOSThread()
		}
	}
}

func getFileName() string {
  _, filePath, _, _ := runtime.Caller(0)
  fileName := filepath.Base(filePath)
	fileNameWithoutExt := fileName[:len(fileName)-len(filepath.Ext(fileName))]
  return fileNameWithoutExt
}

func creatFolder(folderName string) {
  if _, err := os.Stat(folderName); os.IsNotExist(err) {
		// フォルダが存在しない場合にのみ作成
		err := os.Mkdir(folderName, 0755)
		if err != nil {
			fmt.Println("フォルダの作成に失敗しました:", err)
			return
		}

		fmt.Println("フォルダが作成されました")
	}
}

func extractGoroutineData(stackStr string) (gs []goroutineData) {
	lines := strings.Split(stackStr, "\n")
	childReg := regexp.MustCompile(`goroutine\s(\d+)\s\[`)
	parentReg := regexp.MustCompile(`goroutine\s(\d+)`)
	var g goroutineData
	for i, line := range lines {
		matchChild := childReg.FindStringSubmatch(line)
		if len(matchChild) > 1 {
			number := matchChild[1]
			g.GoroutineId = number
      if number == "1" {
        g.GoroutineName = "main"
        g.parentGoroutineId = "-1"
        gs = append(gs, g)
      }
		} else if strings.HasPrefix(line, "created by") {
			matchParent := parentReg.FindStringSubmatch(line)
			if len(matchParent) > 1 {
				number := matchParent[1]
				g.parentGoroutineId = number
        g.GoroutineName = lines[i - 2]
				gs = append(gs, g)
				// fmt.Printf(g.toString())
			}
		}
	}
	return gs
}
