package pm

import (
	"fmt"
	"internal/agrt"
	"internal/define"
	"internal/pbgo"
	"time"

	"github.com/j32u4ukh/cntr"
	"github.com/j32u4ukh/gos"
	"github.com/j32u4ukh/gos/ans"
	"github.com/j32u4ukh/gos/base"
	"google.golang.org/protobuf/proto"
)

// TODO: 不緩存別人的貼文，因為不知道是否完整，每次讀取別人緩存時，都會再問一次 DBA，沒有在此緩存別人貼文的必要。
type PostMessageServer struct {
	Tcp *ans.Tcp0Anser
	// key: user id; value: post ids
	postIds map[int32]*cntr.Set[uint64]

	// pmRoots 以 post id 為鍵值來管理貼文，包含所有貼文緩存，也包含 pmLeaves 中的所有貼文
	// pmLeaves 以 parent id 為鍵值來管理貼文，只包含有 parent id 的貼文緩存
	// key: post id; value: PostMessage
	pmRoots map[uint64]*pbgo.PostMessage
	// 回覆他人的貼文，parent id 為被回覆的貼文的 post id
	// key1: parent id, key2: post ids
	pmLeaves map[uint64][]*pbgo.PostMessage

	// key: server id, value: conn id
	serverIdDict  map[int32]int32
	heartbeatTime time.Time
}

func NewPostMessageServer() *PostMessageServer {
	s := &PostMessageServer{
		postIds:       make(map[int32]*cntr.Set[uint64]),
		pmRoots:       make(map[uint64]*pbgo.PostMessage),
		pmLeaves:      make(map[uint64][]*pbgo.PostMessage),
		serverIdDict:  make(map[int32]int32),
		heartbeatTime: time.Now(),
	}
	return s
}

func (s *PostMessageServer) Handler(work *base.Work) {
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
		s.handleSystem(work, agreement)
	case define.NormalCommand:
		s.handleNormal(work, agreement)
	case define.CommissionCommand:
		s.handleCommission(work, agreement)
	default:
		logger.Warn("Unsupport command: %d\n", agreement.Cmd)
		work.Finish()
	}
}

func (s *PostMessageServer) Run() {

}

func (s *PostMessageServer) handleSystem(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	// 回應心跳包
	case define.Heartbeat:
		work.Body.Clear()
		agreement.Msg = "OK"
		// bs, _ := agreement.Marshal()
		// work.Body.AddByteArray(bs)
		// work.SendTransData()
		_, err := agrt.SendWork(work, agreement)
		if err != nil {
			logger.Error("Failed to send work, err: %+v", err)
		}
	case define.Introduction:
		if agreement.Cipher != define.CIPHER {
			logger.Error("Cipher: %s, Identity: %d", agreement.Cipher, agreement.Identity)
			gos.Disconnect(define.DbaPort, work.Index)
		} else {
			s.serverIdDict[agreement.Identity] = work.Index
			logger.Info("Hello %s from %d", define.ServerName(agreement.Identity), work.Index)
		}
		work.Finish()
	default:
		logger.Warn("Unsupport system service: %d\n", agreement.Service)
		work.Finish()
	}
}

func (s *PostMessageServer) handleNormal(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	default:
		logger.Warn("Unsupport normal service: %d", agreement.Service)
		work.Finish()
	}
}

