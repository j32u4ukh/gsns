package main

import (
	"fmt"
	"internal/gsns"

	"github.com/j32u4ukh/glog"
)

func main() {
	fmt.Println("Hello, gsns!")
	logger := glog.GetLogger("../log", "gsns", glog.DebugLevel, false)
	logger.SetOptions(glog.DefaultOption(true, true), glog.UtcOption(8))
	err := gsns.Init(logger)
	if err != nil {
		logger.Error(fmt.Sprintf("gsns 初始化失敗, err: %+v", err))
		return
	}
	gsns.Run()
}
