package account

import (
	"fmt"
	"internal/agrt"
	"internal/define"
	"internal/pbgo"
	"time"

	"github.com/j32u4ukh/cntr"
	"github.com/j32u4ukh/gos"
	"github.com/j32u4ukh/gos/base"
	"google.golang.org/protobuf/proto"
)

func (s *AccountServer) DbaHandler(work *base.Work) {
	agreement := agrt.GetAgreement()
	defer agrt.PutAgreement(agreement)
	err := agreement.Init(work)
	if err != nil {
		work.Finish()
		logger.Error("Failed to unmarshal agreement, err: %+v", err)
		return
	}
	logger.Info("Cmd: %d, Service: %d", agreement.Cmd, agreement.Service)
	switch byte(agreement.Cmd) {
	case define.SystemCommand:
		s.handleDbaSystemCommand(work, agreement)
	case define.NormalCommand:
		s.handleDbaNormalCommand(work, agreement)
	case define.CommissionCommand:
		s.handleDbaCommission(work, agreement)
	default:
		fmt.Printf("Unsupport command: %d\n", agreement.Cmd)
		work.Finish()
	}
}

func (s *AccountServer) handleDbaSystemCommand(work *base.Work, agreement *agrt.Agreement) {
	switch uint16(agreement.Service) {
	// 回應心跳包
	case define.Heartbeat:
		fmt.Printf("Heart response Now: %+v\n", time.Now())
		work.Finish()
	default:
		fmt.Printf("Unsupport service: %d\n", agreement.Service)
		work.Finish()
	}
}

func (s *AccountServer) handleDbaNormalCommand(work *base.Work, agreement *agrt.Agreement) {
	switch uint16(agreement.Service) {
	// 取得用戶資訊
	case define.GetUserData:
		logger.Debug("GetUserData")
		if agreement.ReturnCode != 0 {
			work.Finish()
			return
		}
		for _, account := range agreement.Accounts {
			logger.Debug("account: %+v", account)
			// 將用戶資訊加入緩存
			s.accounts.Set(account.Index, account.Account, account)
		}
	default:
		logger.Warn("Unsupport service: %d\n", agreement.Service)
		work.Finish()
	}
}

func (s *AccountServer) handleDbaCommission(work *base.Work, agreement *agrt.Agreement) {
	switch uint16(agreement.Service) {
	case define.Register:
		work.Finish()

		// 將新註冊用戶加入緩存管理
		account := proto.Clone(agreement.Accounts[0]).(*pbgo.Account)
		s.accounts.Set(account.Index, account.Account, account)
		logger.Info("New account created : %+v", account)

		td := base.NewTransData()
		td.AddByte(byte(agreement.Cmd))
		td.AddUInt16(uint16(agreement.Service))
		td.AddInt32(agreement.Cid)
		td.AddUInt16(uint16(agreement.ReturnCode))

		if agreement.ReturnCode != 0 {
			logger.Error("ReturnCode: %d", agreement.ReturnCode)
		} else {
			// 隱藏密碼相關資訊，無須提供給 GSNS
			agreement.Accounts[0].Password = ""
			bs, _ := proto.Marshal(agreement.Accounts[0])
			td.AddByteArray(bs)
		}

		data := td.FormData()

		// 將註冊結果回傳主伺服器
		err := gos.SendToClient(define.AccountPort, s.MainServerId, &data, td.GetLength())

		if err != nil {
			logger.Error("Failed to send to client %d: %v\nError: %+v", s.MainServerId, data, err)
			return
		}

	case define.SetUserData:
		var err error
		work.Finish()
		td := base.NewTransData()
		td.AddByte(byte(agreement.Cmd))
		td.AddUInt16(uint16(agreement.Service))
		td.AddInt32(agreement.Cid)
		td.AddUInt16(uint16(agreement.ReturnCode))

		if agreement.ReturnCode != 0 {
			logger.Error("ReturnCode: %d", agreement.ReturnCode)
		} else {
			account := agreement.Accounts[0]
			logger.Debug("New account: %+v", account)

			// 隱藏密碼相關資訊，無須提供給 GSNS
			clone := proto.Clone(account).(*pbgo.Account)
			clone.Password = ""
			bs, _ := proto.Marshal(clone)
			td.AddByteArray(bs)

			// 更新用戶帳號緩存
			bivalue := cntr.NewBivalue(account.Index, account.Account, account)
			s.accounts.UpdateByKey1(account.Index, bivalue)
		}
		data := td.FormData()

		// 將註冊結果回傳主伺服器
		err = gos.SendToClient(define.AccountPort, s.MainServerId, &data, td.GetLength())

		if err != nil {
			logger.Error("Failed to send to client %d: %v\nError: %+v", s.MainServerId, data, err)
			return
		}

	default:
		fmt.Printf("Unsupport commission: %d\n", agreement.Service)
		work.Finish()
	}
}
