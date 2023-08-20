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
		serverLogger.Error("Failed to unmarshal agreement, err: %+v", err)
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
		clientLogger.Warn("Unsupport command: %d\n", agreement.Cmd)
		work.Finish()
	}
}

func (s *AccountServer) handleSystem(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	// 回應心跳包
	case define.Heartbeat:
		agreement.ReturnCode = 0
		agreement.Msg = "OK"
		_, err := agrt.SendWork(work, agreement)
		if err != nil {
			_, _, agreement.Msg = define.ErrorMessage(define.Error.CannotSendMessage, "work")
			serverLogger.Error("%s, err: %+v", agreement.Msg, err)
		}
	case define.Introduction:
		if agreement.Cipher != define.CIPHER {
			_, _, agreement.Msg = define.ErrorMessage(define.Error.WrongConnectionIdentity, agreement.Cipher, agreement.Identity)
			clientLogger.Error(agreement.Msg)
			gos.Disconnect(define.DbaPort, work.Index)
		} else {
			s.serverIdDict[agreement.Identity] = work.Index
			serverLogger.Info("Hello %s from %d", define.ServerName(agreement.Identity), work.Index)
		}
		work.Finish()
	default:
		clientLogger.Warn("Unsupport service: %d", agreement.Service)
		work.Finish()
	}
}

func (s *AccountServer) handleNormal(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	default:
		clientLogger.Warn("Unsupport normal service: %d", agreement.Service)
		work.Finish()
	}
}

func (s *AccountServer) handleCommission(work *base.Work, agreement *agrt.Agreement) {
	serverLogger.Info("Service: %d, Cid: %d", agreement.Service, agreement.Cid)

	switch agreement.Service {
	case define.Register:
		s.handleCommissionRequest(work, agreement)

	case define.Login:
		defer serverLogger.Info("Check Login agreement: %+v", agreement)
		data := agreement.Accounts[0]
		var account *pbgo.Account
		var ok bool

		// 檢查是否有用戶帳號緩存
		if account, ok = s.accounts.GetByKey2(data.Account); ok {
			serverLogger.Info("Account in cache: %+v", account)

			// 檢查密碼是否正確
			if data.Password != account.Password {
				_, agreement.ReturnCode, agreement.Msg = define.ErrorMessage(define.Error.WrongParameter, "password", data.Password)
				agreement.Accounts = agreement.Accounts[:0]
				clientLogger.Error(agreement.Msg)
			} else {
				agreement.ReturnCode = define.Error.None
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

			s.handleRequest(work, agreement)
			// _, err = agrt.SendWork(work, agreement)
			// if err != nil {
			// 	_, _, agreement.Msg = define.ErrorMessage(define.Error.CannotSendMessage, "work")
			// 	serverLogger.Error("%s, err: %+v", agreement.Msg, err)
			// } else {
			// 	serverLogger.Info("Send define.Login response: (%d) %+v", agreement.ReturnCode, agreement)
			// }
		} else {
			// 若不存在用戶帳號緩存
			serverLogger.Info("不存在用戶帳號緩存")
			s.handleCommissionRequest(work, agreement)
		}

	// 設置用戶資料
	case define.SetUserData:
		newAccount := agreement.Accounts[0]
		account, ok := s.accounts.GetByKey1(newAccount.Index)
		if !ok {
			_, agreement.ReturnCode, agreement.Msg = define.ErrorMessage(define.Error.NotFoundUser, "index", newAccount.Index)
			s.handleRequest(work, agreement)

			// _, err = agrt.SendWork(work, agreement)
			// if err != nil {
			// 	_, _, agreement.Msg = define.ErrorMessage(define.Error.CannotSendMessage, "work")
			// 	serverLogger.Error("%s, err: %+v", agreement.Msg, err)
			// } else {
			// 	serverLogger.Info("Send define.SetUserData response: (%d) %+v", agreement.ReturnCode, agreement)
			// }
			return
		}

		// 填入原始密碼
		newAccount.Password = account.Password
		serverLogger.Info("newAccount: %+v", newAccount)
		serverLogger.Info("Accounts[0]: %+v", agreement.Accounts[0])

		// ==================================================
		// 更新緩存後，再將更新請求傳送給 DBA server
		// ==================================================
		s.handleCommissionRequest(work, agreement)

	case define.GetOtherUsers:
		s.handleCommissionRequest(work, agreement)

	case define.Subscribe:
		s.handleCommissionRequest(work, agreement)

	default:
		fmt.Printf("Unsupport commission: %d", agreement.Service)
		work.Finish()
	}
}

func (s *AccountServer) handleRequest(work *base.Work, agreement *agrt.Agreement) {
	_, err := agrt.SendWork(work, agreement)
	if err != nil {
		_, _, agreement.Msg = define.ErrorMessage(define.Error.CannotSendMessage, "work")
		serverLogger.Error("%s, err: %+v", agreement.Msg, err)
	} else {
		serverLogger.Info("Send %s request: %+v", define.ServiceName(agreement.Service), agreement)
	}
}

func (s *AccountServer) handleCommissionRequest(work *base.Work, agreement *agrt.Agreement) {
	_, err := agrt.SendToServer(define.DbaServer, agreement)
	if err != nil {
		_, agreement.ReturnCode, agreement.Msg = define.ErrorMessage(define.Error.CannotSendMessage, "to Dba server")
		serverLogger.Error("%s, err: %+v", agreement.Msg, err)
		_, err = agrt.SendWork(work, agreement)
		if err != nil {
			_, _, agreement.Msg = define.ErrorMessage(define.Error.CannotSendMessage, "work")
			serverLogger.Error("%s, err: %+v", agreement.Msg, err)
			work.Finish()
		} else {
			serverLogger.Info("Send %s response(%d): %+v", define.ServiceName(agreement.Service), agreement.ReturnCode, agreement.Msg)
		}
	} else {
		serverLogger.Info("Send %s request: %+v", define.ServiceName(agreement.Service), agreement)
		work.Finish()
	}
}
