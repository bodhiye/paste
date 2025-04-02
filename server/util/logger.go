package util

import (
	"context"
	"fmt"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetReportCaller(true) // 让 logrus 在输出日志时自动包含调用日志记录函数的位置信息（文件名和行号）
	// 设置日志的输出格式
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true, // 日志中包含完整的时间戳
		// 自定义方法，用于美化调用者信息的显示格式
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line) // 函数名() 文件名:行号
		},
	})
}

// 定义一个非导出结构体，作为 context.WithValue 的键
type (
	loggerKey struct{}
)

func GetLogger(ctx context.Context) *logrus.Entry {
	logger := ctx.Value(loggerKey{})

	// 创建一个基于logrus标准（全局）日志记录器 (logrus.StandardLogger()) 的新 *logrus.Entry 并返回
	if logger == nil {
		return logrus.NewEntry(logrus.StandardLogger())
	}	

	return logger.(*logrus.Entry)
}

// withLogger 将一个 *logrus.Entry 实例添加到 context 中，并返回一个新的 context
func withLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

// EnsureWithLogger 确保 context 中存在一个 *logrus.Entry 实例
// 如果 context 中已存在 logger，则直接返回该 context 和 logger
// 如果不存在，则创建一个新的 logger（包含请求 ID），将其添加到 context 中
// 然后返回新的 context 和新创建的 logger
func EnsureWithLogger(ctx context.Context) (context.Context, *logrus.Entry) {
	logger := ctx.Value(loggerKey{})

	if logger == nil {
		var rid string
		ctx, rid = EnsureWithReqId(ctx)

		// .WithField(REQID, rid) 为这个新的 logger 实例添加了一个固定的"ReqID"请求ID字段
		newLoggerEntry := logrus.NewEntry(logrus.StandardLogger()).WithField(REQID, rid)

		// 将新创建的、带有请求 ID 的 logger 添加到 context 中
		ctx = withLogger(ctx, newLoggerEntry)
		logger = newLoggerEntry
	}
	return ctx, logger.(*logrus.Entry)
}
