package main

import (
	"fmt"
	pm "internal/post_message"

	"github.com/j32u4ukh/glog/v2"
	"github.com/j32u4ukh/gos"
	gosDefine "github.com/j32u4ukh/gos/define"
	gosUtils "github.com/j32u4ukh/gos/utils"
)

func main() {
	fmt.Println("Hello, post message!")
	option := glog.DefaultOption(false, false, 8, "../log")

	gosUtils.GosConfig.AnswerConnectNumbers[gosDefine.Tcp0] = 10000
	gosUtils.GosConfig.AnswerWorkNumbers[gosDefine.Tcp0] = 10000
	gosUtils.GosConfig.AskerWorkNumbers[gosDefine.Tcp0] = 10000
	gosLogger := glog.SetLogger(0, "gos", glog.DebugLevel)
	gosLogger.SetSkip(3)
	gosLogger.SetOptions(option)
	gos.SetLogger(gosLogger)
	logger := glog.SetLogger(1, "post-message-server", glog.DebugLevel)
	logger.SetOptions(option)
	logger = glog.SetLogger(2, "post-message-client", glog.DebugLevel)
	logger.SetOptions(option)

	err := pm.Init()
	if err != nil {
		logger.Error(fmt.Sprintf("Account server 初始化失敗, err: %+v", err))
		return
	}
	pm.Run()
}
