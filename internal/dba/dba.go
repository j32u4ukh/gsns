package dba

import (
	"fmt"
	"internal/define"
	"time"

	"github.com/j32u4ukh/gos"
	"github.com/j32u4ukh/gos/ans"
	"github.com/j32u4ukh/gos/base"
	gosDefine "github.com/j32u4ukh/gos/define"
)

var s *DbaServer

func Init() {
	var port int = 1022
	anser, err := gos.Listen(gosDefine.Tcp0, int32(port))
	fmt.Printf("DbaServer | Listen to port %d\n", port)

	if err != nil {
		fmt.Printf("DbaServer | Error: %+v\n", err)
		return
	}

	s = &DbaServer{}
	tcpAnser := anser.(*ans.Tcp0Anser)
	tcpAnser.SetWorkHandler(s.Handler)
	fmt.Printf("DbaServer | 伺服器初始化完成\n")

	// =============================================
	// 開始所有已註冊的監聽
	// =============================================
	gos.StartListen()
	fmt.Printf("DbaServer | 開始監聽\n")
}

func Run() {
	var start time.Time
	var during, frameTime time.Duration = 0, 20 * time.Millisecond

	for {
		start = time.Now()

		gos.RunAns()
		s.Run()

		during = time.Since(start)
		if during < frameTime {
			time.Sleep(frameTime - during)
		}
	}
}

type DbaServer struct {
}

func (s *DbaServer) Handler(work *base.Work) {
	cmd := work.Body.PopByte()

	switch cmd {
	case 0:
		s.handleSystemCommand(work)
	case 1:
	case define.CommissionCommand:
		s.handleCommission(work)

	default:
		fmt.Printf("Unsupport command: %d\n", cmd)
		work.Finish()
	}
}

func (rrs *DbaServer) Run() {

}

func (s *DbaServer) handleSystemCommand(work *base.Work) {
	service := work.Body.PopUInt16()

	switch service {
	// 回應心跳包
	case 0:
		fmt.Printf("Heart beat! Now: %+v\n", time.Now())
		work.Body.Clear()
		work.Body.AddByte(0)
		work.Body.AddUInt16(0)
		work.Body.AddString("OK")
		work.SendTransData()
	default:
		fmt.Printf("Unsupport service: %d\n", service)
		work.Finish()
	}
}

func (s *DbaServer) handleCommission(work *base.Work) {
	commission := work.Body.PopUInt16()
	var cid int32 = work.Body.PopInt32()

	switch commission {
	case 1023:
		work.Body.Clear()
		work.Body.AddByte(1)
		work.Body.AddUInt16(1023)
		work.Body.AddInt32(cid)
		work.Body.AddString("Commission completed.")
		work.SendTransData()

	case define.Register:
		// TODO: 建立使用者資料

	default:
		fmt.Printf("Unsupport commission: %d\n", commission)
		work.Finish()
	}
}
