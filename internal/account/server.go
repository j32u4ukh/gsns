package account

import (
	"fmt"
	"internal/define"
	"time"

	"github.com/j32u4ukh/gos"
	"github.com/j32u4ukh/gos/ans"
	"github.com/j32u4ukh/gos/base"
)

type AccountServer struct {
	Tcp          *ans.Tcp0Anser
	MainServerId int32
}

func (s *AccountServer) Handler(work *base.Work) {
	cmd := work.Body.PopByte()
	logger.Info("cmd: %d", cmd)

	switch cmd {
	case define.SystemCommand:
		s.handleSystemCommand(work)
	case define.CommissionCommand:
		s.handleCommission(work)
	default:
		logger.Warn("Unsupport command: %d\n", cmd)
		work.Finish()
	}
}

func (rrs *AccountServer) Run() {

}

func (s *AccountServer) handleSystemCommand(work *base.Work) {
	service := work.Body.PopUInt16()

	switch service {
	// 回應心跳包
	case define.Heartbeat:
		logger.Debug("Heart beat! Now: %+v\n", time.Now())
		work.Body.Clear()
		work.Body.AddByte(0)
		work.Body.AddUInt16(0)
		work.Body.AddString("OK")
		work.SendTransData()
	default:
		logger.Warn("Unsupport service: %d\n", service)
		work.Finish()
	}
}

func (s *AccountServer) handleCommission(work *base.Work) {
	commission := work.Body.PopUInt16()
	logger.Info("commission: %d", commission)

	switch commission {
	case 1023:
		cid := work.Body.PopInt32()
		work.Body.Clear()

		work.Body.AddByte(1)
		work.Body.AddUInt16(1023)
		work.Body.AddInt32(cid)
		work.Body.AddString("Commission completed.")
		work.SendTransData()

	case define.Register:
		s.MainServerId = work.Index
		cid := work.Body.PopInt32()
		bs := work.Body.PopByteArray()
		logger.Info("MainServerId: %d, cid: %d, bs: %+v", s.MainServerId, cid, bs)
		work.Finish()

		td := base.NewTransData()
		td.AddByte(define.CommissionCommand)
		td.AddUInt16(define.Register)
		td.AddInt32(cid)

		// Account data for register
		td.AddByteArray(bs)

		data := td.FormData()
		logger.Info("data: %+v", data)

		// 將註冊數據傳到 Dba 伺服器
		err := gos.SendToServer(define.DbaServer, &data, td.GetLength())

		if err != nil {
			logger.Error("Failed to send to server %d: %v\nError: %+v", define.DbaServer, data, err)
			return
		}

	case define.Login:
		s.MainServerId = work.Index
		cid := work.Body.PopInt32()
		bs := work.Body.PopByteArray()
		logger.Info("MainServerId: %d, cid: %d, bs: %+v", s.MainServerId, cid, bs)
		work.Finish()

		td := base.NewTransData()
		td.AddByte(define.CommissionCommand)
		td.AddUInt16(define.Login)
		td.AddInt32(cid)

		// Account data for register
		td.AddByteArray(bs)

		data := td.FormData()
		logger.Info("data: %+v", data)

		// 將註冊數據傳到 Dba 伺服器
		err := gos.SendToServer(define.DbaServer, &data, td.GetLength())

		if err != nil {
			logger.Error("Failed to send to server %d: %v\nError: %+v", define.DbaServer, data, err)
			return
		}

	default:
		fmt.Printf("Unsupport commission: %d\n", commission)
		work.Finish()
	}
}
