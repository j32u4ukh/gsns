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

type PostMessageServer struct {
	Tcp *ans.Tcp0Anser
	// 獨立的貼文、不是回覆他人的貼文
	pmRoots map[uint64]*pbgo.PostMessage
	// 回覆他人的貼文，parent id 為被回覆的貼文的 post id
	// key1: post id, key2: parent id
	pmLeaves *cntr.BikeyMap[uint64, uint64, *pbgo.PostMessage]

	// key: server id, value: conn id
	serverIdDict  map[int32]int32
	heartbeatTime time.Time
}

func NewPostMessageServer() *PostMessageServer {
	s := &PostMessageServer{
		pmRoots:       make(map[uint64]*pbgo.PostMessage),
		pmLeaves:      cntr.NewBikeyMap[uint64, uint64, *pbgo.PostMessage](),
		serverIdDict:  make(map[int32]int32),
		heartbeatTime: time.Now(),
	}
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
		bs, _ := agreement.Marshal()
		work.Body.AddByteArray(bs)
		work.SendTransData()
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
		work.Finish()

		// ==================================================
		// 準備將請求轉送給 DBA server
		// ==================================================
		td := base.NewTransData()
		bs, _ := agreement.Marshal()
		td.AddByteArray(bs)
		data := td.FormData()
		logger.Info("data: %+v", data)

		// 將註冊數據傳到 Dba 伺服器
		err := gos.SendToServer(define.DbaServer, &data, td.GetLength())

		if err != nil {
			logger.Error("Failed to send to server %d: %v\nError: %+v", define.DbaServer, data, err)
			return
		}

	case define.GetPost:
		work.Body.Clear()
		if len(agreement.PostMessages) == 0 {
			agreement.ReturnCode = 1
			agreement.Msg = "Not found posts' id."
		} else {
			for i, pm := range agreement.PostMessages {
				if root, ok := s.pmRoots[pm.Id]; ok {
					agreement.PostMessages[i] = proto.Clone(root).(*pbgo.PostMessage)
				}
			}
			agreement.ReturnCode = 0
		}
		bs, _ := agreement.Marshal()
		work.Body.AddByteArray(bs)
		work.SendTransData()

	default:
		fmt.Printf("Unsupport commission service: %d", agreement.Service)
		work.Finish()
	}
}
