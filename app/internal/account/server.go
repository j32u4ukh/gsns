package account

import (
	"fmt"
	"internal/agrt"
	"internal/define"
	"internal/pbgo"
	"time"

	"github.com/j32u4ukh/cntr"
	"github.com/j32u4ukh/gos"
	"github.com/j32u4ukh/gos/ans"
	"github.com/j32u4ukh/gos/base"
	"google.golang.org/protobuf/proto"
)

type AccountServer struct {
	Tcp *ans.Tcp0Anser
	// key1: user index, key2: account name;
	// key1 不可變更，但 key2 可以更新
	accounts *cntr.BikeyMap[int32, string, *pbgo.Account]
	// 社群關係
	Edges map[int32]*cntr.Set[int32]
	// key: server id, value: conn id
	serverIdDict  map[int32]int32
	heartbeatTime time.Time
}

func NewAccountServer() *AccountServer {
	s := &AccountServer{
		accounts:      cntr.NewBikeyMap[int32, string, *pbgo.Account](),
		Edges:         make(map[int32]*cntr.Set[int32]),
		serverIdDict:  make(map[int32]int32),
		heartbeatTime: time.Now(),
	}
	return s
}

func (s *AccountServer) Handler(work *base.Work) {
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
		s.handleSystem(work, agreement)
	case define.NormalCommand:
		s.handleNormal(work, agreement)
	case define.CommissionCommand:
		s.handleCommission(work, agreement)
	default:
		logger.Warn("Unsupport command: %d\n", agreement.Cmd)
		work.Finish()
	}
}

func (rrs *AccountServer) Run() {

}

func (s *AccountServer) handleSystem(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	// 回應心跳包
	case define.Heartbeat:
		agreement.Msg = "OK"
		_, err := agrt.SendWork(work, agreement)
		// TODO: CannotSendMessage
		if err != nil {
			logger.Error("Failed to send work, err: %+v", err)
		}
	case define.Introduction:
		// TODO: WrongConnectionIdentity
		if agreement.Cipher != define.CIPHER {
			logger.Error("Cipher: %s, Identity: %d", agreement.Cipher, agreement.Identity)
			gos.Disconnect(define.DbaPort, work.Index)
		} else {
			s.serverIdDict[agreement.Identity] = work.Index
			logger.Info("Hello %s from %d", define.ServerName(agreement.Identity), work.Index)
		}
		work.Finish()
	default:
		logger.Warn("Unsupport service: %d\n", agreement.Service)
		work.Finish()
	}
}

func (s *AccountServer) handleNormal(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	default:
		logger.Warn("Unsupport normal service: %d", agreement.Service)
		work.Finish()
	}
}

