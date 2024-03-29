package pm

import (
	"fmt"
	"internal/agrt"
	"internal/define"
	"internal/pbgo"
	"time"

	"github.com/j32u4ukh/cntr"
	"github.com/j32u4ukh/gos/base"
	"google.golang.org/protobuf/proto"
)

func (s *PostMessageServer) DbaHandler(work *base.Work) {
	agreement := agrt.GetAgreement()
	defer agrt.PutAgreement(agreement)
	err := agreement.Init(work)
	if err != nil {
		work.Finish()
		serverLogger.Error("Failed to unmarshal agreement, err: %+v", err)
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
			serverLogger.Info("Heart response Now: %+v", time.Now())
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
	// 用戶登入後，取得貼文數據並緩存下來
	case define.GetMyPosts:
		work.Finish()
		for i, pm := range agreement.PostMessages {
			serverLogger.Debug("%d) %+v", i, pm)
			s.cachePost(pm)
		}
	default:
		clientLogger.Warn("Unsupport service: %d", agreement.Service)
		work.Finish()
	}
}

func (s *PostMessageServer) handleDbaCommission(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	case define.AddPost:
		work.Finish()

		if agreement.ReturnCode == define.Error.None {
			post := proto.Clone(agreement.PostMessages[0]).(*pbgo.PostMessage)
			serverLogger.Info("New post: %+v", post)
			s.cachePost(post)
		} else {
			serverLogger.Info("ReturnCode: %d, Msg: %s", agreement.ReturnCode, agreement.Msg)
		}
		s.responseToGsns(agreement)
	case define.GetPost:
		work.Finish()
		serverLogger.Info("Receive define.GetPost response(%d): %+v", agreement.ReturnCode, agreement)

		if agreement.ReturnCode == define.Error.None {
			for _, pm := range agreement.PostMessages {
				s.cachePost(proto.Clone(pm).(*pbgo.PostMessage))
			}
		} else {
			serverLogger.Info("ReturnCode: %d, Msg: %s", agreement.ReturnCode, agreement.Msg)
		}
		s.responseToGsns(agreement)
	case define.ModifyPost:
		work.Finish()
		serverLogger.Info("Receive define.ModifyPost response(%d): %+v", agreement.ReturnCode, agreement)

		if agreement.ReturnCode == 0 {
			pm := proto.Clone(agreement.PostMessages[0]).(*pbgo.PostMessage)
			s.pmRoots[pm.Id] = pm
			s.postIds[pm.UserId].Add(pm.Id)
		}
		s.responseToGsns(agreement)
	case define.GetSubscribedPosts:
		work.Finish()
		s.responseToGsns(agreement)
	default:
		fmt.Printf("Unsupport commission: %d\n", agreement.Service)
		work.Finish()
	}
}

func (s *PostMessageServer) cachePost(pm *pbgo.PostMessage) {
	serverLogger.Info("Cache post: %+v", pm)
	s.pmRoots[pm.Id] = pm
	if pm.ParentId != 0 {
		if _, ok := s.pmLeaves[pm.ParentId]; !ok {
			s.pmLeaves[pm.ParentId] = []*pbgo.PostMessage{}
		}
		s.pmLeaves[pm.ParentId] = append(s.pmLeaves[pm.ParentId], pm)
	}
	if _, ok := s.postIds[pm.UserId]; !ok {
		s.postIds[pm.UserId] = cntr.NewSet[uint64]()
	}
	s.postIds[pm.UserId].Add(pm.Id)
}

func (s *PostMessageServer) responseToGsns(agreement *agrt.Agreement) {
	_, err := agrt.SendToClient(define.PostMessagePort, s.serverIdDict[define.GsnsServer], agreement)
	if err != nil {
		_, _, msg := define.ErrorMessage(define.Error.CannotSendMessage, "to Gsns server")
		serverLogger.Error("%s, err: %+v", msg, err)
	} else {
		serverLogger.Info("Send %s response(%d): %+v", define.ServiceName(agreement.Service), agreement.ReturnCode, agreement)
	}
}
