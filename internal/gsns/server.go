package gsns

import (
	"internal/gsns/account"
	pm "internal/gsns/post_message"

	"github.com/j32u4ukh/gos/ans"
)

type MainServer struct {
	HttpAnswer *ans.HttpAnser
	AMgr       *account.AccountMgr
	PMgr       *pm.PostMessageMgr
}

func newMainServer() *MainServer {
	m := &MainServer{
		AMgr: account.NewAccountMgr(logger),
		PMgr: pm.NewPostMessageMgr(logger),
	}
	m.PMgr.SetFuncGetUserByToken(m.AMgr.GetUserByToken)
	return m
}

func (s *MainServer) SetHttpAnswer(a *ans.HttpAnser) {
	s.HttpAnswer = a

	// 帳號相關節點
	s.AMgr.SetHttpAnswer(a)
	rAccount := a.Router.NewRouter("/account")
	s.AMgr.HttpAccountHandler(rAccount)
	rSocial := a.Router.NewRouter("/social")
	s.AMgr.HttpSocialHandler(rSocial)

	// 貼文相關節點
	s.PMgr.SetHttpAnswer(a)
	rPost := a.Router.NewRouter("/post")
	s.PMgr.HttpHandler(rPost)
}

func (s *MainServer) Run() {

}
