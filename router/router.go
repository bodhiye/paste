package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"paste.org.cn/paste/db"
	"paste.org.cn/paste/service"
)

func Init(r *gin.Engine, pasteDB db.Paste) {
	paste := &service.Paste{
		Paste: pasteDB,
	}

	r.POST("v1/paste", paste.PostPaste)
	r.POST("v1/paste/once", paste.PostPasteOnce)
	r.GET("v1/paste/:key", paste.GetPaste)

	// health check
	r.Any("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "paste ok!")
	})
}
