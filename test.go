package main

import (
  "runtime"
  "path/filepath"
)

func main() {
    _, file, _, _ := runtime.Caller(0)
    fileName := filepath.Base(file)

  	// 拡張子を除いた部分だけを取り出す
  	fileNameWithoutExt := fileName[:len(fileName)-len(filepath.Ext(fileName))]
    println(fileNameWithoutExt)
}
