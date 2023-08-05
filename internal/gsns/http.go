package gsns

import (
	"fmt"
	"internal/agrt"
	"internal/define"
	"internal/pbgo"
	"strconv"

	"github.com/j32u4ukh/gos"
	"github.com/j32u4ukh/gos/ans"
	"github.com/j32u4ukh/gos/base"
	"github.com/j32u4ukh/gos/base/ghttp"
)

// TODO: HTTP 請求處理過程中若失敗，要返回錯誤訊息給客戶端，而非印出日誌或直接返回
// [endpoint]/social
func (s *MainServer) HttpSocialHandler(router *ans.Router) {
	// 取得其他用戶的清單
	router.GET("/other_users", s.getOtherUsers)
	router.POST("/subscribe", s.subscribe)
}

// [endpoint]/social/other_users
func (s *MainServer) getOtherUsers(c *ghttp.Context) {
	var sToken string
	var ok bool

	if sToken, ok = c.Params["token"]; !ok {
		msg := "Not found parameter token"
		logger.Error(msg)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 1,
			"msg": msg,
		})
		s.Http.Send(c)
		return
	}

	token, err := strconv.ParseUint(sToken, 10, 64)

	if err != nil {
		msg := "Invalid token"
		logger.Error(msg)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 2,
			"msg": msg,
		})
		s.Http.Send(c)
		return
	}

	user, ok := s.AMgr.GetUserByToken(token)

	if !ok {
		msg := fmt.Sprintf("Not found user with token(%d)", token)
		logger.Error(msg)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 3,
			"msg": msg,
		})
		s.Http.Send(c)
		return
	}

	agreement := agrt.GetAgreement()
	defer agrt.PutAgreement(agreement)
	agreement.Cmd = define.CommissionCommand
	agreement.Service = define.GetOtherUsers
	agreement.Cid = c.GetId()
	agreement.Accounts = append(agreement.Accounts, &pbgo.Account{
		Index: user.Index,
	})
	bs, err := agreement.Marshal()

	if err != nil {
		msg := "Failed to marshal agreement."
		logger.Error(fmt.Sprintf("%s, err: %+v", msg, err))
		c.Json(ghttp.StatusInternalServerError, ghttp.H{
			"ret": 4,
			"msg": msg,
		})
		s.Http.Send(c)
		return
	}

	// 寫入 agreement
	td := base.NewTransData()
	td.AddByteArray(bs)
	data := td.FormData()

	// 將註冊數據傳到 Account 伺服器
	err = gos.SendToServer(define.AccountServer, &data, int32(len(data)))

	if err != nil {
		msg := "Failed to send request to account server"
		logger.Error(fmt.Sprintf("%s, err: %+v", msg, err))
		c.Json(ghttp.StatusInternalServerError, ghttp.H{
			"ret": 5,
			"msg": msg,
		})
		s.Http.Send(c)
	} else {
		logger.Info("Send define.GetOtherUsers request: %+v", agreement)
		s.Http.Finish(c)
	}
}

// [endpoint]/social/subscribe
func (s *MainServer) subscribe(c *ghttp.Context) {
	ip := &SocialProtocol{}
	c.ReadJson(ip)

	if ip.Token == 0 || ip.TargetId == 0 {
		msg := "缺少參數"
		logger.Error(msg)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 1,
			"msg": msg,
		})
		s.Http.Send(c)
		return
	}

	user, ok := s.AMgr.GetUserByToken(ip.Token)

	if !ok {
		msg := fmt.Sprintf("Not found token %d", ip.Token)
		logger.Error(msg)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 2,
			"msg": msg,
		})
		s.Http.Send(c)
		return
	}

	// 避免重複訂閱，先檢查訂閱對象的 ID 再送出請求
	if edges, ok := s.AMgr.Edges[user.Index]; ok {
		if edges.Contains(ip.TargetId) {
			msg := fmt.Sprintf("User %s has subscribed user %d", user.Name, ip.TargetId)
			logger.Info(msg)
			c.Json(ghttp.StatusOK, ghttp.H{
				"ret": 0,
				"msg": msg,
			})
			s.Http.Send(c)
			return
		}
	}

	agreement := agrt.GetAgreement()
	defer agrt.PutAgreement(agreement)
	agreement.Cmd = define.CommissionCommand
	agreement.Service = define.Subscribe
	agreement.Cid = c.GetId()
	agreement.Edges = append(agreement.Edges, &pbgo.Edge{
		UserId: user.Index,
		Target: ip.TargetId,
	})

	// 寫入 agreement
	bs, err := agreement.Marshal()

	if err != nil {
		msg := "Failed to marshal agreement."
		logger.Error(fmt.Sprintf("%s, err: %+v", msg, err))
		c.Json(ghttp.StatusInternalServerError, ghttp.H{
			"ret": 3,
			"msg": msg,
		})
		s.Http.Send(c)
		return
	}

	td := base.NewTransData()
	td.AddByteArray(bs)
	data := td.FormData()

	// 將註冊數據傳到 Account 伺服器
	err = gos.SendToServer(define.AccountServer, &data, int32(len(data)))

	if err != nil {
		msg := "Failed to send request to account server"
		logger.Error("%s, err: %+v", msg, err)
		c.Json(ghttp.StatusInternalServerError, ghttp.H{
			"ret": 4,
			"msg": msg,
		})
		s.Http.Send(c)
		return
	} else {
		logger.Info("Send define.Subscribe request: %+v", agreement)
		s.Http.Finish(c)
	}
}
