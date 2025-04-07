package proto

import "time"

// Snippet 结构体表示片段类型
type Snippet struct {
	Langtype string `json:"langtype" form:"langtype" bson:"langtype"` // 代码/文本（如 "go", "python", "text"，"markdown"等）
	Content  string `json:"content" form:"content" bson:"content"`    // 片段内容，最大支持十万个字符
}

// ImageInfo 结构体图片信息
type ImageInfo struct {
	FileName    string    `json:"file_name" bson:"file_name"`                 // 原文件名
	ObjectKey   string    `json:"object_key" bson:"object_key"`               // 对象存储中的键
	MimeType    string    `json:"mime_type" bson:"mime_type"`                 // 文件MIME类型
	Width       int       `json:"width" bson:"width"`                         // 宽度(px)
	Height      int       `json:"height" bson:"height"`                       // 高度(px)
	SizeBytes   int64     `json:"size_bytes" bson:"size_bytes"`               // 文件大小(字节)
	URL         string    `json:"url" bson:"url"`                               // 临时使用，不存储到数据库
	UploadedAt  time.Time `json:"uploaded_at" bson:"uploaded_at"`             // 上传时间
	AccessCount int64     `json:"access_count,omitempty" bson:"access_count"` // 访问计数
}

// PostPasteReq 结构体表示创建新的分享请求的请求体
type PostPasteReq struct {
	Title       string      `form:"title" json:"title"`                           // 分享标题
	Description string      `form:"description" json:"description"`               // 分享描述
	Snippets    []Snippet   `form:"snippets,omitempty" json:"snippets,omitempty"` // 多段代码内容，前后端都需要限制片段字符内容长度和数量
	Images      []ImageInfo `form:"-" json:"images,omitempty"`                    // 多张截图分享内容，前后端需要限制图片大小和数量（10M,5张）
	Password    string      `form:"password,omitempty" json:"password,omitempty"` // 访问分享内容的可选密码（omitempty 表示如果为空则不序列化）
	ExpireIn    string      `form:"expirein,omitempty" json:"expirein,omitempty"` // 过期时间，格式如："1h", "1d", "7d", "30d", "365d"
}

// PostPasteResp 结构体表示创建分享请求的响应体
type PostPasteResp struct {
	Code    int    `json:"code"`              // 状态码
	Key     string `json:"key,omitempty"`     // 分享内容的唯一标识符（可选）
	Message string `json:"message,omitempty"` // 服务器返回的消息（可选）
}

// GetPasteResp 结构体表示获取分享内容的响应体
type GetPasteResp struct {
	Code        int         `json:"code"`                  // 状态码
	Snippets    []Snippet   `json:"snippets,omitempty"`    // 返回多个片段
	Images      []ImageInfo `json:"images,omitempty"`      // 返回多张图片
	Message     string      `json:"message,omitempty"`     // 服务器返回的消息（可选）
	Title       string      `json:"title,omitempty"`       // 分享标题（可选）
	Description string      `json:"description,omitempty"` // 分享描述（可选）
}
