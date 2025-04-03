package db

import (
	"time"

	"paste.org.cn/paste/server/proto"
)

type PasteEntry struct {
	Key         string            `json:"key" bson:"key"`                         // 唯一标识
	Title       string            `json:"title" bson:"title"`                     // 分享标题
	Description string            `json:"description" bson:"description"`         // 分享描述
	Snippets    []proto.Snippet   `json:"snippets" bson:"snippets,omitempty"`     // 多段代码内容
	ImageURLs   []proto.ImageInfo `json:"image_urls" bson:"image_urls,omitempty"` //多张截图分享内容
	Password    string            `json:"-" bson:"password,omitempty"`            // 密码保护
	ClientIP    string            `json:"-" bson:"client_ip"`                     // 客户端 IP
	Once        bool              `json:"-" bson:"once,omitempty"`                // 是否一次性阅读
	CreatedAt   time.Time         `json:"-" bson:"created_at"`                    // 创建时间
	ExpireAt    time.Time         `json:"-" bson:"expire_at,omitempty"`           // 过期时间
}
