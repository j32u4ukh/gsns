package pm

import (
	"fmt"
	"internal/agrt"
	"internal/define"
	"internal/pbgo"

	"github.com/j32u4ukh/gos/ans"
	"github.com/j32u4ukh/gos/base/ghttp"
)

// [endpoint]/post
// TODO: HTTP 請求處理過程中若失敗，要返回錯誤訊息給客戶端，而非印出日誌或直接返回
func (m *PostMessageMgr) HttpHandler(router *ans.Router) {
	router.POST("/", m.addNewPost)
	router.PATCH("/", m.modifyPost)
	router.GET("/<post_id int>", m.getPost)
	router.POST("/mine", m.getMyPosts)
}

// 用於新增貼文
// [endpoint]/post
func (m *PostMessageMgr) addNewPost(c *ghttp.Context) {
	pmp := &PostMessageProtocol{}
	c.ReadJson(pmp)
	m.logger.Info("PostMessageProtocol: %+v", pmp)

	if pmp.Token == "" || pmp.Content == "" {
		msg := fmt.Sprintf("缺少參數, PostMessage: %+v", pmp)
		m.logger.Error(msg)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 1,
			"msg": msg,
		})
		return
	}
	user, ok := m.getUserByTokenFunc(pmp.Token)

	if !ok {
		msg := fmt.Sprintf("Not found token %s", pmp.Token)
		m.logger.Error(msg)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"msg": msg,
		})
		// m.httpAnswer.Send(c)
		return
	}

	agreement := agrt.GetAgreement()
	defer agrt.PutAgreement(agreement)
	agreement.Cmd = define.CommissionCommand
	agreement.Service = define.AddPost
	agreement.Cid = c.GetId()
	agreement.PostMessages = append(agreement.PostMessages, &pbgo.PostMessage{
		ParentId: pmp.ParentId,
		UserId:   user.Index,
		Content:  pmp.Content,
	})

	_, err := agrt.SendToServer(define.PostMessageServer, agreement)
	if err != nil {
		msg := "Failed to send to PostMessage server"
		m.logger.Error(fmt.Sprintf("%s, err: %+v", msg, err))
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"err": msg,
		})
		// m.httpAnswer.Send(c)
	} else {
		m.logger.Info("Send define.AddPost request: %+v", agreement)
		// // 將當前 Http 的工作結束
		// m.httpAnswer.Finish(c)
	}
}

// 用於讀取特定貼文
// [endpoint]/post/<post_id int>
func (m *PostMessageMgr) getPost(c *ghttp.Context) {
	value := c.GetValue("post_id")
	if value == nil {
		msg := "Failed to get post id."
		m.logger.Error(msg)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 1,
			"msg": msg,
		})
		// m.httpAnswer.Send(c)
		return
	}

	post_id := value.(int64)
	m.logger.Info("post_id: %d", post_id)

	agreement := agrt.GetAgreement()
	defer agrt.PutAgreement(agreement)
	agreement.Cmd = define.CommissionCommand
	agreement.Service = define.GetPost
	agreement.Cid = c.GetId()
	agreement.PostMessages = append(agreement.PostMessages, &pbgo.PostMessage{
		Id: uint64(post_id),
	})

	_, err := agrt.SendToServer(define.PostMessageServer, agreement)
	if err != nil {
		m.logger.Error("err: %+v", err)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 3,
			"msg": "Failed to send data to PostMessage server.",
		})
		// m.httpAnswer.Send(c)
	} else {
		m.logger.Info("Send define.GetPost request: %+v", agreement)
		// // 將當前 Http 的工作結束
		// m.httpAnswer.Finish(c)
	}
}

// 用於讀取貼文
// [endpoint]/post/mine
func (m *PostMessageMgr) getMyPosts(c *ghttp.Context) {
	pmp := &PostMessageProtocol{}
	c.ReadJson(pmp)
	m.logger.Info("PostMessageProtocol: %+v", pmp)

	if pmp.Token == "" {
		msg := fmt.Sprintf("缺少參數, PostMessage: %+v", pmp)
		m.logger.Error(msg)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 1,
			"msg": msg,
		})
		return
	}

	user, ok := m.getUserByTokenFunc(pmp.Token)

	if !ok {
		msg := fmt.Sprintf("Not found token %s", pmp.Token)
		m.logger.Error(msg)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 2,
			"msg": msg,
		})
		// m.httpAnswer.Send(c)
		return
	}

	agreement := agrt.GetAgreement()
	defer agrt.PutAgreement(agreement)
	agreement.Cmd = define.CommissionCommand
	agreement.Service = define.GetMyPosts
	agreement.Cid = c.GetId()
	agreement.Accounts = append(agreement.Accounts, &pbgo.Account{
		Index: user.Index,
	})

	_, err := agrt.SendToServer(define.PostMessageServer, agreement)
	if err != nil {
		m.logger.Error("Failed to send to PostMessage, err: %+v", err)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"err": "Failed to send to server.",
		})
		// m.httpAnswer.Send(c)
	} else {
		m.logger.Info("Send define.GetMyPosts request: %+v", agreement)
		// // 將當前 Http 的工作結束
		// m.httpAnswer.Finish(c)
	}
}

// 用於編輯貼文
// [endpoint]/post
func (m *PostMessageMgr) modifyPost(c *ghttp.Context) {
	pmp := &PostMessageProtocol{}
	c.ReadJson(pmp)

	if pmp.Token == "" || pmp.PostId == 0 || pmp.Content == "" {
		msg := fmt.Sprintf("缺少參數, PostMessage: %+v", pmp)
		m.logger.Error(msg)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 1,
			"msg": msg,
		})
		return
	}

	user, ok := m.getUserByTokenFunc(pmp.Token)

	if !ok {
		msg := fmt.Sprintf("Not found token %s", pmp.Token)
		m.logger.Error(msg)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"msg": msg,
		})
		// m.httpAnswer.Send(c)
		return
	}

	agreement := agrt.GetAgreement()
	defer agrt.PutAgreement(agreement)
	agreement.Cmd = define.CommissionCommand
	agreement.Service = define.ModifyPost
	agreement.Cid = c.GetId()
	agreement.PostMessages = append(agreement.PostMessages, &pbgo.PostMessage{
		Id:       pmp.PostId,
		ParentId: pmp.ParentId,
		UserId:   user.Index,
		Content:  pmp.Content,
	})

	_, err := agrt.SendToServer(define.PostMessageServer, agreement)
	if err != nil {
		msg := "Failed to send to PostMessage server"
		m.logger.Error("%s, err: %+v", msg, err)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"msg": msg,
		})
		// m.httpAnswer.Send(c)
	} else {
		m.logger.Info("Send define.ModifyPost request: %+v", agreement)
		// // 將當前 Http 的工作結束
		// m.httpAnswer.Finish(c)
	}
}
