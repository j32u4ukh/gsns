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
	switch agreement.Cmd {
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
	switch agreement.Service {
	// 回應心跳包
	case define.Heartbeat:
		if time.Now().After(s.heartbeatTime) {
			logger.Info("Heart response Now: %+v", time.Now())
			s.heartbeatTime = time.Now().Add(1 * time.Minute)
		}
		work.Finish()
	default:
		fmt.Printf("Unsupport service: %d\n", agreement.Service)
		work.Finish()
	}
}

func (s *AccountServer) handleDbaNormalCommand(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	// 取得用戶資訊
	case define.GetUserData:
		logger.Debug("GetUserData")
		if agreement.ReturnCode == 0 {
			for _, account := range agreement.Accounts {
				logger.Debug("account: %+v", account)
				// 將用戶資訊加入緩存
				s.accounts.Set(account.Index, account.Account, account)
			}
		}
		work.Finish()

	default:
		logger.Warn("Unsupport service: %d\n", agreement.Service)
		work.Finish()
	}
}

func (s *AccountServer) handleDbaCommission(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	case define.Register:
		work.Finish()

		// 將新註冊用戶加入緩存管理
		account := proto.Clone(agreement.Accounts[0]).(*pbgo.Account)
		s.accounts.Set(account.Index, account.Account, account)
		logger.Info("New account created : %+v", account)
		// 隱藏密碼相關資訊，無須提供給 GSNS
		agreement.Accounts[0].Password = ""

		td := base.NewTransData()
		bs, _ := agreement.Marshal()
		td.AddByteArray(bs)
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

		if agreement.ReturnCode != 0 {
			logger.Error("ReturnCode: %d", agreement.ReturnCode)
			agreement.Accounts = agreement.Accounts[:0]
		} else {
			account := proto.Clone(agreement.Accounts[0]).(*pbgo.Account)
			logger.Debug("New account: %+v", account)

			// 更新用戶帳號緩存
			bivalue := cntr.NewBivalue(account.Index, account.Account, account)
			s.accounts.UpdateByKey1(account.Index, bivalue)

			// 隱藏密碼相關資訊，無須提供給 GSNS
			agreement.Accounts[0].Password = ""
		}

		bs, _ := agreement.Marshal()
		td.AddByteArray(bs)
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
