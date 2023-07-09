package gsns

import (
	"internal/gsns/account"

	"github.com/j32u4ukh/gos/ans"
)

type MainServer struct {
	HttpAnswer *ans.HttpAnser
	AMgr       *account.AccountMgr
}

func newMainServer() *MainServer {
	m := &MainServer{
		AMgr: account.NewAccountMgr(logger),
	}
	return m
}

func (s *MainServer) SetHttpAnswer(a *ans.HttpAnser) {
	s.HttpAnswer = a
	s.AMgr.SetHttpAnswer(a)
}

func (s *MainServer) HttpHandler(router *ans.Router) {
	// 帳號相關節點
	rAccount := router.NewRouter("/account")
	s.AMgr.HttpHandler(rAccount)
}

func (s *MainServer) Run() {

}
