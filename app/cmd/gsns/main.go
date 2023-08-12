package main

import (
	"fmt"
	"internal/gsns"

	"github.com/j32u4ukh/glog/v2"
	"github.com/j32u4ukh/gos"
)

func main() {
	fmt.Println("Hello, gsns!")
	logger := glog.SetLogger(0, "gsns", glog.DebugLevel)
	logger.SetFolder("../log")
	logger.SetOptions(glog.DefaultOption(true, true), glog.UtcOption(8))

	gosLogger := glog.SetLogger(1, "gos", glog.DebugLevel)
	gosLogger.SetFolder("../log")
	gosLogger.SetOptions(glog.DefaultOption(true, true), glog.UtcOption(8))
	gosLogger.SetSkip(3)
	gos.SetLogger(gosLogger)

	err := gsns.Init()
	if err != nil {
		logger.Error(fmt.Sprintf("gsns 初始化失敗, err: %+v", err))
		return
	}
	gsns.Run()
}
