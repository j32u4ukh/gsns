package account

import (
	"fmt"
	"internal/define"
	"internal/pbgo"
	"time"

	"github.com/j32u4ukh/gos"
	"github.com/j32u4ukh/gos/base"
	"google.golang.org/protobuf/proto"
)

func (s *AccountServer) DbaHandler(work *base.Work) {
	cmd := work.Body.PopByte()
	logger.Info("cmd: %d", cmd)

	switch cmd {
	case define.SystemCommand:
		s.handleDbaSystemCommand(work)
	case define.NormalCommand:
		s.handleDbaNormalCommand(work)
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

func (s *AccountServer) handleDbaNormalCommand(work *base.Work) {
	service := work.Body.PopUInt16()

	switch service {
	// 取得用戶資訊
	case define.GetUserData:
		logger.Debug("GetUserData")
		returnCode := work.Body.PopUInt16()
		if returnCode != 0 {
			work.Finish()
			return
		}
		bs := work.Body.PopByteArray()
		accounts := &pbgo.AccountArray{}
		err := proto.Unmarshal(bs, accounts)
		if err != nil {
			logger.Error("Failed to unmarshal AccountArray, err: %+v", err)
			work.Finish()
			return
		}
		for _, account := range accounts.Accounts {
			logger.Debug("account: %+v", account)
			// 將用戶資訊加入緩存
			s.accounts[account.Account] = account
		}
	default:
		logger.Warn("Unsupport service: %d\n", service)
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

		// 將新註冊用戶加入緩存管理
		s.accounts[account.Account] = account
		logger.Info("New account created : %+v", account)

		// ==================================================
		// 準備將回應返還給 Main server
		// ==================================================
		td := base.NewTransData()
		td.AddByte(define.CommissionCommand)
		td.AddUInt16(define.Register)
		td.AddInt32(cid)
		td.AddUInt16(returnCode)

		// Account data for register
		clone := proto.Clone(account).(*pbgo.Account)
		clone.Password = ""
		bs, _ = proto.Marshal(clone)
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
