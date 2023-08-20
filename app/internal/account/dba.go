package account

import (
	"fmt"
	"internal/agrt"
	"internal/define"
	"internal/pbgo"
	"time"

	"github.com/j32u4ukh/cntr"
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
	default:
		logger.Warn("Unsupport service: %d\n", agreement.Service)
		work.Finish()
	}
}

func (s *AccountServer) handleDbaCommission(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	case define.Register:
		work.Finish()

		if agreement.ReturnCode == define.Error.None {
			// 將新註冊用戶加入緩存管理
			account := proto.Clone(agreement.Accounts[0]).(*pbgo.Account)
			s.accounts.Set(account.Index, account.Account, account)
			logger.Info("New account created : %+v", account)
			// 隱藏密碼相關資訊，無須提供給 GSNS
			agreement.Accounts[0].Password = ""
		} else {
			logger.Info("ReturnCode: %d, Msg: %s", agreement.ReturnCode, agreement.Msg)
		}
		s.responseToGsns(agreement)

	case define.Login:
		work.Finish()
		logger.Info("Recieved login response: %+v", agreement)

		if agreement.ReturnCode == define.Error.None {
			// 將新註冊用戶加入緩存管理
			account := proto.Clone(agreement.Accounts[0]).(*pbgo.Account)
			if !s.accounts.ContainKey1(account.Index) {
				logger.Info("加入緩存")
				s.accounts.Add(account.Index, account.Account, account)
			} else {
				logger.Info("更新緩存")
				s.accounts.UpdateByKey1(account.Index, cntr.NewBivalue(account.Index, account.Account, account))
			}
			logger.Info("Login account: %+v", account)

			// 隱藏密碼相關資訊，無須提供給 GSNS
			agreement.Accounts[0].Password = ""

			// 載入社群關係
			if _, ok := s.Edges[account.Index]; !ok {
				s.Edges[account.Index] = cntr.NewSet[int32]()
			}

			for _, edge := range agreement.Edges {
				s.Edges[edge.UserId].Add(edge.Target)
			}
		} else {
			logger.Info("ReturnCode: %d, Msg: %s", agreement.ReturnCode, agreement.Msg)
		}
		s.responseToGsns(agreement)

	case define.SetUserData:
		work.Finish()

		if agreement.ReturnCode == define.Error.None {
			account := proto.Clone(agreement.Accounts[0]).(*pbgo.Account)
			logger.Debug("New account: %+v", account)

			// 更新用戶帳號緩存
			bivalue := cntr.NewBivalue(account.Index, account.Account, account)
			s.accounts.UpdateByKey1(account.Index, bivalue)

			// 隱藏密碼相關資訊，無須提供給 GSNS
			agreement.Accounts[0].Password = ""
		} else {
			logger.Info("ReturnCode: %d, Msg: %s", agreement.ReturnCode, agreement.Msg)
		}
		s.responseToGsns(agreement)

	case define.GetOtherUsers:
		work.Finish()
		if agreement.ReturnCode != define.Error.None {
			logger.Info("ReturnCode: %d, Msg: %s", agreement.ReturnCode, agreement.Msg)
		}
		s.responseToGsns(agreement)

	case define.Subscribe:
		work.Finish()
		if agreement.ReturnCode == 0 {
			// 更新社群關係緩存
			for _, edge := range agreement.Edges {
				s.Edges[edge.UserId].Add(edge.Target)
			}
		} else {
			logger.Info("ReturnCode: %d, Msg: %s", agreement.ReturnCode, agreement.Msg)
		}
		s.responseToGsns(agreement)
	default:
		fmt.Printf("Unsupport commission: %d\n", agreement.Service)
		work.Finish()
	}
}

func (s *AccountServer) responseToGsns(agreement *agrt.Agreement) {
	_, err := agrt.SendToClient(define.AccountPort, s.serverIdDict[define.GsnsServer], agreement)
	if err != nil {
		_, _, msg := define.ErrorMessage(define.Error.CannotSendMessage, "to Gsns server")
		logger.Error("%s, err: %+v", msg, err)
	} else {
		logger.Info("Send %s response(%d): %+v", define.ServiceName(agreement.Service), agreement.ReturnCode, agreement)
	}
}
