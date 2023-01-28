package gsns

import (
	"fmt"
	"time"

	"github.com/j32u4ukh/gos"
	"github.com/j32u4ukh/gos/ans"
	"github.com/j32u4ukh/gos/ask"
	"github.com/j32u4ukh/gos/base"
	"github.com/j32u4ukh/gos/base/ghttp"
	"github.com/j32u4ukh/gos/define"
)

const EDbaServer int32 = 0

var s *MainServer
var dbaAsker *ask.Tcp0Asker

func Init() {
	// ==================================================
	// Http Server: 接受來自客戶端的請求
	// ==================================================
	var port int = 1023
	anser, err := gos.Listen(define.Http, int32(port))
	fmt.Printf("Init | Listen to port %d\n", port)

	if err != nil {
		fmt.Printf("Init | Failed to listen port %d, error: %+v\n", port, err)
		return
	}

	s = &MainServer{}
	s.HttpAnswer = anser.(*ans.HttpAnser)
	s.HttpHandler(s.HttpAnswer.Router)
	fmt.Printf("Init | Http Anser 伺服器初始化完成\n")

	// ==================================================
	// 送往 Dba Server 的請求，將數據依序寫入緩存
	// ==================================================
	asker, err := gos.Bind(EDbaServer, "127.0.0.1", 1022, define.Tcp0)

	if err != nil {
		fmt.Printf("Init | Bind error: %+v\n", err)
		return
	}

	dbaAsker = asker.(*ask.Tcp0Asker)
	dbaAsker.SetWorkHandler(s.DbaHandler)
	fmt.Printf("Init | DbaServer Asker 伺服器初始化完成\n")
	fmt.Printf("Init | 伺服器初始化完成\n")

	// =============================================
	// 開始所有已註冊的監聽
	// =============================================
	gos.StartListen()
	fmt.Printf("Init | 開始所有已註冊的監聽\n")

	// =============================================
	// 開始所有已註冊的連線
	// =============================================
	err = gos.StartConnect()

	if err != nil {
		fmt.Printf("Init | 與 DbaServer 連線時發生錯誤, error: %+v\n", err)
		return
	}

	fmt.Printf("Init | 成功與 DbaServer 連線\n")
}

func Run() {
	var start time.Time
	var during, frameTime time.Duration = 0, 20 * time.Millisecond

	for {
		start = time.Now()

		gos.RunAns()
		gos.RunAsk()
		s.Run()

		during = time.Since(start)
		if during < frameTime {
			time.Sleep(frameTime - during)
		}
	}
}

type MainServer struct {
	HttpAnswer *ans.HttpAnser
}

func (s *MainServer) HttpHandler(router *ans.Router) {
	router.GET("/", func(c *ghttp.Context) {
		c.Json(200, ghttp.H{
			"index": 1,
			"msg":   "GET | /",
		})
	})
	router.POST("/", func(c *ghttp.Context) {
		c.Json(200, ghttp.H{
			"index": 2,
			"msg":   "POST | /",
		})
	})

	r1 := router.NewRouter("/abc")

	r1.GET("/get", func(c *ghttp.Context) {
		c.Json(200, ghttp.H{
			"index": 3,
			"msg":   "GET | /abc/get",
		})
	})
	r1.POST("/post", func(c *ghttp.Context) {
		c.Json(200, ghttp.H{
			"index": 4,
			"msg":   "POST | /abc/post",
		})
	})
	r1.POST("/delay_response", func(c *ghttp.Context) {
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

func (s *MainServer) Run() {

}
