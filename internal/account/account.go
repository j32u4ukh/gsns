package account

import (
	"fmt"
	"internal/define"
	"time"

	"github.com/j32u4ukh/glog"
	"github.com/j32u4ukh/gos"
	"github.com/j32u4ukh/gos/ans"
	"github.com/j32u4ukh/gos/ask"
	"github.com/j32u4ukh/gos/base"
	gosDefine "github.com/j32u4ukh/gos/define"
	"github.com/pkg/errors"
)

var as *AccountServer
var dbaAsker *ask.Tcp0Asker
var logger *glog.Logger

func Init(lg *glog.Logger) error {
	logger = lg

	// ==================================================
	// 與 Dba Server 建立 TCP 連線，將數據依序寫入緩存
	// ==================================================
	var port int32 = 1021
	anser, err := gos.Listen(gosDefine.Tcp0, port)
	fmt.Printf("AccountServer | Listen to port %d\n", port)

	if err != nil {
		return errors.Wrapf(err, "Failed to listen port %d", port)
	}

	as = &AccountServer{}
	as.Tcp = anser.(*ans.Tcp0Anser)
	as.Tcp.SetWorkHandler(as.Handler)
	logger.Info("完成與 Dba Server 建立 TCP 連線")

	// ==================================================
	// 與 Dba Server 建立 TCP 連線，將數據依序寫入緩存
	// ==================================================
	var address string = "127.0.0.1"
	port = 1022
	asker, err := gos.Bind(define.DbaServer, address, 1022, gosDefine.Tcp0)

	if err != nil {
		return errors.Wrapf(err, "Failed to bind address %s:%d", address, port)
	}

	dbaAsker = asker.(*ask.Tcp0Asker)
	// dbaAsker.SetWorkHandler(as.DbaHandler)
	logger.Info("DbaServer Asker 伺服器初始化完成")
	logger.Info("伺服器初始化完成")

	fmt.Printf("AccountServer | 伺服器初始化完成\n")

	// =============================================
	// 開始所有已註冊的監聽
	// =============================================
	gos.StartListen()
	fmt.Printf("AccountServer | 開始監聽\n")
	return nil
}

func Run() {
	var start time.Time
	var during, frameTime time.Duration = 0, 20 * time.Millisecond

	for {
		start = time.Now()

		gos.RunAns()
		as.Run()

		during = time.Since(start)
		if during < frameTime {
			time.Sleep(frameTime - during)
		}
	}
}

type AccountServer struct {
	Tcp *ans.Tcp0Anser
}

func (s *AccountServer) Handler(work *base.Work) {
	cmd := work.Body.PopByte()

	switch cmd {
	case 0:
		s.handleSystemCommand(work)
	case define.CommissionCommand:
		s.handleCommission(work)
	default:
		fmt.Printf("Unsupport command: %d\n", cmd)
		work.Finish()
	}
}

func (rrs *AccountServer) Run() {

}

func (s *AccountServer) handleSystemCommand(work *base.Work) {
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

func (s *AccountServer) handleCommission(work *base.Work) {
	commission := work.Body.PopUInt16()

	switch commission {
	case 1023:
		cid := work.Body.PopInt32()
		work.Body.Clear()

		work.Body.AddByte(1)
		work.Body.AddUInt16(1023)
		work.Body.AddInt32(cid)
		work.Body.AddString("Commission completed.")
		work.SendTransData()

	default:
		fmt.Printf("Unsupport commission: %d\n", commission)
		work.Finish()
	}
}
