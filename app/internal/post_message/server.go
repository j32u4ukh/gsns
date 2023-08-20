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
		serverLogger.Error("Failed to unmarshal agreement, err: %+v", err)
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
		clientLogger.Warn("Unsupport command: %d", agreement.Cmd)
		work.Finish()
	}
}

func (s *PostMessageServer) handleSystem(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	// 回應心跳包
	case define.Heartbeat:
		agreement.ReturnCode = define.Error.None
		agreement.Msg = "OK"
		_, err := agrt.SendWork(work, agreement)
		if err != nil {
			_, _, agreement.Msg = define.ErrorMessage(define.Error.CannotSendMessage, "work")
			serverLogger.Error("%s, err: %+v", agreement.Msg, err)
		}
	case define.Introduction:
		if agreement.Cipher != define.CIPHER {
			_, _, agreement.Msg = define.ErrorMessage(define.Error.WrongConnectionIdentity, agreement.Cipher, agreement.Identity)
			clientLogger.Error(agreement.Msg)
			gos.Disconnect(define.DbaPort, work.Index)
		} else {
			s.serverIdDict[agreement.Identity] = work.Index
			serverLogger.Info("Hello %s from %d", define.ServerName(agreement.Identity), work.Index)
		}
		work.Finish()
	default:
		clientLogger.Warn("Unsupport system service: %d\n", agreement.Service)
		work.Finish()
	}
}

func (s *PostMessageServer) handleNormal(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	default:
		clientLogger.Warn("Unsupport normal service: %d", agreement.Service)
		work.Finish()
	}
}

func (s *PostMessageServer) handleCommission(work *base.Work, agreement *agrt.Agreement) {
	serverLogger.Info("Service: %d, Cid: %d", agreement.Service, agreement.Cid)

	switch agreement.Service {
	case define.AddPost:
		s.handleCommissionRequest(work, agreement)

	// TODO: 若該 post id 存在於緩存當中，則可直接返回，不需要再問 DBA
	case define.GetPost:
		var err error
		pm := agreement.PostMessages[0]
		if root, ok := s.pmRoots[pm.Id]; ok {
			agreement.ReturnCode = define.Error.None
			agreement.PostMessages[0] = proto.Clone(root).(*pbgo.PostMessage)
		} else {
			_, err = agrt.SendToServer(define.DbaServer, agreement)
			if err != nil {
				_, agreement.ReturnCode, agreement.Msg = define.ErrorMessage(define.Error.CannotSendMessage, "to Dba server")
				serverLogger.Error("%s, err: %+v", agreement.Msg, err)
			} else {
				serverLogger.Info("Send define.GetPost request: %+v", agreement)
				work.Finish()
				return
			}
		}
		s.handleRequest(work, agreement)

	case define.GetMyPosts:
		userId := agreement.Accounts[0].Index

		// 取得用戶 userId 的貼文 ID 列表
		if postIds, ok := s.postIds[userId]; ok {
			agreement.ReturnCode = define.Error.None

			// 根據貼文 ID 列表，依序讀取對應的貼文
			for postId := range postIds.Elements {
				if pm, ok := s.pmRoots[postId]; ok {
					agreement.PostMessages = append(agreement.PostMessages, pm)
				}
			}
		} else {
			_, agreement.ReturnCode, agreement.Msg = define.ErrorMessage(define.Error.NotFoundUser, "user_id", userId)
		}

		s.handleRequest(work, agreement)

	case define.ModifyPost:
		nPost := len(agreement.PostMessages)
		if nPost == 1 {
			s.handleCommissionRequest(work, agreement)
		} else {
			_, agreement.ReturnCode, agreement.Msg = define.ErrorMessage(define.Error.TooManyParameter, 1, nPost)
			s.handleRequest(work, agreement)
		}

	case define.GetSubscribedPosts:
		s.handleCommissionRequest(work, agreement)

	default:
		fmt.Printf("Unsupport commission service: %d", agreement.Service)
		work.Finish()
	}
}

func (s *PostMessageServer) handleRequest(work *base.Work, agreement *agrt.Agreement) {
	_, err := agrt.SendWork(work, agreement)
	if err != nil {
		_, _, agreement.Msg = define.ErrorMessage(define.Error.CannotSendMessage, "work")
		serverLogger.Error("%s, err: %+v", agreement.Msg, err)
	} else {
		serverLogger.Info("Send %s request: %+v", define.ServiceName(agreement.Service), agreement)
	}
}

func (s *PostMessageServer) handleCommissionRequest(work *base.Work, agreement *agrt.Agreement) {
	_, err := agrt.SendToServer(define.DbaServer, agreement)
	if err != nil {
		_, agreement.ReturnCode, agreement.Msg = define.ErrorMessage(define.Error.CannotSendMessage, "to Dba server")
		serverLogger.Error("%s, err: %+v", agreement.Msg, err)
		_, err = agrt.SendWork(work, agreement)
		if err != nil {
			_, _, agreement.Msg = define.ErrorMessage(define.Error.CannotSendMessage, "work")
			serverLogger.Error("%s, err: %+v", agreement.Msg, err)
			work.Finish()
		} else {
			serverLogger.Info("Send %s response(%d): %+v", define.ServiceName(agreement.Service), agreement.ReturnCode, agreement)
		}
	} else {
		serverLogger.Info("Send %s request: %+v", define.ServiceName(agreement.Service), agreement)
		work.Finish()
	}
}
