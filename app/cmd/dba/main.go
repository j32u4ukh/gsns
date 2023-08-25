package main

import (
	"fmt"
	"internal/dba"

	"github.com/j32u4ukh/glog/v2"
	"github.com/j32u4ukh/gos"
	gosDefine "github.com/j32u4ukh/gos/define"
	gosUtils "github.com/j32u4ukh/gos/utils"
)

func main() {
	fmt.Println("Hello, dba!")
	option := glog.DefaultOption(false, false, 8, "../log")

	gosUtils.GosConfig.AnswerConnectNumbers[gosDefine.Tcp0] = 10000
	gosUtils.GosConfig.AnswerWorkNumbers[gosDefine.Tcp0] = 10000
	gosUtils.GosConfig.AskerWorkNumbers[gosDefine.Tcp0] = 10000
	gosLogger := glog.SetLogger(0, "gos", glog.DebugLevel)
	gosLogger.SetSkip(3)
	gosLogger.SetOptions(option)
	gos.SetLogger(gosLogger)
	logger := glog.SetLogger(1, "dba-server", glog.DebugLevel)
	logger.SetOptions(option)
	logger = glog.SetLogger(2, "dba-client", glog.DebugLevel)
	logger.SetOptions(option)

	err := dba.Init()
	if err != nil {
		logger.Error(fmt.Sprintf("Dba server 初始化失敗, err: %+v", err))
		return
	}
	dba.Run()
}
