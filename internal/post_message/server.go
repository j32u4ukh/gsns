package pm

import (
	"fmt"
	"internal/agrt"
	"internal/define"
	"time"

	"github.com/j32u4ukh/gos/ans"
	"github.com/j32u4ukh/gos/base"
)

type PostMessageServer struct {
	Tcp          *ans.Tcp0Anser
	MainServerId int32
	// key1: user index, key2: account name;
	// key1 不可變更，但 key2 可以更新
	// accounts *cntr.BikeyMap[int32, string, *pbgo.Account]
}

func NewPostMessageServer() *PostMessageServer {
	s := &PostMessageServer{}
	return s
}

func (s *PostMessageServer) Handler(work *base.Work) {
	agreement := agrt.GetAgreement()
	defer agrt.PutAgreement(agreement)
	bs := work.Body.PopByteArray()
	err := agreement.Unmarshal(bs)
	if err != nil {
		work.Finish()
		logger.Error("Failed to unmarshal agreement, err: %+v", err)
		return
	}
	switch agreement.Cmd {
	case define.SystemCommand:
		s.handleSystemCommand(work, agreement)
	case define.NormalCommand:
		s.handleNormalCommand(work, agreement)
	case define.CommissionCommand:
		s.handleCommissionCommand(work, agreement)
	default:
		logger.Warn("Unsupport command: %d\n", agreement.Cmd)
		work.Finish()
	}
}

func (s *PostMessageServer) Run() {

}

func (s *PostMessageServer) handleSystemCommand(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	// 回應心跳包
	case define.Heartbeat:
		logger.Debug("Heart beat! Now: %+v\n", time.Now())
		work.Body.Clear()
		agreement.Msg = "OK"
		bs, _ := agreement.Marshal()
		work.Body.AddByteArray(bs)
		work.SendTransData()
	default:
		logger.Warn("Unsupport system service: %d\n", agreement.Service)
		work.Finish()
	}
}

func (s *PostMessageServer) handleNormalCommand(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	default:
		logger.Warn("Unsupport normal service: %d", agreement.Service)
		work.Finish()
	}
}

func (s *PostMessageServer) handleCommissionCommand(work *base.Work, agreement *agrt.Agreement) {
	logger.Info("Service: %d, Cid: %d", agreement.Service, agreement.Cid)

	switch agreement.Service {
	default:
		fmt.Printf("Unsupport commission service: %d", agreement.Service)
		work.Finish()
	}
}
