package main

import (
	"fmt"
	"internal/account"

	"github.com/j32u4ukh/glog/v2"
	"github.com/j32u4ukh/gos"
)

func main() {
	fmt.Println("Hello, account!")
	logger := glog.SetLogger(0, "gsns", glog.DebugLevel)
	logger.SetOptions(glog.UtcOption(8))
	logger.SetOptions(glog.FolderOption("../log", glog.ShiftDayAndSize, 1, 5*glog.MB))
	logger.SetOptions(glog.BasicOption(&glog.Option{
		Level:     glog.DebugLevel,
		ToConsole: true,
		ToFile:    false,
		FileInfo:  true,
		LineInfo:  true,
	}))
	logger.SetOptions(glog.BasicOption(&glog.Option{
		Level:     glog.InfoLevel,
		ToConsole: true,
		ToFile:    false,
		FileInfo:  true,
		LineInfo:  true,
	}))
	logger.SetOptions(glog.BasicOption(&glog.Option{
		Level:     glog.WarnLevel,
		ToConsole: true,
		ToFile:    true,
		FileInfo:  true,
		LineInfo:  true,
	}))
	logger.SetOptions(glog.BasicOption(&glog.Option{
		Level:     glog.ErrorLevel,
		ToConsole: true,
		ToFile:    true,
		FileInfo:  true,
		LineInfo:  true,
	}))

	gosLogger := glog.SetLogger(1, "gos", glog.DebugLevel)
	gosLogger.SetSkip(3)
	gosLogger.SetOptions(glog.UtcOption(8))
	gosLogger.SetOptions(glog.FolderOption("../log", glog.ShiftDayAndSize, 1, 5*glog.MB))
	gosLogger.SetOptions(glog.BasicOption(&glog.Option{
		Level:     glog.DebugLevel,
		ToConsole: true,
		ToFile:    false,
		FileInfo:  true,
		LineInfo:  true,
	}))
	gosLogger.SetOptions(glog.BasicOption(&glog.Option{
		Level:     glog.InfoLevel,
		ToConsole: true,
		ToFile:    false,
		FileInfo:  true,
		LineInfo:  true,
	}))
	gosLogger.SetOptions(glog.BasicOption(&glog.Option{
		Level:     glog.WarnLevel,
		ToConsole: true,
		ToFile:    true,
		FileInfo:  true,
		LineInfo:  true,
	}))
	gosLogger.SetOptions(glog.BasicOption(&glog.Option{
		Level:     glog.ErrorLevel,
		ToConsole: true,
		ToFile:    true,
		FileInfo:  true,
		LineInfo:  true,
	}))
	gos.SetLogger(gosLogger)

	err := account.Init()
	if err != nil {
		logger.Error(fmt.Sprintf("Account server 初始化失敗, err: %+v", err))
		return
	}
	account.Run()
}
