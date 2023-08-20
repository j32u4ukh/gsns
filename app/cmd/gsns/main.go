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
	options := []glog.IOption{
		glog.UtcOption(8),
		glog.FolderOption("../log", glog.ShiftDayAndSize, 1, 5*glog.MB),
		glog.BasicOption(&glog.Option{
			Level:     glog.DebugLevel,
			ToConsole: true,
			ToFile:    false,
			FileInfo:  true,
			LineInfo:  true,
		}),
		glog.BasicOption(&glog.Option{
			Level:     glog.InfoLevel,
			ToConsole: true,
			ToFile:    false,
			FileInfo:  true,
			LineInfo:  true,
		}),
		glog.BasicOption(&glog.Option{
			Level:     glog.WarnLevel,
			ToConsole: true,
			ToFile:    true,
			FileInfo:  true,
			LineInfo:  true,
		}),
		glog.BasicOption(&glog.Option{
			Level:     glog.ErrorLevel,
			ToConsole: true,
			ToFile:    true,
			FileInfo:  true,
			LineInfo:  true,
		}),
	}

	gosLogger := glog.SetLogger(0, "gos", glog.DebugLevel)
	gosLogger.SetSkip(3)
	gosLogger.SetOptions(options...)
	gosUtils.GosConfig.AnswerConnectNumbers[gosDefine.Tcp0] = 10000
	gosUtils.GosConfig.AnswerWorkNumbers[gosDefine.Tcp0] = 10000
	gosUtils.GosConfig.AskerWorkNumbers[gosDefine.Tcp0] = 10000

	gosUtils.GosConfig.AnswerConnectNumbers[gosDefine.Http] = 10000
	gosUtils.GosConfig.AnswerWorkNumbers[gosDefine.Http] = 10000
	gos.SetLogger(gosLogger)
	logger := glog.SetLogger(1, "gsns-server", glog.DebugLevel)
	logger.SetOptions(options...)
	// logger.SetSkip(3)
	logger = glog.SetLogger(2, "gsns-client", glog.DebugLevel)
	logger.SetOptions(options...)
	// logger.SetSkip(3)

	err := gsns.Init()
	if err != nil {
		logger.Error(fmt.Sprintf("gsns 初始化失敗, err: %+v", err))
		return
	}
	gsns.Run()
}
