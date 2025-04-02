package db

import (
	"time"
)

type PasteEntry struct {
    Key       string    `json:"key" bson:"key"`             // 唯一标识
    Langtype  string    `json:"langtype" bson:"langtype"`   // 代码语言类型
    Content   string    `json:"content" bson:"content"`     // 代码内容
    Password  string    `json:"-" bson:"password,omitempty"` // 密码保护
    ClientIP  string    `json:"-" bson:"client_ip"`         // 客户端 IP
    Once      bool      `json:"-" bson:"once,omitempty"`    // 是否一次性阅读
    CreatedAt time.Time `json:"-" bson:"created_at"`        // 创建时间
    ExpireAt  time.Time `json:"-" bson:"expire_at,omitempty"` // 过期时间
}
