package account

import (
	"fmt"
	"internal/agrt"
	"internal/define"
	"internal/pbgo"
	"internal/utils"

	"github.com/j32u4ukh/gos/ans"
	"github.com/j32u4ukh/gos/base/ghttp"
)

// TODO: HTTP 請求處理過程中若失敗，要返回錯誤訊息給客戶端，而非印出日誌或直接返回
// [endpoint]/account
func (m *AccountMgr) HttpAccountHandler(router *ans.Router) {

	// TODO: 註冊：帳密以及個人資訊。儲存以 SHA256 加密後的密碼，而非儲存原始密碼。
	router.POST("/register", m.register)

	// TODO: token 應該要有時效
	router.POST("/login", m.login)

	// TODO:
	router.POST("/logout", m.logout)

	// 取得用戶資訊
	router.POST("/get_user_info", m.getUserInfo)
	router.POST("/set_user_info", m.setUserInfo)
}

func (m *AccountMgr) register(c *ghttp.Context) {
	ap := &AccountProtocol{}
	err := c.ReadJson(ap)

	if err != nil {
		msg := utils.JsonResponse(c, define.Error.InvalidBodyData)
		m.clientLogger.Error("%s, err: %+v", msg, err)
		return
	}

	m.serverLogger.Info("AccountProtocol: %+v", ap)

	// 帳號名稱(Account) 和 密碼原文(Password) 為必須，個人資訊(Info) 可以不填
	// TODO: 在前端就加密
	if ap.Account == "" || ap.Password == "" {
		m.clientLogger.Error(utils.JsonResponse(c, define.Error.MissingParameters, "account or password"))
		return
	}

	agreement := agrt.GetAgreement()
	defer agrt.PutAgreement(agreement)
	agreement.Cmd = define.CommissionCommand
	agreement.Service = define.Register
	agreement.Cid = c.GetId()
	agreement.Accounts = append(agreement.Accounts, &pbgo.Account{
		Account:  ap.Account,
		Password: ap.Password,
		Info:     ap.Info,
	})

	_, err = agrt.SendToServer(define.AccountServer, agreement)
	if err != nil {
		m.serverLogger.Error("%s, err: %+v", utils.JsonResponse(c, define.Error.CannotSendMessage, "to Account server"), err)
	} else {
		m.serverLogger.Info("Send define.Register request: %+v", agreement)
	}
}

func (m *AccountMgr) login(c *ghttp.Context) {
	ap := &AccountProtocol{}
	err := c.ReadJson(ap)
	if err != nil {
		msg := utils.JsonResponse(c, define.Error.InvalidBodyData)
		m.clientLogger.Error("%s, err: %+v", msg, err)
		return
	}
	m.serverLogger.Info("AccountProtocol: %+v", ap)

	if ap.Account == "" || ap.Password == "" {
		m.clientLogger.Error(utils.JsonResponse(c, define.Error.MissingParameters, "account or password"))
		return
	}

	agreement := agrt.GetAgreement()
	defer agrt.PutAgreement(agreement)
	agreement.Cmd = int32(define.CommissionCommand)
	agreement.Service = int32(define.Login)
	agreement.Cid = c.GetId()
	agreement.Accounts = append(agreement.Accounts, &pbgo.Account{
		Account:  ap.Account,
		Password: ap.Password,
	})

	_, err = agrt.SendToServer(define.AccountServer, agreement)
	if err != nil {
		msg := utils.JsonResponse(c, define.Error.CannotSendMessage, "to Account server")
		m.serverLogger.Error("%s, err: %+v", msg, err)
	} else {
		m.serverLogger.Info("Send define.Login request: %+v", agreement)
	}
}

func (m *AccountMgr) logout(c *ghttp.Context) {
	ap := &AccountProtocol{}
	err := c.ReadJson(ap)
	if err != nil {
		msg := utils.JsonResponse(c, define.Error.InvalidBodyData)
		m.serverLogger.Error("%s, err: %+v", msg, err)
		return
	}

	if ap.Token == "" {
		m.clientLogger.Error(utils.JsonResponse(c, define.Error.MissingParameters, "token"))
		return
	}

	user, ok := m.users.GetByKey2(ap.Token)
	if !ok {
		m.clientLogger.Error(utils.JsonResponse(c, define.Error.NotFoundUser, "token", ap.Token))
	} else {
		m.users.DelByKey2(ap.Token)
		c.Json(200, ghttp.H{
			"msg": fmt.Sprintf("User %s logout success.", user.Name),
		})
	}
}

func (m *AccountMgr) getUserInfo(c *ghttp.Context) {
	ap := &AccountProtocol{}
	err := c.ReadJson(ap)
	if err != nil {
		msg := utils.JsonResponse(c, define.Error.InvalidBodyData)
		m.clientLogger.Error("%s, err: %+v", msg, err)
		return
	}

	if ap.Token == "" {
		m.clientLogger.Error(utils.JsonResponse(c, define.Error.MissingParameters, "token"))
		return
	}

	user, ok := m.users.GetByKey2(ap.Token)
	if !ok {
		m.clientLogger.Error(utils.JsonResponse(c, define.Error.NotFoundUser, "token", ap.Token))
	} else {
		c.Json(ghttp.StatusOK, ghttp.H{
			"name": user.Name,
			"info": user.Info,
		})
	}
}

func (m *AccountMgr) setUserInfo(c *ghttp.Context) {
	ap := &AccountProtocol{}
	err := c.ReadJson(ap)
	if err != nil {
		msg := utils.JsonResponse(c, define.Error.InvalidBodyData)
		m.clientLogger.Error("%s, err: %+v", msg, err)
		return
	}

	if ap.Token == "" {
		m.clientLogger.Error(utils.JsonResponse(c, define.Error.MissingParameters, "token"))
		return
	}

	user, ok := m.users.GetByKey2(ap.Token)

	if !ok {
		m.clientLogger.Error(utils.JsonResponse(c, define.Error.NotFoundUser, "token", ap.Token))
		return
	}

	agreement := agrt.GetAgreement()
	defer agrt.PutAgreement(agreement)
	agreement.Cmd = int32(define.CommissionCommand)
	agreement.Service = int32(define.SetUserData)
	agreement.Cid = c.GetId()

	// 有數值，才更新
	if ap.Account != "" {
		user.Name = ap.Account
	}

	// 有數值，才更新
	if ap.Info != "" {
		user.Info = ap.Info
	}

	agreement.Accounts = append(agreement.Accounts, &pbgo.Account{
		Index:   user.Index,
		Account: user.Name,
		Info:    user.Info,
	})

	_, err = agrt.SendToServer(define.AccountServer, agreement)
	if err != nil {
		msg := utils.JsonResponse(c, define.Error.CannotSendMessage, "to Account server")
		m.serverLogger.Error("%s, err: %+v", msg, err)
	} else {
		m.serverLogger.Info("Send define.SetUserData request: %+v", agreement)
	}
}
