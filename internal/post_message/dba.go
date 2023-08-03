package pm

import (
	"fmt"
	"internal/agrt"
	"internal/define"
	"internal/pbgo"
	"time"

	"github.com/j32u4ukh/cntr"
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
				s.cachePost(pm)
			}
		}
	// 用戶登入後，取得貼文數據並緩存下來
	case define.GetMyPosts:
		work.Finish()
		for i, pm := range agreement.PostMessages {
			logger.Debug("%d) %+v", i, pm)
			s.cachePost(pm)
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
		s.cachePost(post)

		td := base.NewTransData()
		bs, _ := agreement.Marshal()
		td.AddByteArray(bs)
		data := td.FormData()

		// 將註冊結果回傳主伺服器
		err := gos.SendToClient(define.PostMessagePort, s.serverIdDict[define.GsnsServer], &data, int32(len(data)))

		if err != nil {
			logger.Error("Failed to send to client %d: %v\nError: %+v", s.serverIdDict[define.GsnsServer], data, err)
			return
		}
	case define.GetPost:
		work.Finish()
		if agreement.ReturnCode == 0 {
			pm := proto.Clone(agreement.PostMessages[0]).(*pbgo.PostMessage)
			s.pmRoots[pm.Id] = pm
			s.postIds[pm.UserId].Add(pm.Id)
		} else {
			logger.Info("ReturnCode: %d, Msg: %s", agreement.ReturnCode, agreement.Msg)
		}
		td := base.NewTransData()
		bs, _ := agreement.Marshal()
		td.AddByteArray(bs)
		data := td.FormData()

		// 將註冊結果回傳主伺服器
		err := gos.SendToClient(define.PostMessagePort, s.serverIdDict[define.GsnsServer], &data, int32(len(data)))

		if err != nil {
			logger.Error("Failed to send to client %d: %v\nError: %+v", s.serverIdDict[define.GsnsServer], data, err)
			return
		}
	case define.ModifyPost:
		work.Finish()
		if agreement.ReturnCode == 0 {
			pm := proto.Clone(agreement.PostMessages[0]).(*pbgo.PostMessage)
			s.pmRoots[pm.Id] = pm
			s.postIds[pm.UserId].Add(pm.Id)
		} else {
			logger.Info("ReturnCode: %d, Msg: %s", agreement.ReturnCode, agreement.Msg)
		}
		td := base.NewTransData()
		bs, _ := agreement.Marshal()
		td.AddByteArray(bs)
		data := td.FormData()

		// 將註冊結果回傳主伺服器
		err := gos.SendToClient(define.PostMessagePort, s.serverIdDict[define.GsnsServer], &data, int32(len(data)))

		if err != nil {
			logger.Error("Failed to send to client %d: %v\nError: %+v", s.serverIdDict[define.GsnsServer], data, err)
			return
		}
	default:
		fmt.Printf("Unsupport commission: %d\n", agreement.Service)
		work.Finish()
	}
}

func (s *PostMessageServer) cachePost(pm *pbgo.PostMessage) {
	logger.Info("Cache post: %+v", pm)
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
