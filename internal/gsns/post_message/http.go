package pm

import (
	"fmt"

	"github.com/j32u4ukh/gos/ans"
	"github.com/j32u4ukh/gos/base/ghttp"
)

// TODO: HTTP 請求處理過程中若失敗，要返回錯誤訊息給客戶端，而非印出日誌或直接返回
func (m *PostMessageMgr) HttpHandler(router *ans.Router) {

	router.GET("/", func(c *ghttp.Context) {
		fmt.Printf("Hello context %d\n", c.GetId())
	})
}
