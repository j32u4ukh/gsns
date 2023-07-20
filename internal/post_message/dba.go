package pm

import (
	"fmt"
	"internal/agrt"
	"internal/define"
	"internal/pbgo"
	"time"

	"github.com/j32u4ukh/gos"
	"github.com/j32u4ukh/gos/base"
	"google.golang.org/protobuf/proto"
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
	case define.GetPost:
		if agreement.ReturnCode != 0 {
			logger.Error("Failed to query posts, err: %s", agreement.Msg)
		} else {
			for i, pm := range agreement.PostMessages {
				logger.Debug("%d) %+v", i, pm)
			}
		}
	default:
		logger.Warn("Unsupport service: %d\n", agreement.Service)
		work.Finish()
	}
}

func (s *PostMessageServer) handleDbaCommission(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	case define.AddPost:
		work.Finish()
		logger.Debug("#post: %d, returnCode: %d", len(agreement.PostMessages), agreement.ReturnCode)
		post := proto.Clone(agreement.PostMessages[0]).(*pbgo.PostMessage)
		logger.Info("New post: %+v", post)

		if post.ParentId == 0 {
			s.pmRoots[post.Id] = post
		} else {
			s.pmLeaves.Add(post.Id, post.ParentId, post)
		}

		td := base.NewTransData()
		bs, _ := agreement.Marshal()
		td.AddByteArray(bs)
		data := td.FormData()

		// 將註冊結果回傳主伺服器
		err := gos.SendToClient(define.PostMessagePort, s.MainServerId, &data, td.GetLength())

		if err != nil {
			logger.Error("Failed to send to client %d: %v\nError: %+v", s.MainServerId, data, err)
			return
		}
	default:
		fmt.Printf("Unsupport commission: %d\n", agreement.Service)
		work.Finish()
	}
}
