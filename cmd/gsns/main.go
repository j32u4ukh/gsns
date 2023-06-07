package main

import (
	"fmt"
	"internal/gsns"

	"github.com/j32u4ukh/glog/v2"
)

func main() {
	fmt.Println("Hello, gsns!")
	logger := glog.SetLogger(0, "gsns", glog.DebugLevel)
	logger.SetFolder("../log")
	logger.SetOptions(glog.DefaultOption(true, true), glog.UtcOption(8))
	err := gsns.Init()
	if err != nil {
		logger.Error(fmt.Sprintf("gsns 初始化失敗, err: %+v", err))
		return
	}
	gsns.Run()
}
