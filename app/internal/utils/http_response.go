package utils

import (
	"internal/define"

	"github.com/j32u4ukh/gos/base/ghttp"
)

func JsonResponse(c *ghttp.Context, code int32, contents ...any) string {
	status, _, msg := define.ErrorMessage(code, contents...)
	c.Json(status, ghttp.H{
		"code": code,
		"msg":  msg,
	})
	return msg
}
