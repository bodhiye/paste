package util

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"time"
)

const REQID string = "ReqID"

var pid = uint32(time.Now().UnixNano() % 4294967291)

type (
	reqIDKey struct{}
)

func GenReqID() string {
	var b [12]byte
	binary.LittleEndian.PutUint32(b[:], pid)
	binary.LittleEndian.PutUint64(b[4:], uint64(time.Now().UnixNano()))
	return base64.URLEncoding.EncodeToString(b[:])
}

func GetReqID(ctx context.Context) string {
	reqID := ctx.Value(reqIDKey{})

	if reqID == nil {
		reqID = ctx.Value(REQID)
		if reqID == nil {
			return GenReqID()
		}
	}

	return reqID.(string)
}

func WithReqId(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, reqIDKey{}, reqID)
}

func EnsureWithReqId(ctx context.Context) (context.Context, string) {
	reqID := ctx.Value(reqIDKey{})

	if reqID == nil {
		reqID = ctx.Value(REQID)
		if reqID == nil {
			reqID = GenReqID()
			ctx = WithReqId(ctx, reqID.(string))
		}
	}

	return ctx, reqID.(string)
}
