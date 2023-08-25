package main

import (
	"fmt"
	"internal/account"

	"github.com/j32u4ukh/glog/v2"
	"github.com/j32u4ukh/gos"
	gosDefine "github.com/j32u4ukh/gos/define"
	gosUtils "github.com/j32u4ukh/gos/utils"
)

func main() {
	fmt.Println("Hello, account!")
	gosUtils.GosConfig.AnswerConnectNumbers[gosDefine.Tcp0] = 10000
	gosUtils.GosConfig.AnswerWorkNumbers[gosDefine.Tcp0] = 10000
	gosUtils.GosConfig.AskerWorkNumbers[gosDefine.Tcp0] = 10000
	option := glog.DefaultOption(false, false, 8, "../log")
	gosLogger := glog.SetLogger(0, "gos", glog.DebugLevel)
	gosLogger.SetSkip(3)
	gosLogger.SetOptions(option)
	gos.SetLogger(gosLogger)
	logger := glog.SetLogger(1, "account-server", glog.DebugLevel)
	logger.SetOptions(option)
	logger = glog.SetLogger(2, "account-client", glog.DebugLevel)
	logger.SetOptions(option)

	err := account.Init()
	if err != nil {
		logger.Error(fmt.Sprintf("Account server 初始化失敗, err: %+v", err))
		return
	}
	account.Run()
}
