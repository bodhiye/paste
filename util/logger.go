package util

import (
	"context"
	"fmt"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
	})
}

type (
	loggerKey struct{}
)

func GetLogger(ctx context.Context) *logrus.Entry {
	logger := ctx.Value(loggerKey{})

	if logger == nil {
		return logrus.NewEntry(logrus.StandardLogger())
	}

	return logger.(*logrus.Entry)
}

func withLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

func EnsureWithLogger(ctx context.Context) (context.Context, *logrus.Entry) {
	logger := ctx.Value(loggerKey{})

	if logger == nil {
		var rid string
		ctx, rid = EnsureWithReqId(ctx)
		logger = logrus.NewEntry(logrus.StandardLogger()).WithField(REQID, rid)
		ctx = withLogger(ctx, logger.(*logrus.Entry))
	}

	return ctx, logger.(*logrus.Entry)
}
