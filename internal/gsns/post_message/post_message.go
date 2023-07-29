package pm

import (
	"fmt"
	"internal/agrt"
	"internal/define"
	"internal/pbgo"
	"time"

	"github.com/j32u4ukh/glog/v2"
	"github.com/j32u4ukh/gos/ans"
	"github.com/j32u4ukh/gos/base"
	"github.com/j32u4ukh/gos/base/ghttp"
)

type PostMessageProtocol struct {
	Token    uint64
	ParentId uint64
	PostId   uint64
	Content  string
}

// 與 PostMessage 相關的由這個物件來管理
type PostMessageMgr struct {
	httpAnswer         *ans.HttpAnser
	logger             *glog.Logger
	getUserByTokenFunc func(token uint64) (*pbgo.SnsUser, bool)
	heartbeatTime      time.Time
}

func NewPostMessageMgr(lg *glog.Logger) *PostMessageMgr {
	m := &PostMessageMgr{
		heartbeatTime: time.Now(),
		logger:        lg,
	}
	return m
}

func (m *PostMessageMgr) SetHttpAnswer(a *ans.HttpAnser) {
	m.httpAnswer = a
}

func (m *PostMessageMgr) WorkHandler(work *base.Work) {
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
	case define.NormalCommand:
		m.handleNormal(work, agreement)
	case define.CommissionCommand:
		m.handleCommission(work, agreement)
	default:
		fmt.Printf("Unsupport command: %d\n", agreement.Cmd)
		work.Finish()
	}
}

func (m *PostMessageMgr) handleSystem(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	case define.Heartbeat:
		if time.Now().After(m.heartbeatTime) {
			m.logger.Info("Heart response Now: %+v", time.Now())
			m.heartbeatTime = time.Now().Add(1 * time.Minute)
		}
		work.Finish()
	default:
		fmt.Printf("Unsupport system service: %d\n", agreement.Service)
		work.Finish()
	}
}

func (m *PostMessageMgr) handleNormal(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	default:
		fmt.Printf("Unsupport system service: %d\n", agreement.Service)
		work.Finish()
	}
}

func (m *PostMessageMgr) handleCommission(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	case define.AddPost:
		// 利用 cid 取得對應的 Context
		c := m.httpAnswer.GetContext(agreement.Cid)
		m.logger.Debug("returnCode: %d", agreement.ReturnCode)

		if agreement.ReturnCode != 0 {
			c.Json(ghttp.StatusBadGateway, ghttp.H{
				"ret": agreement.ReturnCode,
				"msg": agreement.Msg,
			})
		} else {
			c.Json(ghttp.StatusOK, ghttp.H{
				"ret": 0,
				"msg": fmt.Sprintf("Post successfully: %+v", agreement.PostMessages[0]),
			})
		}
		m.httpAnswer.Send(c)
	case define.GetPost:
		// 利用 cid 取得對應的 Context
		c := m.httpAnswer.GetContext(agreement.Cid)
		m.logger.Debug("returnCode: %d", agreement.ReturnCode)

		if agreement.ReturnCode != 0 {
			c.Json(ghttp.StatusBadGateway, ghttp.H{
				"ret": agreement.ReturnCode,
				"msg": agreement.Msg,
			})
		} else {
			c.Json(ghttp.StatusOK, ghttp.H{
				"ret": 0,
				"pm":  fmt.Sprintf("%+v", agreement.PostMessages[0]),
			})
		}
		m.httpAnswer.Send(c)
		work.Finish()
	case define.ModifyPost:
		// 利用 cid 取得對應的 Context
		c := m.httpAnswer.GetContext(agreement.Cid)
		m.logger.Debug("returnCode: %d", agreement.ReturnCode)

		if agreement.ReturnCode != 0 {
			c.Json(ghttp.StatusBadGateway, ghttp.H{
				"ret": agreement.ReturnCode,
				"msg": agreement.Msg,
			})
		} else {
			c.Json(ghttp.StatusOK, ghttp.H{
				"ret": 0,
				"pm":  fmt.Sprintf("%+v", agreement.PostMessages[0]),
			})
		}
		m.httpAnswer.Send(c)
		work.Finish()
	default:
		fmt.Printf("Unsupport commission service: %d\n", agreement.Service)
		work.Finish()
	}
}

func (m *PostMessageMgr) SetFuncGetUserByToken(f func(token uint64) (*pbgo.SnsUser, bool)) {
	m.getUserByTokenFunc = f
}
