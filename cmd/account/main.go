package main

import (
	"fmt"
	"internal/account"

	"github.com/j32u4ukh/glog/v2"
)

func main() {
	fmt.Println("Hello, dba!")
	logger := glog.SetLogger(0, "gsns", glog.DebugLevel)
	logger.SetFolder("../log")
	logger.SetOptions(glog.DefaultOption(true, true), glog.UtcOption(8))
	err := account.Init()
	if err != nil {
		logger.Error(fmt.Sprintf("Dba server 初始化失敗, err: %+v", err))
		return
	}
	account.Run()
}
