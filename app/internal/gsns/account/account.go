package account

import (
	"fmt"
	"internal/agrt"
	"internal/define"
	"internal/pbgo"
	"internal/utils"
	"time"

	"github.com/j32u4ukh/cntr"
	"github.com/j32u4ukh/glog/v2"
	"github.com/j32u4ukh/gos/ans"
	"github.com/j32u4ukh/gos/base"
	"github.com/j32u4ukh/gos/base/ghttp"
)

type AccountProtocol struct {
	Account  string
	Password string
	Info     string
	Token    string
}

// 與 Account 相關的由這個物件來管理
type AccountMgr struct {
	httpAnswer *ans.HttpAnser
	// key1: user id, key2: token
	users         *cntr.BikeyMap[int32, string, *pbgo.User]
	Edges         map[int32]*cntr.Set[int32]
	serverLogger  *glog.Logger
	clientLogger  *glog.Logger
	heartbeatTime time.Time
}

func NewAccountMgr(slog, clog *glog.Logger) *AccountMgr {
	m := &AccountMgr{
		users:         cntr.NewBikeyMap[int32, string, *pbgo.User](),
		Edges:         make(map[int32]*cntr.Set[int32]),
		serverLogger:  slog,
		clientLogger:  clog,
		heartbeatTime: time.Now(),
	}
	return m
}

func (m *AccountMgr) SetHttpAnswer(a *ans.HttpAnser) {
	m.httpAnswer = a
}

func (m *AccountMgr) WorkHandler(work *base.Work) {
	agreement := agrt.GetAgreement()
	defer agrt.PutAgreement(agreement)
	err := agreement.Init(work)
	if err != nil {
		work.Finish()
		m.serverLogger.Error("Failed to unmarshal agreement, err: %+v", err)
		return
	}
	switch agreement.Cmd {
	case define.SystemCommand:
		m.handleSystem(work, agreement)
	case define.CommissionCommand:
		m.handleAccountCommission(work, agreement)
	default:
		fmt.Printf("Unsupport command: %d\n", agreement.Cmd)
		work.Finish()
	}
}

func (m *AccountMgr) handleSystem(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	case define.Heartbeat:
		if time.Now().After(m.heartbeatTime) {
			m.serverLogger.Info("Heart response Now: %+v", time.Now())
			m.heartbeatTime = time.Now().Add(1 * time.Minute)
		}
		work.Finish()
	default:
		fmt.Printf("Unsupport service: %d\n", agreement.Service)
		work.Finish()
	}
}

func (m *AccountMgr) handleAccountCommission(work *base.Work, agreement *agrt.Agreement) {
	work.Finish()
	switch agreement.Service {
	case define.Register:
		m.responseCommission(agreement, nil)

	case define.Login:
		m.responseCommission(agreement, func(c *ghttp.Context) {
			account := agreement.Accounts[0]
			m.serverLogger.Info("index: %d, name: %s, Account: %+v", account.Index, account.Account, account)
			user := &pbgo.User{
				Index: account.Index,
				Name:  account.Account,
				Info:  account.Info,
				Token: m.getToken(),
			}
			m.serverLogger.Info("New user: %+v", user)
			if m.users.ContainKey1(account.Index) {
				m.users.UpdateByKey1(account.Index, cntr.NewBivalue(account.Index, user.Token, user))
			} else {
				err := m.users.Add(user.Index, user.Token, user)
				if err != nil {
					msg := utils.JsonResponse(c, define.Error.InvalidInsertData, user)
					m.clientLogger.Error("%s, err: %+v", msg, err)
					return
				}
			}

			// 初始化用戶的 Edges
			if _, ok := m.Edges[account.Index]; !ok {
				m.Edges[account.Index] = cntr.NewSet[int32]()
			}

			// 寫入社群資訊
			for _, edge := range agreement.Edges {
				m.Edges[account.Index].Add(edge.Target)
			}

			c.Json(ghttp.StatusOK, ghttp.H{
				"ret":    define.Error.None,
				"msg":    fmt.Sprintf("User %s login success", account.Account),
				"token":  user.Token,
				"index":  user.Index,
				"n_edge": len(agreement.Edges),
			})
		})

	case define.SetUserData:
		m.responseCommission(agreement, func(c *ghttp.Context) {
			account := agreement.Accounts[0]
			m.serverLogger.Info("index: %d, name: %s", account.Index, account.Account)
			// 檢查緩存中是否存在
			user, ok := m.users.GetByKey1(account.Index)
			if !ok {
				m.clientLogger.Error(utils.JsonResponse(c, define.Error.NotFoundUser, "index", account.Index))
			} else {
				user.Name = account.Account
				user.Info = account.Info
				bivalue := cntr.NewBivalue(user.Index, user.Token, user)

				// 更新用戶緩存
				m.users.UpdateByKey1(user.Index, bivalue)

				c.Json(ghttp.StatusOK, ghttp.H{
					"msg":   fmt.Sprintf("User %s update success", account.Account),
					"token": user.Token,
					"info":  user.Info,
				})
			}
		})

	case define.GetOtherUsers:
		m.responseCommission(agreement, func(c *ghttp.Context) {
			ids := []int32{}
			for _, account := range agreement.Accounts {
				ids = append(ids, account.Index)
			}
			c.Json(ghttp.StatusOK, ghttp.H{
				"users": ids,
			})
		})

	case define.Subscribe:
		m.responseCommission(agreement, func(c *ghttp.Context) {
			edge := agreement.Edges[0]
			if _, ok := m.Edges[edge.UserId]; !ok {
				m.Edges[edge.UserId] = cntr.NewSet[int32]()
			}
			m.Edges[edge.UserId].Add(edge.Target)
			c.Json(ghttp.StatusOK, ghttp.H{
				"ret": 0,
				"msg": fmt.Sprintf("User %d subscribe user %d", edge.UserId, edge.Target),
			})
		})
	default:
	}
}

func (m *AccountMgr) responseCommission(agreement *agrt.Agreement, handlerFunc func(c *ghttp.Context)) {
	// 檢視收到的回應
	m.serverLogger.Info("agreement(%d): %+v", agreement.ReturnCode, agreement)
	// 利用 cid 取得對應的 Context
	c := m.httpAnswer.GetContext(agreement.Cid)
	if (agreement.ReturnCode == define.Error.None) && (handlerFunc != nil) {
		handlerFunc(c)
	} else {
		c.Json(define.GetStatus(agreement.ReturnCode), ghttp.H{
			"ret": agreement.ReturnCode,
			"msg": agreement.Msg,
		})
	}
	m.httpAnswer.Send(c)
}
