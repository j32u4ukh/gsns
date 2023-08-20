package main

import (
	"fmt"
	"internal/dba"

	"github.com/j32u4ukh/glog/v2"
	"github.com/j32u4ukh/gos"
)

func main() {
	fmt.Println("Hello, dba!")
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
	gos.SetLogger(gosLogger)
	logger := glog.SetLogger(1, "dba-server", glog.DebugLevel)
	logger.SetOptions(options...)
	// logger.SetSkip(3)
	logger = glog.SetLogger(2, "dba-client", glog.DebugLevel)
	logger.SetOptions(options...)
	// logger.SetSkip(3)

	err := dba.Init()
	if err != nil {
		logger.Error(fmt.Sprintf("Dba server 初始化失敗, err: %+v", err))
		return
	}
	dba.Run()
}
