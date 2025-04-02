package util

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"time"
)

// 定义常量 REQID，用于在上下文作为键
const REQID string = "ReqID"

// 生成一个基于当前时间戳的进程ID，用于生成唯一的请求ID
var pid = uint32(time.Now().UnixNano() % 4294967291)

// 定义一个空的结构体类型 reqIDKey，用于在上下文中作为键
type (
	reqIDKey struct{}
)

// GenReqID 生成一个新的请求ID
// 该ID是由进程ID和当前的纳秒时间戳组成的，并使用Base64 URL编码格式返回
func GenReqID() string {
	// 创建一个长度为12字节的字节数组
	var b [12]byte
	
	// 将进程ID (pid) 写入字节数组的前4字节
	binary.LittleEndian.PutUint32(b[:], pid)
	
	// 将当前的时间戳（纳秒）写入字节数组的第5到12字节
	binary.LittleEndian.PutUint64(b[4:], uint64(time.Now().UnixNano()))
	
	// 将字节数组编码为Base64 URL格式字符串
	return base64.URLEncoding.EncodeToString(b[:])
}

// GetReqID 从给定的上下文中获取请求ID
// 如果上下文中没有请求ID，生成并返回一个新的请求ID
func GetReqID(ctx context.Context) string {
	// 从上下文中获取请求ID，首先尝试使用 reqIDKey
	reqID := ctx.Value(reqIDKey{})

	// 如果没有找到 reqID（即 reqID 为空），尝试获取常量名为 "ReqID" 的值
	if reqID == nil {
		reqID = ctx.Value(REQID)
		if reqID == nil {
			// 如果上下文中仍然没有找到请求ID，生成一个新的请求ID
			return GenReqID()
		}
	}

	// 将上下文中的请求ID转换为字符串并返回
	return reqID.(string)
}

// WithReqId 将请求ID添加到上下文中
// 返回一个包含该请求ID的新上下文
func WithReqId(ctx context.Context, reqID string) context.Context {
	// 使用 reqIDKey 作为键，将请求ID值存储到上下文中
	return context.WithValue(ctx, reqIDKey{}, reqID)
}

// EnsureWithReqId 确保上下文中有一个请求ID
// 如果没有请求ID，则生成一个新的请求ID并将其添加到上下文中
// 返回更新后的上下文和请求ID
func EnsureWithReqId(ctx context.Context) (context.Context, string) {
	// 尝试从上下文中获取请求ID，首先使用 reqIDKey
	reqID := ctx.Value(reqIDKey{})

	// 如果上下文中没有请求ID，尝试获取常量名为 "ReqID" 的值
	if reqID == nil {
		reqID = ctx.Value(REQID)
		if reqID == nil {
			// 如果上下文中仍然没有找到请求ID，生成一个新的请求ID
			reqID = GenReqID()
			// 将新的请求ID添加到上下文中
			ctx = WithReqId(ctx, reqID.(string))
		}
	}

	// 返回更新后的上下文和请求ID
	return ctx, reqID.(string)
}