package gsns

import (
	"fmt"
	"internal/define"

	"github.com/j32u4ukh/gos"
	"github.com/j32u4ukh/gos/base"
	"github.com/j32u4ukh/gos/base/ghttp"
)

func (s *MainServer) register(c *ghttp.Context) {
	// 將當前 Http 的工作結束
	s.HttpAnswer.Finish(c)

	// 形成 "建立使用者" 的請求
	td := base.NewTransData()
	td.AddByte(define.CommissionCommand)
	td.AddUInt16(define.Register)
	td.AddInt32(c.GetId())
	data := td.FormData()
	err := gos.SendToServer(EDbaServer, &data, td.GetLength())

	if err != nil {
		fmt.Printf("(s *MainServer) CommissionHandler | Failed to send to server %d: %v\nError: %+v\n", EDbaServer, data, err)
		return
	}
}

func (s *MainServer) login(c *ghttp.Context) {
	c.Json(200, ghttp.H{
		"index": 2,
		"msg":   "POST | /",
	})
}

func (s *MainServer) logout(c *ghttp.Context) {
	c.Json(200, ghttp.H{
		"index": 2,
		"msg":   "POST | /",
	})
}