func (s *PostMessageServer) handleCommission(work *base.Work, agreement *agrt.Agreement) {
	logger.Info("Service: %d, Cid: %d", agreement.Service, agreement.Cid)

	switch agreement.Service {
	case define.AddPost:
		// ==================================================
		// 準備將請求轉送給 DBA server
		// ==================================================
		// td := base.NewTransData()
		// bs, err := agreement.Marshal()
		// if err != nil {
		// 	agreement.Msg = "Failed to marshal agreement"
		// 	logger.Error("%s, err: %+v", agreement.Msg, err)
		// 	work.Finish()
		// 	return
		// }
		// td.AddByteArray(bs)
		// data := td.FormData()

		// // 將註冊數據傳到 Dba 伺服器
		// err = gos.SendToServer(define.DbaServer, &data, int32(len(data)))

		_, err := agrt.SendToServer(define.DbaServer, agreement)

		if err != nil {
			agreement.ReturnCode = 1
			agreement.Msg = "Failed to send to Dba server"
			_, err = agrt.SendWork(work, agreement)
			if err != nil {
				logger.Error("Failed to send work, err: %+v", err)
				work.Finish()
			} else {
				logger.Info("Send define.AddPost response(%d): %+v", agreement.ReturnCode, agreement)
			}
		} else {
			logger.Info("Send define.AddPost request: %+v", agreement)
			work.Finish()
		}

	// TODO: 若該 post id 存在於緩存當中，則可直接返回，不需要再問 DBA
	case define.GetPost:
		// var bs []byte
		var err error
		pm := agreement.PostMessages[0]
		if root, ok := s.pmRoots[pm.Id]; ok {
			agreement.ReturnCode = 0
			agreement.PostMessages[0] = proto.Clone(root).(*pbgo.PostMessage)
		} else {
			// bs, err = agreement.Marshal()
			// if err != nil {
			// 	logger.Error("Failed to marshal agreement, err: %+v", err)
			// 	work.Finish()
			// 	return
			// }
			// td := base.NewTransData()
			// td.AddByteArray(bs)
			// data := td.FormData()

			// // 將註冊數據傳到 Dba 伺服器
			// err = gos.SendToServer(define.DbaServer, &data, int32(len(data)))

			_, err := agrt.SendToServer(define.DbaServer, agreement)

			if err != nil {
				agreement.ReturnCode = 2
				agreement.Msg = "Failed to query from DbaServer."
				logger.Error("%s, err: %+v", agreement.Msg, err)
			} else {
				logger.Info("Send define.GetPost request: %+v", agreement)
				work.Finish()
				return
			}
		}
		// bs, err = agreement.Marshal()
		// if err != nil {
		// 	logger.Error("Failed to marshal agreement, err: %+v", err)
		// 	work.Finish()
		// 	return
		// }
		// work.Body.AddByteArray(bs)
		// work.SendTransData()

		_, err = agrt.SendWork(work, agreement)
		if err != nil {
			logger.Error("Failed to send work, err: %+v", err)
		} else {
			logger.Info("Send define.GetPost response(%d): %+v", agreement.ReturnCode, agreement)
		}

	case define.GetMyPosts:
		userId := agreement.Accounts[0].Index

		// 取得用戶 userId 的貼文 ID 列表
		if postIds, ok := s.postIds[userId]; ok {
			agreement.ReturnCode = 0

			// 根據貼文 ID 列表，依序讀取對應的貼文
			for postId := range postIds.Elements {
				if pm, ok := s.pmRoots[postId]; ok {
					agreement.PostMessages = append(agreement.PostMessages, pm)
				}
			}
		} else {
			agreement.ReturnCode = 1
			agreement.Msg = fmt.Sprintf("Not found posts belong to user with id: %d", userId)
		}

		// bs, err := agreement.Marshal()
		// if err != nil {
		// 	logger.Error("Failed to marshal agreement, err: %+v", err)
		// 	work.Finish()
		// 	return
		// }
		// work.Body.AddByteArray(bs)
		// work.SendTransData()

		_, err := agrt.SendWork(work, agreement)
		if err != nil {
			logger.Error("Failed to send work, err: %+v", err)
		} else {
			logger.Info("Send define.GetMyPosts response(%d): %+v", agreement.ReturnCode, agreement)
		}

	case define.ModifyPost:
		// var bs []byte
		var err error
		if len(agreement.PostMessages) == 1 {
			// bs, err = agreement.Marshal()
			// if err != nil {
			// 	logger.Error("Failed to marshal agreement, err: %+v", err)
			// 	return
			// }
			// td := base.NewTransData()
			// td.AddByteArray(bs)
			// data := td.FormData()

			// // 將註冊數據傳到 Dba 伺服器
			// err = gos.SendToServer(define.DbaServer, &data, int32(len(data)))
			_, err = agrt.SendToServer(define.DbaServer, agreement)

			if err != nil {
				agreement.ReturnCode = 2
				agreement.Msg = "Failed to query to DbaServer."
				logger.Error("%s, err: %+v", agreement.Msg, err)
			} else {
				logger.Info("Send define.ModifyPost request: %+v", agreement)
				work.Finish()
				return
			}
		} else {
			agreement.ReturnCode = 1
			agreement.Msg = "Not found posts' id."
		}

		// // 若發生錯誤時，將錯誤訊息回傳 GSNS 伺服器
		// bs, err = agreement.Marshal()
		// if err != nil {
		// 	logger.Error("Failed to marshal agreement, err: %+v", err)
		// 	return
		// }
		// work.Body.AddByteArray(bs)
		// work.SendTransData()

		_, err = agrt.SendWork(work, agreement)
		if err != nil {
			logger.Error("Failed to send work, err: %+v", err)
		} else {
			logger.Info("Send define.ModifyPost response(%d): %+v", agreement.ReturnCode, agreement)
		}

	case define.GetSubscribedPosts:
		// ==================================================
		// 準備將請求轉送給 DBA server
		// ==================================================
		// td := base.NewTransData()
		// bs, err := agreement.Marshal()
		// if err != nil {
		// 	agreement.Msg = "Failed to marshal agreement"
		// 	logger.Error("%s, err: %+v", agreement.Msg, err)
		// 	work.Finish()
		// 	return
		// }
		// td.AddByteArray(bs)
		// data := td.FormData()

		// // 將註冊數據傳到 Dba 伺服器
		// err = gos.SendToServer(define.DbaServer, &data, int32(len(data)))

		_, err := agrt.SendToServer(define.DbaServer, agreement)

		if err != nil {
			agreement.ReturnCode = 1
			agreement.Msg = "Failed to send to Dba server"
			logger.Error("%s, err: %+v", agreement.Msg, err)
			// var bs []byte
			// bs, err = agreement.Marshal()
			// if err != nil {
			// 	agreement.Msg = "Failed to marshal agreement"
			// 	logger.Error("%s, err: %+v", agreement.Msg, err)
			// 	work.Finish()
			// 	return
			// }
			// work.Body.AddByteArray(bs)
			// work.SendTransData()

			_, err := agrt.SendWork(work, agreement)
			if err != nil {
				logger.Error("Failed to send work, err: %+v", err)
				work.Finish()
			} else {
				logger.Info("Send define.GetSubscribedPosts response(%d): %+v", agreement.ReturnCode, agreement)
			}
		} else {
			logger.Info("Send define.GetSubscribedPosts request: %+v", agreement)
			work.Finish()
		}

	default:
		fmt.Printf("Unsupport commission service: %d", agreement.Service)
		work.Finish()
	}
}
