package gsns

import "github.com/j32u4ukh/gos/base/ghttp"

func (s *MainServer) register(c *ghttp.Context) {
	c.Json(200, ghttp.H{
		"index": 2,
		"msg":   "POST | /",
	})
}

func (s *MainServer) login(c *ghttp.Context) {
	c.Json(200, ghttp.H{
		"index": 2,
		"msg":   "POST | /",
	})
}

func (s *MainServer) logout(c *ghttp.Context) {
	c.Json(200, ghttp.H{
		"index": 2,
		"msg":   "POST | /",
	})
}
