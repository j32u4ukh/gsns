package account

import (
	"fmt"
	"internal/define"
	"internal/pbgo"

	"github.com/j32u4ukh/glog/v2"
	"github.com/j32u4ukh/gos/ans"
	"github.com/j32u4ukh/gos/base"
	"github.com/j32u4ukh/gos/base/ghttp"
)

// 與 Account 相關的由這個物件來管理
type AccountMgr struct {
	httpAnswer *ans.HttpAnser
	users      map[int32]*pbgo.SnsUser
	logger     *glog.Logger
}

func NewAccountMgr(lg *glog.Logger) *AccountMgr {
	m := &AccountMgr{
		users:  map[int32]*pbgo.SnsUser{},
		logger: lg,
	}
	return m
}

func (m *AccountMgr) SetHttpAnswer(a *ans.HttpAnser) {
	m.httpAnswer = a
}

func (m *AccountMgr) WorkHandler(work *base.Work) {
	cmd := work.Body.PopByte()

	switch cmd {
	case define.SystemCommand:
		m.handleSystemCommand(work)
	case define.CommissionCommand:
		m.handleAccountCommission(work)
	default:
		fmt.Printf("Unsupport command: %d\n", cmd)
		work.Finish()
	}
}

func (m *AccountMgr) handleSystemCommand(work *base.Work) {
	service := work.Body.PopUInt16()

	switch service {
	case 0:
		response := work.Body.PopString()
		fmt.Printf("Heart beat response: %s\n", response)
		work.Finish()
	default:
		fmt.Printf("Unsupport service: %d\n", service)
		work.Finish()
	}
}

func (m *AccountMgr) handleAccountCommission(work *base.Work) {
	commission := work.Body.PopUInt16()
	switch commission {
	case define.Register:
		c := m.httpAnswer.GetContext(-1)
		c.Cid = work.Body.PopInt32()
		returnCode := work.Body.PopUInt16()
		m.logger.Debug("returnCode: %d", returnCode)
		work.Finish()

		c.Json(200, ghttp.H{
			"index": 5,
			"msg":   fmt.Sprintf("POST | register returnCode: %d", returnCode),
		})
		m.httpAnswer.Send(c)
	case define.Login:
		c := m.httpAnswer.GetContext(-1)
		c.Cid = work.Body.PopInt32()
		returnCode := work.Body.PopUInt16()
		token := work.Body.PopUInt64()
		m.logger.Debug("returnCode: %d", returnCode)
		work.Finish()

		if returnCode == 0 {
			c.Json(200, ghttp.H{
				"msg":   "Login success",
				"token": token,
			})
		} else {
			c.Json(200, ghttp.H{
				"msg":   "Login failed",
				"token": -1,
			})
		}

		m.httpAnswer.Send(c)
	default:
	}
}
