package main

import (
	"github.com/j32u4ukh/gos/ans"
	"github.com/j32u4ukh/gos/base/ghttp"
)

type Mgr struct {
	Answer *ans.HttpAnser
}

func (m *Mgr) Handler(router *ans.Router) {
	router.GET("/", func(c *ghttp.Context) {
		c.Json(200, ghttp.H{
			"index": 1,
			"msg":   "GET | /",
		})
		m.Answer.Send(c)
	})
	router.POST("/", func(c *ghttp.Context) {
		c.Json(200, ghttp.H{
			"index": 2,
			"msg":   "POST | /",
		})
		m.Answer.Send(c)
	})

	r1 := router.NewRouter("/abc")

	r1.GET("/get", func(c *ghttp.Context) {
		c.Json(200, ghttp.H{
			"index": 3,
			"msg":   "GET | /abc/get",
		})
		m.Answer.Send(c)
	})
	r1.POST("/post", func(c *ghttp.Context) {
		c.Json(200, ghttp.H{
			"index": 4,
			"msg":   "POST | /abc/post",
		})
		m.Answer.Send(c)
	})
}
