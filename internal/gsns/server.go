package gsns

import (
	"fmt"
	"internal/gsns/user"
	"internal/pbgo"

	"github.com/j32u4ukh/gos"
	"github.com/j32u4ukh/gos/ans"
	"github.com/j32u4ukh/gos/base"
	"github.com/j32u4ukh/gos/base/ghttp"
)

type MainServer struct {
	HttpAnswer *ans.HttpAnser
	userMgr    *user.UserMgr
}

func newMainServer(nUser int32) *MainServer {
	m := &MainServer{
		userMgr: user.NewUserMgr(nUser, logger),
	}
	return m
}

func (s *MainServer) AddUser(user *pbgo.SnsUser) {
	s.userMgr.AddUser(user)
}

func (s *MainServer) HttpHandler(router *ans.Router) {
	// 帳號相關節點
	rAccount := router.NewRouter("/account")

	// TODO: 註冊：帳密以及個人資訊。儲存以 SHA256 加密後的密碼，而非儲存原始密碼。
	rAccount.POST("/register", s.register)

	// TODO:
	rAccount.POST("/login", s.logout)

	// TODO:
	rAccount.POST("/logout", s.logout)

	// 轉交工作範例
	rAccount.POST("/delay_response", func(c *ghttp.Context) {
		s.HttpAnswer.Finish(c)
		s.CommissionHandler(1023, c.GetId())
	})
}

func (s *MainServer) DbaHandler(work *base.Work) {
	cmd := work.Body.PopByte()

	switch cmd {
	case 0:
		s.handleSystemCommand(work)
	case 1:
		s.handleCommission(work)
	default:
		fmt.Printf("Unsupport command: %d\n", cmd)
		work.Finish()
	}
}

func (s *MainServer) handleSystemCommand(work *base.Work) {
	service := work.Body.PopUInt16()

	switch service {
	case 0:
		response := work.Body.PopString()
		fmt.Printf("Heart beat response: %s\n", response)
		work.Finish()
	default:
		fmt.Printf("Unsupport service: %d\n", service)
		work.Finish()
	}
}

func (s *MainServer) handleCommission(work *base.Work) {
	commission := work.Body.PopUInt16()

	switch commission {
	case 1023:
		c := s.HttpAnswer.GetContext(-1)
		c.Cid = work.Body.PopInt32()
		response := work.Body.PopString()
		fmt.Printf("(m *Mgr) handleCommission | response: %s\n", response)
		work.Finish()

		c.Json(200, ghttp.H{
			"index": 5,
			"msg":   fmt.Sprintf("POST | /abc/delay_response: %s", response),
		})
		s.HttpAnswer.Send(c)

	default:
		fmt.Printf("Unsupport commission: %d\n", commission)
		work.Finish()
	}
}

func (s *MainServer) CommissionHandler(site int32, cid int32) {
	if site == 1023 {
		td := base.NewTransData()
		td.AddByte(1)
		td.AddUInt16(1023)
		td.AddInt32(cid)
		data := td.FormData()
		err := gos.SendToServer(EDbaServer, &data, td.GetLength())

		if err != nil {
			fmt.Printf("(s *MainServer) CommissionHandler | Failed to send to server %d: %v\nError: %+v\n", EDbaServer, data, err)
			return
		}
	}
}

func (s *MainServer) AccountHandler(work *base.Work) {
	cmd := work.Body.PopByte()

	switch cmd {
	case 0:
		s.handleSystemCommand(work)
	case 1:
		s.handleCommission(work)
	default:
		fmt.Printf("Unsupport command: %d\n", cmd)
		work.Finish()
	}
}

func (s *MainServer) Run() {

}
