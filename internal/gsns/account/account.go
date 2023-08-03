package account

import (
	"fmt"
	"internal/agrt"
	"internal/define"
	"internal/pbgo"
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
	Token    uint64
}

// 與 Account 相關的由這個物件來管理
type AccountMgr struct {
	httpAnswer *ans.HttpAnser
	// key1: user id, key2: token
	users         *cntr.BikeyMap[int32, uint64, *pbgo.SnsUser]
	edges         map[int32]*cntr.Set[int32]
	logger        *glog.Logger
	heartbeatTime time.Time
}

func NewAccountMgr(lg *glog.Logger) *AccountMgr {
	m := &AccountMgr{
		users:         cntr.NewBikeyMap[int32, uint64, *pbgo.SnsUser](),
		edges:         make(map[int32]*cntr.Set[int32]),
		logger:        lg,
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
		m.logger.Error("Failed to unmarshal agreement, err: %+v", err)
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
			m.logger.Info("Heart response Now: %+v", time.Now())
			m.heartbeatTime = time.Now().Add(1 * time.Minute)
		}
		work.Finish()
	default:
		fmt.Printf("Unsupport service: %d\n", agreement.Service)
		work.Finish()
	}
}

func (m *AccountMgr) handleAccountCommission(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	case define.Register:
		work.Finish()
		// 利用 cid 取得對應的 Context
		c := m.httpAnswer.GetContext(agreement.Cid)
		m.logger.Debug("returnCode: %d", agreement.ReturnCode)

		if agreement.ReturnCode != 0 {
			c.Json(ghttp.StatusBadGateway, ghttp.H{
				"ret": 1,
				"msg": fmt.Sprintf("returnCode: %d", agreement.ReturnCode),
			})
		} else {
			c.Json(ghttp.StatusOK, ghttp.H{
				"ret": 0,
				"msg": fmt.Sprintf("registered account: %+v", agreement.Accounts[0]),
			})
		}
		m.httpAnswer.Send(c)

	case define.Login:
		work.Finish()

		// 取得空閒的 HTTP 連線物件
		c := m.httpAnswer.GetContext(agreement.Cid)
		defer m.httpAnswer.Send(c)
		m.logger.Debug("returnCode: %d", agreement.ReturnCode)

		if agreement.ReturnCode == 0 {
			account := agreement.Accounts[0]
			m.logger.Info("index: %d, name: %s, Account: %+v", account.Index, account.Account)
			user := &pbgo.SnsUser{
				Index: account.Index,
				Name:  account.Account,
				Info:  account.Info,
				Token: m.getToken(),
			}
			m.logger.Info("New user: %+v", user)
			if m.users.ContainKey1(account.Index) {
				m.users.UpdateByKey1(account.Index, cntr.NewBivalue(account.Index, user.Token, user))
			} else {
				err := m.users.Add(user.Index, user.Token, user)
				if err != nil {
					m.logger.Error("Failed to add user, err: %+v", err)
					c.Json(ghttp.StatusInternalServerError, ghttp.H{
						"msg":   "Login failed",
						"token": -1,
					})
					return
				}
			}
			c.Json(ghttp.StatusOK, ghttp.H{
				"msg":   fmt.Sprintf("User %s login success", account.Account),
				"token": user.Token,
				"info":  user.Info,
			})
		} else {
			m.logger.Error("Failed to login, ReturnCode: %d", agreement.ReturnCode)
			c.Json(ghttp.StatusInternalServerError, ghttp.H{
				"msg":   "Login failed",
				"token": -1,
			})
		}

	case define.SetUserData:
		work.Finish()

		// 取得空閒的 HTTP 連線物件
		c := m.httpAnswer.GetContext(agreement.Cid)
		m.logger.Debug("returnCode: %d", agreement.ReturnCode)

		if agreement.ReturnCode == 0 {
			account := agreement.Accounts[0]
			m.logger.Info("index: %d, name: %s", account.Index, account.Account)
			// 檢查緩存中是否存在
			user, ok := m.users.GetByKey1(account.Index)
			if !ok {
				c.Json(ghttp.StatusInternalServerError, ghttp.H{
					"err": fmt.Sprintf("Not found user %s in cache.", account.Account),
				})
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
		} else {
			c.Json(ghttp.StatusInternalServerError, ghttp.H{
				"err": fmt.Sprintf("Return code %d", agreement.ReturnCode),
			})
		}
		m.httpAnswer.Send(c)

	case define.GetOtherUsers:
		work.Finish()
		c := m.httpAnswer.GetContext(agreement.Cid)
		c.Json(ghttp.StatusOK, ghttp.H{
			"users": agreement.Accounts,
		})
		m.httpAnswer.Send(c)

	case define.Subscribe:
		// 取得空閒的 HTTP 連線物件
		c := m.httpAnswer.GetContext(agreement.Cid)
		m.logger.Debug("returnCode: %d", agreement.ReturnCode)

		if agreement.ReturnCode == 0 {
			edge := agreement.Edges[0]
			if _, ok := m.edges[edge.UserId]; !ok {
				m.edges[edge.UserId] = cntr.NewSet[int32]()
			}
			m.edges[edge.UserId].Add(edge.Target)
			c.Json(ghttp.StatusOK, ghttp.H{
				"ret": 0,
				"msg": fmt.Sprintf("User %d subscribe user %d", edge.UserId, edge.Target),
			})
		} else {
			c.Json(ghttp.StatusInternalServerError, ghttp.H{
				"err": fmt.Sprintf("Return code %d", agreement.ReturnCode),
			})
		}

		work.Finish()
		m.httpAnswer.Send(c)
	default:
	}
}
