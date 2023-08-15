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
	// case define.GetPost:
	// 	work.Finish()
	// 	if agreement.ReturnCode != 0 {
	// 		logger.Error("ReturnCode: %d, err: %s", agreement.ReturnCode, agreement.Msg)
	// 	} else {
	// 		for i, pm := range agreement.PostMessages {
	// 			logger.Debug("%d) %+v", i, pm)
	// 			s.cachePost(pm)
	// 		}
	// 	}
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

		if agreement.ReturnCode == 0 {
			post := proto.Clone(agreement.PostMessages[0]).(*pbgo.PostMessage)
			logger.Info("New post: %+v", post)
			s.cachePost(post)
		}

		// td := base.NewTransData()
		// bs, err := agreement.Marshal()
		// if err != nil {
		// 	logger.Error("Failed to marshal agreement, err: %+v", err)
		// 	return
		// }
		// td.AddByteArray(bs)
		// data := td.FormData()

		// // 將註冊結果回傳主伺服器
		// err = gos.SendToClient(define.PostMessagePort, s.serverIdDict[define.GsnsServer], &data, int32(len(data)))

		_, err := agrt.SendToClient(define.PostMessagePort, s.serverIdDict[define.GsnsServer], agreement)

		if err != nil {
			logger.Error("Failed to send Gsns serve, err: %+v", err)
		} else {
			logger.Info("Send define.AddPost response(%d): %+v", agreement.ReturnCode, agreement)
		}
	case define.GetPost:
		work.Finish()
		logger.Info("Receive define.GetPost response(%d): %+v", agreement.ReturnCode, agreement)

		if agreement.ReturnCode == 0 {
			for _, pm := range agreement.PostMessages {
				s.cachePost(proto.Clone(pm).(*pbgo.PostMessage))
			}
		}

		// td := base.NewTransData()
		// bs, err := agreement.Marshal()
		// if err != nil {
		// 	logger.Error("Failed to marshal agreement, err: %+v", err)
		// 	return
		// }
		// td.AddByteArray(bs)
		// data := td.FormData()

		// // 將註冊結果回傳主伺服器
		// err = gos.SendToClient(define.PostMessagePort, s.serverIdDict[define.GsnsServer], &data, int32(len(data)))

		_, err := agrt.SendToClient(define.PostMessagePort, s.serverIdDict[define.GsnsServer], agreement)

		if err != nil {
			logger.Error("Failed to send to Gsns server, err: %+v", err)
			return
		} else {
			logger.Info("Send define.GetPost response(%d): %+v", agreement.ReturnCode, agreement)
		}
	case define.ModifyPost:
		work.Finish()
		logger.Info("Receive define.ModifyPost response(%d): %+v", agreement.ReturnCode, agreement)

		if agreement.ReturnCode == 0 {
			pm := proto.Clone(agreement.PostMessages[0]).(*pbgo.PostMessage)
			s.pmRoots[pm.Id] = pm
			s.postIds[pm.UserId].Add(pm.Id)
		}

		// td := base.NewTransData()
		// bs, err := agreement.Marshal()
		// if err != nil {
		// 	logger.Error("Failed to marshal agreement, err: %+v", err)
		// 	return
		// }
		// td.AddByteArray(bs)
		// data := td.FormData()

		// // 將註冊結果回傳主伺服器
		// err = gos.SendToClient(define.PostMessagePort, s.serverIdDict[define.GsnsServer], &data, int32(len(data)))

		_, err := agrt.SendToClient(define.PostMessagePort, s.serverIdDict[define.GsnsServer], agreement)

		if err != nil {
			logger.Error("Failed to send to Gsns server, err: %+v", err)
		} else {
			logger.Info("Send define.ModifyPost response(%d): %+v", agreement.ReturnCode, agreement)
		}
	case define.GetSubscribedPosts:
		work.Finish()
		// td := base.NewTransData()
		// bs, err := agreement.Marshal()
		// if err != nil {
		// 	logger.Error("Failed to marshal agreement, err: %+v", err)
		// 	return
		// }
		// td.AddByteArray(bs)
		// data := td.FormData()

		// // 將註冊結果回傳主伺服器
		// err = gos.SendToClient(define.PostMessagePort, s.serverIdDict[define.GsnsServer], &data, int32(len(data)))

		_, err := agrt.SendToClient(define.PostMessagePort, s.serverIdDict[define.GsnsServer], agreement)

		if err != nil {
			logger.Error("Failed to send to Gsns server, err: %+v", err)
		} else {
			logger.Info("Send define.GetSubscribedPosts response(%d): %+v", agreement.ReturnCode, agreement)
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
