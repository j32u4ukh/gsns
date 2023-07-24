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
	switch agreement.Cmd {
	case define.SystemCommand:
		s.handleDbaSystem(work, agreement)
	case define.NormalCommand:
		s.handleDbaNormal(work, agreement)
	case define.CommissionCommand:
		s.handleDbaCommission(work, agreement)
	default:
		fmt.Printf("Unsupport command: %d\n", agreement.Cmd)
		work.Finish()
	}
}

func (s *PostMessageServer) handleDbaSystem(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	// 回應心跳包
	case define.Heartbeat:
		if time.Now().After(s.heartbeatTime) {
			logger.Info("Heart response Now: %+v", time.Now())
			s.heartbeatTime = time.Now().Add(1 * time.Minute)
		}
		work.Finish()
	default:
		fmt.Printf("Unsupport service: %d\n", agreement.Service)
		work.Finish()
	}
}

func (s *PostMessageServer) handleDbaNormal(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	case define.GetPost:
		work.Finish()
		if agreement.ReturnCode != 0 {
			logger.Error("ReturnCode: %d, err: %s", agreement.ReturnCode, agreement.Msg)
		} else {
			for i, pm := range agreement.PostMessages {
				logger.Debug("%d) %+v", i, pm)
				s.pmRoots[pm.Id] = pm
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
		err := gos.SendToClient(define.PostMessagePort, s.serverIdDict[define.GsnsServer], &data, td.GetLength())

		if err != nil {
			logger.Error("Failed to send to client %d: %v\nError: %+v", s.serverIdDict[define.GsnsServer], data, err)
			return
		}
	case define.GetPost:
		work.Finish()

		if agreement.ReturnCode == 0 {
			pm := agreement.PostMessages[0]
			s.pmRoots[pm.Id] = proto.Clone(pm).(*pbgo.PostMessage)
		} else {
			logger.Info("ReturnCode: %d, Msg: %s", agreement.ReturnCode, agreement.Msg)
		}

		logger.Info("Response agreement: %+v", agreement)
		td := base.NewTransData()
		bs, _ := agreement.Marshal()
		td.AddByteArray(bs)
		data := td.FormData()

		// 將註冊結果回傳主伺服器
		err := gos.SendToClient(define.PostMessagePort, s.serverIdDict[define.GsnsServer], &data, td.GetLength())

		if err != nil {
			logger.Error("Failed to send to client %d: %v\nError: %+v", s.serverIdDict[define.GsnsServer], data, err)
			return
		}
	default:
		fmt.Printf("Unsupport commission: %d\n", agreement.Service)
		work.Finish()
	}
}
