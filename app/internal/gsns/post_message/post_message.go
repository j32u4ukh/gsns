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
	Token    string
	ParentId uint64 `json:"parent_id"`
	PostId   uint64 `json:"post_id"`
	Content  string
}

// 與 PostMessage 相關的由這個物件來管理
type PostMessageMgr struct {
	httpAnswer         *ans.HttpAnser
	logger             *glog.Logger
	getUserByTokenFunc func(token string) (*pbgo.User, bool)
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
	work.Finish()
	switch agreement.Service {
	case define.AddPost:
		m.responseCommission(agreement, func(c *ghttp.Context) {
			c.Json(ghttp.StatusOK, ghttp.H{
				"ret": define.Error.None,
				"msg": fmt.Sprintf("Post successfully: %+v", agreement.PostMessages[0]),
			})
		})

	case define.GetPost:
		m.responseCommission(agreement, func(c *ghttp.Context) {
			c.Json(ghttp.StatusOK, ghttp.H{
				"ret": define.Error.None,
				"pms": agreement.PostMessages,
			})
		})

	case define.GetMyPosts:
		m.responseCommission(agreement, func(c *ghttp.Context) {
			c.Json(ghttp.StatusOK, ghttp.H{
				"ret": define.Error.None,
				"pms": agreement.PostMessages,
			})
		})

	case define.ModifyPost:
		m.responseCommission(agreement, func(c *ghttp.Context) {
			c.Json(ghttp.StatusOK, ghttp.H{
				"ret": define.Error.None,
				"pm":  fmt.Sprintf("%+v", agreement.PostMessages[0]),
			})
		})

	case define.GetSubscribedPosts:
		m.responseCommission(agreement, func(c *ghttp.Context) {
			c.Json(ghttp.StatusOK, ghttp.H{
				"n_post": len(agreement.PostMessages),
				"posts":  agreement.PostMessages,
			})
		})

	default:
		fmt.Printf("Unsupport commission service: %d\n", agreement.Service)
	}
}

func (m *PostMessageMgr) SetFuncGetUserByToken(f func(token string) (*pbgo.User, bool)) {
	m.getUserByTokenFunc = f
}

func (m *PostMessageMgr) responseCommission(agreement *agrt.Agreement, handlerFunc func(c *ghttp.Context)) {
	// 檢視收到的回應
	m.logger.Info("agreement(%d): %+v", agreement.ReturnCode, agreement)
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
