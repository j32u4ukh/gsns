package account

import (
	"fmt"
	"internal/define"
	"internal/pbgo"
	"time"

	"github.com/j32u4ukh/gos"
	"github.com/j32u4ukh/gos/ans"
	"github.com/j32u4ukh/gos/base"
	"google.golang.org/protobuf/proto"
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
	case define.Heartbeat:
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

func (s *AccountServer) DbaHandler(work *base.Work) {
	cmd := work.Body.PopByte()
	logger.Info("cmd: %d", cmd)

	switch cmd {
	case define.SystemCommand:
		s.handleDbaSystemCommand(work)
	case define.CommissionCommand:
		s.handleDbaCommission(work)
	default:
		fmt.Printf("Unsupport command: %d\n", cmd)
		work.Finish()
	}
}

func (s *AccountServer) handleDbaSystemCommand(work *base.Work) {
	service := work.Body.PopUInt16()

	switch service {
	// 回應心跳包
	case define.Heartbeat:
		fmt.Printf("Heart response Now: %+v\n", time.Now())
		work.Finish()
	default:
		fmt.Printf("Unsupport service: %d\n", service)
		work.Finish()
	}
}

func (s *AccountServer) handleDbaCommission(work *base.Work) {
	commission := work.Body.PopUInt16()

	switch commission {
	case define.Register:
		cid := work.Body.PopInt32()
		returnCode := work.Body.PopUInt16()
		bs := work.Body.PopByteArray()
		work.Finish()

		account := &pbgo.Account{}
		err := proto.Unmarshal(bs, account)

		if err != nil {
			return
		}

		logger.Info("New account created : %+v", account)

		td := base.NewTransData()
		td.AddByte(define.CommissionCommand)
		td.AddUInt16(define.Register)
		td.AddInt32(cid)
		td.AddUInt16(returnCode)

		// Account data for register
		td.AddByteArray(bs)

		data := td.FormData()

		// 將註冊結果回傳主伺服器
		err = gos.SendToClient(define.AccountPort, s.MainServerId, &data, td.GetLength())

		if err != nil {
			logger.Error("Failed to send to client %d: %v\nError: %+v", s.MainServerId, data, err)
			return
		}
	case define.Login:
		cid := work.Body.PopInt32()
		returnCode := work.Body.PopUInt16()
		token := work.Body.PopUInt64()
		work.Finish()

		td := base.NewTransData()
		td.AddByte(define.CommissionCommand)
		td.AddUInt16(define.Login)
		td.AddInt32(cid)
		td.AddUInt16(returnCode)
		td.AddUInt64(token)
		data := td.FormData()

		// 將註冊結果回傳主伺服器
		err := gos.SendToClient(define.AccountPort, s.MainServerId, &data, td.GetLength())

		if err != nil {
			logger.Error("Failed to send to client %d: %v\nError: %+v", s.MainServerId, data, err)
			return
		}

	default:
		fmt.Printf("Unsupport commission: %d\n", commission)
		work.Finish()
	}
}
