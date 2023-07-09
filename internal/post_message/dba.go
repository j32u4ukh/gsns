package pm

import (
	"fmt"
	"internal/agrt"
	"internal/define"
	"time"

	"github.com/j32u4ukh/gos/base"
)

func (s *PostMessageServer) DbaHandler(work *base.Work) {
	agreement := agrt.GetAgreement()
	defer agrt.PutAgreement(agreement)
	err := agreement.Init(work)
	if err != nil {
		work.Finish()
		logger.Error("Failed to unmarshal agreement, err: %+v", err)
		return
	}
	logger.Info("Cmd: %d, Service: %d", agreement.Cmd, agreement.Service)
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

func (s *PostMessageServer) handleDbaSystemCommand(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	// 回應心跳包
	case define.Heartbeat:
		fmt.Printf("Heart response Now: %+v\n", time.Now())
		work.Finish()
	default:
		fmt.Printf("Unsupport service: %d\n", agreement.Service)
		work.Finish()
	}
}

func (s *PostMessageServer) handleDbaNormalCommand(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	default:
		logger.Warn("Unsupport service: %d\n", agreement.Service)
		work.Finish()
	}
}

func (s *PostMessageServer) handleDbaCommission(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	default:
		fmt.Printf("Unsupport commission: %d\n", agreement.Service)
		work.Finish()
	}
}