func (s *AccountServer) handleCommission(work *base.Work, agreement *agrt.Agreement) {
	logger.Info("Service: %d, Cid: %d", agreement.Service, agreement.Cid)

	switch agreement.Service {
	case define.Register:
		// 將註冊數據傳到 Dba 伺服器
		_, err := agrt.SendToServer(define.DbaServer, agreement)

		// TODO: CannotSendMessage
		if err != nil {
			agreement.ReturnCode = 1
			agreement.Msg = fmt.Sprintf("Failed to send to server: %d", define.DbaServer)
			_, err = agrt.SendWork(work, agreement)

			// TODO: CannotSendMessage
			if err != nil {
				logger.Error("Failed to send work, err: %+v", err)
				work.Finish()
			} else {
				logger.Info("Send define.Register response(%d): %+v", agreement.ReturnCode, agreement.Msg)
			}
		} else {
			logger.Info("Send define.Register request: %+v", agreement)
			work.Finish()
		}

	case define.Login:
		defer logger.Info("Check Login agreement: %+v", agreement)
		data := agreement.Accounts[0]
		var account *pbgo.Account
		var ok bool
		var err error

		// 檢查是否有用戶帳號緩存
		if account, ok = s.accounts.GetByKey2(data.Account); ok {
			logger.Info("Account in cache: %+v", account)

			// 檢查密碼是否正確
			if data.Password != account.Password {
				// Return code
				agreement.ReturnCode = 1
				agreement.Msg = fmt.Sprintf("Password %s is not correct.", data.Password)
				agreement.Accounts = agreement.Accounts[:0]
				logger.Error(agreement.Msg)
			} else {
				agreement.ReturnCode = 0
				agreement.Accounts[0] = proto.Clone(account).(*pbgo.Account)
				agreement.Accounts[0].Password = ""
				if _, ok := s.Edges[account.Index]; !ok {
					s.Edges[account.Index] = cntr.NewSet[int32]()
				}
				// 載入社群關係
				edges := s.Edges[account.Index]
				for edge := range edges.Elements {
					agreement.Edges = append(agreement.Edges, &pbgo.Edge{
						UserId: account.Index,
						Target: edge,
					})
				}
			}
			_, err = agrt.SendWork(work, agreement)

			// TODO: CannotSendMessage
			if err != nil {
				logger.Error("Failed to send work, err: %+v", err)
			} else {
				logger.Info("Send define.Login response: (%d) %+v", agreement.ReturnCode, agreement)
			}
			return
		}

		// 若不存在用戶帳號緩存
		logger.Info("不存在用戶帳號緩存")
		_, err = agrt.SendToServer(define.DbaServer, agreement)

		// TODO: CannotSendMessage
		if err != nil {
			agreement.ReturnCode = 2
			agreement.Msg = "Failed to send to Dba server"
			agreement.Accounts = agreement.Accounts[:0]
			_, err = agrt.SendWork(work, agreement)

			// TODO: CannotSendMessage
			if err != nil {
				logger.Error("Failed to send work, err: %+v", err)
				work.Finish()
			} else {
				logger.Info("Send define.Login response: (%d) %+v", agreement.ReturnCode, agreement)
			}
		} else {
			logger.Info("Send define.Login request: %+v", agreement)
			work.Finish()
		}

	// 設置用戶資料
	case define.SetUserData:
		var err error
		newAccount := agreement.Accounts[0]
		account, ok := s.accounts.GetByKey1(newAccount.Index)
		if !ok {
			agreement.ReturnCode = 1
			agreement.Msg = fmt.Sprintf("找不到 user(%d)", newAccount.Index)
			_, err = agrt.SendWork(work, agreement)
			// TODO: CannotSendMessage
			if err != nil {
				logger.Error("Failed to send work, err: %+v", err)
			} else {
				logger.Info("Send define.SetUserData response: (%d) %+v", agreement.ReturnCode, agreement)
			}
			return
		}

		// 填入原始密碼
		newAccount.Password = account.Password
		logger.Info("newAccount: %+v", newAccount)
		logger.Info("Accounts[0]: %+v", agreement.Accounts[0])

		// ==================================================
		// 更新緩存後，再將更新請求傳送給 DBA server
		// ==================================================
		_, err = agrt.SendToServer(define.DbaServer, agreement)

		// TODO: CannotSendMessage
		if err != nil {
			agreement.ReturnCode = 2
			agreement.Msg = fmt.Sprintf("Failed to send to %s", define.ServerName(define.DbaServer))
			_, err = agrt.SendWork(work, agreement)
			// TODO: CannotSendMessage
			if err != nil {
				logger.Error("Failed to send work, err: %+v", err)
				work.Finish()
			} else {
				logger.Info("Send define.SetUserData response: (%d) %+v", agreement.ReturnCode, agreement)
			}
		} else {
			logger.Info("Send define.SetUserData request: %+v", agreement)
			// 當前工作直接結束，無須回應
			work.Finish()
		}

	case define.GetOtherUsers:
		_, err := agrt.SendToServer(define.DbaServer, agreement)
		// TODO: CannotSendMessage
		if err != nil {
			agreement.ReturnCode = 1
			agreement.Msg = "Failed to send to Dba server."
			_, err = agrt.SendWork(work, agreement)

			// TODO: CannotSendMessage
			if err != nil {
				logger.Error("Failed to send work, err: %+v", err)
				work.Finish()
			} else {
				logger.Info("Send define.GetOtherUsers response: (%d) %+v", agreement.ReturnCode, agreement)
			}
		} else {
			logger.Info("Send define.GetOtherUsers request: %+v", agreement)
			work.Finish()
		}

	case define.Subscribe:
		_, err := agrt.SendToServer(define.DbaServer, agreement)
		// TODO: CannotSendMessage
		if err != nil {
			agreement.ReturnCode = 1
			agreement.Msg = "Failed to send to Dba server."
			_, err = agrt.SendWork(work, agreement)
			// TODO: CannotSendMessage
			if err != nil {
				logger.Error("Failed to send work, err: %+v", err)
				work.Finish()
			} else {
				logger.Info("Send define.Subscribe response: (%d) %+v", agreement.ReturnCode, agreement)
			}
		} else {
			logger.Info("Send define.Subscribe request: %+v", agreement)
			work.Finish()
		}
	default:
		fmt.Printf("Unsupport commission: %d", agreement.Service)
		work.Finish()
	}
}
