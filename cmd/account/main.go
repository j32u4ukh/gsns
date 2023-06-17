package main

import (
	"fmt"
	"internal/account"

	"github.com/j32u4ukh/glog/v2"
)

func main() {
	fmt.Println("Hello, account!")
	logger := glog.SetLogger(0, "gsns", glog.DebugLevel)
	logger.SetFolder("../log")
	logger.SetOptions(glog.UtcOption(8))
	logger.SetOptions(glog.BasicOption(glog.DebugLevel, true, false, true, true))
	logger.SetOptions(glog.BasicOption(glog.InfoLevel, true, false, true, true))
	logger.SetOptions(glog.BasicOption(glog.WarnLevel, true, true, true, true))
	logger.SetOptions(glog.BasicOption(glog.ErrorLevel, true, true, true, true))
	err := account.Init()
	if err != nil {
		logger.Error(fmt.Sprintf("Account server 初始化失敗, err: %+v", err))
		return
	}
	account.Run()
}
