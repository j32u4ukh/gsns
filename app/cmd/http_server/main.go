package main

import (
	"fmt"
	"time"

	"github.com/j32u4ukh/gos"
	"github.com/j32u4ukh/gos/ans"
	"github.com/j32u4ukh/gos/define"
)

func main() {
	fmt.Println("Starting http server...")
	var port int = 1023
	anser, err := gos.Listen(define.Http, int32(port))
	fmt.Printf("RunAns | Listen to port %d\n", port)

	if err != nil {
		fmt.Printf("Error: %+v\n", err)
		return
	}

	httpAnswer := anser.(*ans.HttpAnser)
	mrg := &Mgr{
		Answer: httpAnswer,
	}
	mrg.Handler(httpAnswer.Router)

	fmt.Printf("(s *Service) RunAns | 伺服器初始化完成\n")
	gos.StartListen()
	fmt.Printf("(s *Service) RunAns | 開始監聽\n")
	gos.SetFrameTime(200 * time.Millisecond)
	gos.Run(nil)
}
