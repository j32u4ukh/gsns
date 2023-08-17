package gsns

import (
	"fmt"
	"internal/agrt"
	"internal/define"
	"internal/pbgo"
	"strconv"
	"time"

	"github.com/j32u4ukh/cntr"
	"github.com/j32u4ukh/gos/ans"
	"github.com/j32u4ukh/gos/base/ghttp"
)

// [endpoint]/social
func (s *MainServer) HttpSocialHandler(router *ans.Router) {
	// 取得其他用戶的清單
	router.GET("/other_users", s.getOtherUsers)
	router.POST("/subscribe", s.subscribe)
	router.POST("/subscribed_posts", s.getSubscribedPosts)
}

// [endpoint]/social/other_users
func (s *MainServer) getOtherUsers(c *ghttp.Context) {
	var sUserId string
	var ok bool

	// TODO: MissingParameters
	if sUserId, ok = c.Params["user_id"]; !ok {
		msg := "Not found parameter user_id"
		logger.Error(msg)
		c.Json(ghttp.StatusNotFound, ghttp.H{
			"ret": 1,
			"msg": msg,
		})
		return
	}

	// TODO: NotFoundUser
	userId, err := strconv.ParseInt(sUserId, 10, 64)

	if err != nil {
		msg := "Invalid user id."
		logger.Error(msg)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 2,
			"msg": msg,
		})
		return
	}

	agreement := agrt.GetAgreement()
	defer agrt.PutAgreement(agreement)
	agreement.Cmd = define.CommissionCommand
	agreement.Service = define.GetOtherUsers
	agreement.Cid = c.GetId()
	agreement.Accounts = append(agreement.Accounts, &pbgo.Account{
		Index: int32(userId),
	})

	_, err = agrt.SendToServer(define.AccountServer, agreement)
	if err != nil {
		// TODO: CannotSendToServer
		msg := "Failed to send request to account server"
		logger.Error(fmt.Sprintf("%s, err: %+v", msg, err))
		c.Json(ghttp.StatusInternalServerError, ghttp.H{
			"ret": 5,
			"msg": msg,
		})
	} else {
		logger.Info("Send define.GetOtherUsers request: %+v", agreement)
	}
}

// [endpoint]/social/subscribe
func (s *MainServer) subscribe(c *ghttp.Context) {
	ip := &SocialProtocol{}
	c.ReadJson(ip)

	// TODO: MissingParameters
	if ip.Token == "" || ip.TargetId == 0 {
		msg := "缺少參數"
		logger.Error(msg)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 1,
			"msg": msg,
		})
		return
	}

	user, ok := s.AMgr.GetUserByToken(ip.Token)

	// TODO: NotFoundUser
	if !ok {
		msg := fmt.Sprintf("Not found token %s", ip.Token)
		logger.Error(msg)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 2,
			"msg": msg,
		})
		return
	}

	// TODO: InvalidTarget
	if user.Index == ip.TargetId {
		msg := fmt.Sprintf("不能訂閱自己 User(%d), Target(%d)", user.Index, ip.TargetId)
		logger.Error(msg)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 3,
			"msg": msg,
		})
		return
	}

	// 避免重複訂閱，先檢查訂閱對象的 ID 再送出請求
	if edges, ok := s.AMgr.Edges[user.Index]; ok {
		if edges.Contains(ip.TargetId) {
			msg := fmt.Sprintf("User %s has subscribed user %d", user.Name, ip.TargetId)
			logger.Info(msg)
			// TODO: DuplicateEntity
			c.Json(ghttp.StatusConflict, ghttp.H{
				"ret": 0,
				"msg": msg,
			})
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

	// TODO: CannotSendMessage
	_, err := agrt.SendToServer(define.AccountServer, agreement)
	if err != nil {
		msg := "Failed to send request to account server"
		logger.Error("%s, err: %+v", msg, err)
		c.Json(ghttp.StatusInternalServerError, ghttp.H{
			"ret": 5,
			"msg": msg,
		})
		return
	} else {
		logger.Info("Send define.Subscribe request: %+v", agreement)
	}
}

// [endpoint]/social/subscribed_posts
func (s *MainServer) getSubscribedPosts(c *ghttp.Context) {
	ip := &SocialProtocol{}
	c.ReadJson(ip)

	// TODO: MissingParameters
	if ip.Token == "" {
		msg := "缺少參數"
		logger.Error(msg)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 1,
			"msg": msg,
		})
		return
	}

	var user *pbgo.User
	var ok bool
	user, ok = s.AMgr.GetUserByToken(ip.Token)

	// TODO: NotFoundUser
	if !ok {
		msg := fmt.Sprintf("Not found token %s", ip.Token)
		logger.Error(msg)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 2,
			"msg": msg,
		})
		return
	}

	agreement := agrt.GetAgreement()
	defer agrt.PutAgreement(agreement)

	// 取得訂閱對象的 ID
	var edges *cntr.Set[int32]
	returnImmediate := false

	if edges, ok = s.AMgr.Edges[user.Index]; !ok {
		returnImmediate = true
	} else if edges.Length() == 0 {
		returnImmediate = true
	}

	if returnImmediate {
		c.Json(ghttp.StatusOK, ghttp.H{
			"n_post": len(agreement.PostMessages),
			"posts":  agreement.PostMessages,
		})
		return
	}

	agreement.Cmd = define.CommissionCommand
	agreement.Service = define.GetSubscribedPosts
	agreement.Cid = c.GetId()
	var err error

	if ip.StartUtc != 0 {
		agreement.StartUtc = ip.StartUtc
	} else {
		agreement.StartUtc = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC).Unix()
	}

	if ip.StopUtc != 0 {
		agreement.StopUtc = ip.StopUtc
	} else {
		agreement.StopUtc = time.Now().UTC().Unix()
	}

	for edge := range edges.Elements {
		agreement.Accounts = append(agreement.Accounts, &pbgo.Account{
			Index: edge,
		})
	}

	// TODO: CannotSendMessage
	_, err = agrt.SendToServer(define.PostMessageServer, agreement)
	if err != nil {
		msg := "Failed to sned to PostMessage server."
		logger.Error("%s, err: %+v", msg, err)
		c.Json(ghttp.StatusInternalServerError, ghttp.H{
			"ret": 4,
			"msg": msg,
		})
	} else {
		logger.Info("Send define.GetSubscribedPosts request: %+v", agreement)
	}
}
