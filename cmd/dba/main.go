package main

import (
	"fmt"
	"internal/dba"

	"github.com/j32u4ukh/glog/v2"
	"github.com/j32u4ukh/gos"
)

func main() {
	fmt.Println("Hello, dba!")
	logger := glog.SetLogger(0, "gsns", glog.DebugLevel)
	logger.SetFolder("../log")
	logger.SetOptions(glog.UtcOption(8))
	logger.SetOptions(glog.BasicOption(glog.DebugLevel, true, false, true, true))
	logger.SetOptions(glog.BasicOption(glog.InfoLevel, true, false, true, true))
	logger.SetOptions(glog.BasicOption(glog.WarnLevel, true, true, true, true))
	logger.SetOptions(glog.BasicOption(glog.ErrorLevel, true, true, true, true))

	gosLogger := glog.SetLogger(1, "gos", glog.DebugLevel)
	gosLogger.SetFolder("../log")
	gosLogger.SetOptions(glog.DefaultOption(true, true), glog.UtcOption(8))
	gosLogger.SetSkip(3)
	gos.SetLogger(gosLogger)

	err := dba.Init()
	if err != nil {
		logger.Error(fmt.Sprintf("Dba server 初始化失敗, err: %+v", err))
		return
	}
	dba.Run()
}
