package pm

import (
	"fmt"
	"internal/agrt"
	"internal/define"

	"github.com/j32u4ukh/glog/v2"
	"github.com/j32u4ukh/gos/ans"
	"github.com/j32u4ukh/gos/base"
)

type PostMessageProtocol struct {
}

// 與 PostMessage 相關的由這個物件來管理
type PostMessageMgr struct {
	httpAnswer *ans.HttpAnser
	// // key1: user id, key2: token
	// users  *cntr.BikeyMap[int32, uint64, *pbgo.SnsUser]
	logger *glog.Logger
}

func NewPostMessageMgr(lg *glog.Logger) *PostMessageMgr {
	m := &PostMessageMgr{
		// users:  cntr.NewBikeyMap[int32, uint64, *pbgo.SnsUser](),
		logger: lg,
	}
	return m
}

func (m *PostMessageMgr) SetHttpAnswer(a *ans.HttpAnser) {
	m.httpAnswer = a
}

func (m *PostMessageMgr) WorkHandler(work *base.Work) {
	agreement := agrt.GetAgreement()
	defer agrt.PutAgreement(agreement)
	bs := work.Body.PopByteArray()
	err := agreement.Unmarshal(bs)
	if err != nil {
		work.Finish()
		m.logger.Error("Failed to unmarshal agreement, err: %+v", err)
		return
	}
	switch agreement.Cmd {
	case define.SystemCommand:
		m.handleSystemCommand(work, agreement)
	case define.CommissionCommand:
		m.handleAccountCommission(work, agreement)
	default:
		fmt.Printf("Unsupport command: %d\n", agreement.Cmd)
		work.Finish()
	}
}

func (m *PostMessageMgr) handleSystemCommand(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	case define.Heartbeat:
		fmt.Printf("Heart beat response: %s\n", agreement.Msg)
		work.Finish()
	default:
		fmt.Printf("Unsupport system service: %d\n", agreement.Service)
		work.Finish()
	}
}

func (m *PostMessageMgr) handleAccountCommission(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	default:
		fmt.Printf("Unsupport commission service: %d\n", agreement.Service)
		work.Finish()
	}
}
