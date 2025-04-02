package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"paste.org.cn/paste/server/db"
	"paste.org.cn/paste/server/service"
)

// 注册路由
func Init(r *gin.Engine, pasteDB db.Paste) {
	paste := &service.Paste{
		Paste: pasteDB,
	}

	r.POST("/v1/paste", paste.PostPaste) //创建分享内容
	r.POST("/v1/paste/once", paste.PostPasteOnce) //创建一次性分享内容
	r.GET("/v1/paste/:key", paste.GetPaste) //获取分享内容

	// health check
	r.Any("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "paste ok!")
	})
}
