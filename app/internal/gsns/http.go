package gsns

import (
	"fmt"
	"internal/agrt"
	"internal/define"
	"internal/pbgo"
	"internal/utils"
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

	if sUserId, ok = c.Params["user_id"]; !ok {
		clientLogger.Error(utils.JsonResponse(c, define.Error.MissingParameters, "user_id"))
		return
	}

	userId, err := strconv.ParseInt(sUserId, 10, 64)
	if err != nil {
		clientLogger.Error(utils.JsonResponse(c, define.Error.NotFoundUser, "user_id", sUserId))
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

	s.handleRequest(c, define.AccountServer, agreement)
	// _, err = agrt.SendToServer(define.AccountServer, agreement)
	// if err != nil {
	// 	msg := utils.JsonResponse(c, define.Error.CannotSendMessage, "to Account server")
	// 	serverLogger.Error(fmt.Sprintf("%s, err: %+v", msg, err))
	// } else {
	// 	serverLogger.Info("Send define.GetOtherUsers request: %+v", agreement)
	// }
}

// [endpoint]/social/subscribe
func (s *MainServer) subscribe(c *ghttp.Context) {
	ip := &SocialProtocol{}
	err := c.ReadJson(ip)

	if err != nil {
		msg := utils.JsonResponse(c, define.Error.InvalidBodyData)
		clientLogger.Error("%s, err: %+v", msg, err)
		return
	}

	if ip.Token == "" || ip.TargetId == 0 {
		clientLogger.Error(utils.JsonResponse(c, define.Error.MissingParameters, "token or target_id"))
		return
	}

	user, ok := s.AMgr.GetUserByToken(ip.Token)
	if !ok {
		clientLogger.Error(utils.JsonResponse(c, define.Error.NotFoundUser, "token", ip.Token))
		return
	}

	if user.Index == ip.TargetId {
		clientLogger.Error(utils.JsonResponse(c, define.Error.InvalidTarget, fmt.Sprintf("User(%d)", user.Index), fmt.Sprintf("User(%d)", ip.TargetId)))
		return
	}

	// 避免重複訂閱，先檢查訂閱對象的 ID 再送出請求
	if edges, ok := s.AMgr.Edges[user.Index]; ok {
		if edges.Contains(ip.TargetId) {
			clientLogger.Info(utils.JsonResponse(c, define.Error.DuplicateEntity, fmt.Sprintf("User %s has subscribed user %d", user.Name, ip.TargetId)))
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

	s.handleRequest(c, define.AccountServer, agreement)
	// _, err = agrt.SendToServer(define.AccountServer, agreement)
	// if err != nil {
	// 	msg := utils.JsonResponse(c, define.Error.CannotSendMessage, "to Dba server")
	// 	serverLogger.Error("%s, err: %+v", msg, err)
	// 	return
	// } else {
	// 	serverLogger.Info("Send define.Subscribe request: %+v", agreement)
	// }
}

// [endpoint]/social/subscribed_posts
func (s *MainServer) getSubscribedPosts(c *ghttp.Context) {
	ip := &SocialProtocol{}
	err := c.ReadJson(ip)

	if err != nil {
		msg := utils.JsonResponse(c, define.Error.InvalidBodyData)
		clientLogger.Error("%s, err: %+v", msg, err)
		return
	}

	if ip.Token == "" {
		clientLogger.Error(utils.JsonResponse(c, define.Error.MissingParameters, "token"))
		return
	}

	user, ok := s.AMgr.GetUserByToken(ip.Token)
	if !ok {
		clientLogger.Error(utils.JsonResponse(c, define.Error.NotFoundUser, "token", ip.Token))
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

	s.handleRequest(c, define.PostMessageServer, agreement)
	// _, err = agrt.SendToServer(define.PostMessageServer, agreement)
	// if err != nil {
	// 	msg := utils.JsonResponse(c, define.Error.CannotSendMessage, "to PostMessage server")
	// 	serverLogger.Error("%s, err: %+v", msg, err)
	// } else {
	// 	serverLogger.Info("Send define.GetSubscribedPosts request: %+v", agreement)
	// }
}

func (s *MainServer) handleRequest(c *ghttp.Context, serverId int32, agreement *agrt.Agreement) {
	_, err := agrt.SendToServer(serverId, agreement)
	if err != nil {
		msg := utils.JsonResponse(c, define.Error.CannotSendMessage, fmt.Sprintf("to %s server", define.ServerName(serverId)))
		serverLogger.Error("%s, err: %+v", msg, err)
	} else {
		serverLogger.Info("Send %s request: %+v", define.ServiceName(agreement.Service), agreement)
	}
}
