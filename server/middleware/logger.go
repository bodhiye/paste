package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"paste.org.cn/paste/server/util"
)

// 使用gin框架提供的 LoggerWithFormatter 函数创建的自定义日志中间件,该中间件定义了日志的输出格式
var LogInfo = gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
	return fmt.Sprintf("%s [%s][%s] %s - \"%s %s %s %d %s \"%s\" %s\"\n",
		param.TimeStamp.Format("2006/01/02 15:04:05.999999"), // 请求时间，格式为 "年/月/日 时:分:秒.微秒"
		param.Keys[util.REQID],  // 请求唯一标识符
		"INFO",   // 日志级别
		param.ClientIP,  // 客户端 IP 地址
		param.Method,    // http请求方法
		param.Path,      // 请求路径
		param.Request.Proto,   // 请求的协议版本，如“http/1.1"
		param.StatusCode,      // 响应状态码
		param.Latency,        // 请求处理的延长时间
		param.Request.UserAgent(),  // 客户端的 user-agent 信息
		param.ErrorMessage,  // 错误信息
	)
})

// 中间件，用于为每个请求生成或提取唯一的请求 ID。
func ReqID(c *gin.Context) {
    // 从请求头中获取名为 util.REQID 的请求 ID
    reqid := c.Request.Header.Get(util.REQID)
    if reqid == "" {
        // 如果请求头中没有提供请求 ID，则生成一个新的请求 ID
        reqid = util.GenReqID()
    }
    // 将请求 ID 设置到 Gin 的上下文中，供后续处理使用
    c.Set(util.REQID, reqid)
    // 将请求 ID 添加到响应头中，方便客户端获取
    c.Header(util.REQID, reqid)
    // 继续处理请求的后续流程
    c.Next()
}
