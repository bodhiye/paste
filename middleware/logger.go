package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"paste.org.cn/paste/util"
)

var LogInfo = gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
	return fmt.Sprintf("%s [%s][%s] %s - \"%s %s %s %d %s \"%s\" %s\"\n",
		param.TimeStamp.Format("2006/01/02 15:04:05.999999"),
		param.Keys[util.REQID],
		"INFO",
		param.ClientIP,
		param.Method,
		param.Path,
		param.Request.Proto,
		param.StatusCode,
		param.Latency,
		param.Request.UserAgent(),
		param.ErrorMessage,
	)
})

func ReqID(c *gin.Context) {
	reqid := c.Request.Header.Get(util.REQID)
	if reqid == "" {
		reqid = util.GenReqID()
	}
	c.Set(util.REQID, reqid)
	c.Header(util.REQID, reqid)
	c.Next()
}
