package main

import (
	"fmt"
	"internal/gsns"

	"github.com/j32u4ukh/glog/v2"
	"github.com/j32u4ukh/gos"
	gosDefine "github.com/j32u4ukh/gos/define"
	gosUtils "github.com/j32u4ukh/gos/utils"
)

func main() {
	fmt.Println("Hello, gsns!")
	option := glog.DefaultOption(false, false, 8, "../log")

	gosLogger := glog.SetLogger(0, "gos", glog.DebugLevel)
	gosLogger.SetSkip(3)
	gosLogger.SetOptions(option)
	gosUtils.GosConfig.AnswerConnectNumbers[gosDefine.Tcp0] = 10000
	gosUtils.GosConfig.AnswerWorkNumbers[gosDefine.Tcp0] = 10000
	gosUtils.GosConfig.AskerWorkNumbers[gosDefine.Tcp0] = 10000
	gosUtils.GosConfig.AnswerConnectNumbers[gosDefine.Http] = 10000
	gosUtils.GosConfig.AnswerWorkNumbers[gosDefine.Http] = 10000
	gos.SetLogger(gosLogger)
	logger := glog.SetLogger(1, "gsns-server", glog.DebugLevel)
	logger.SetOptions(option)
	logger = glog.SetLogger(2, "gsns-client", glog.DebugLevel)
	logger.SetOptions(option)

	err := gsns.Init()
	if err != nil {
		logger.Error(fmt.Sprintf("gsns 初始化失敗, err: %+v", err))
		return
	}
	gsns.Run()
}
