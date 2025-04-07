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

	r.POST("/v1/paste", paste.PostPaste)          //创建文本片段分享内容
	r.POST("/v1/paste/once", paste.PostPasteOnce) //创建一次性文本片段分享内容
	r.GET("/v1/paste/:key", paste.GetPaste)       //获取文本片段分享内容

	r.POST("/v1/image", paste.PostImage)                  // 创建图片分享内容
	r.POST("/v1/image/once", paste.PostImageOnce)         //创建一次性图片分享内容
	r.GET("/v1/image/:key", paste.GetImage) 	          // 获取图片分享内容

	// health check
	r.Any("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "paste ok!")
	})
}
