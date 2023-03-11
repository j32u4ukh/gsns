package main

import (
	"fmt"
	"internal/account"

	"github.com/j32u4ukh/glog"
)

func main() {
	fmt.Println("Hello, dba!")
	logger := glog.GetLogger("../log", "gsns", glog.DebugLevel, false)
	logger.SetOptions(glog.DefaultOption(true, true), glog.UtcOption(8))
	err := account.Init(logger)
	if err != nil {
		logger.Error(fmt.Sprintf("Dba server 初始化失敗, err: %+v", err))
		return
	}
	account.Run()
}
