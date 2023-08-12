package client

import (
	"fmt"
	"time"

	"github.com/j32u4ukh/gos"
	"github.com/j32u4ukh/gos/ans"
	"github.com/j32u4ukh/gos/ask"
	"github.com/j32u4ukh/gos/base/ghttp"
	"github.com/j32u4ukh/gos/define"
)

func Init() {
	fmt.Println("Starting http server...")
	asker, err := gos.Bind(0, "127.0.0.1", 1023, define.Http, nil, nil, nil)

	if err != nil {
		fmt.Printf("BindError: %+v\n", err)
		return
	}

	http := asker.(*ask.HttpAsker)
	fmt.Printf("http: %+v\n", http)
}

func Run() {
	req, err := ghttp.NewRequest(ghttp.MethodGet, "127.0.0.1:1023/abc/get", nil)

	if err != nil {
		fmt.Printf("NewRequestError: %+v\n", err)
		return
	}

	fmt.Printf("req: %+v\n", req)
	var site int32
	site, err = gos.SendRequest(req, func(c *ghttp.Context) {
		fmt.Printf("I'm Context, Query: %s\n", c.Query)
	})

	fmt.Printf("site: %d\n", site)

	if err != nil {
		fmt.Printf("SendRequestError: %+v\n", err)
		return
	}

	var start time.Time
	var during, frameTime time.Duration = 0, 200 * time.Millisecond

	for {
		start = time.Now()

		gos.RunAsk()

		during = time.Since(start)
		if during < frameTime {
			time.Sleep(frameTime - during)
		}
	}
}

type Mgr struct {
}

func (m *Mgr) HttpHandler(router *ans.Router) {
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
}
