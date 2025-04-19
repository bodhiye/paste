package proto

// Snippet 结构体表示片段类型
type Snippet struct {
	Langtype string `json:"langtype" bson:"langtype"` // 代码/文本（如 "go", "python", "text"，"markdown"等）
	Content  string `json:"content" bson:"content"`    // 片段内容，最大支持十万个字符
}

// ImageFile 结构体表示图片类型
type ImageFile struct {
	Filename    string `json:"filename" bson:"filename"`             // 图片原始文件名
	URL         string `json:"url" bson:"url"`                            // 文件访问URL路径
	Size        int64  `json:"size" bson:"size"`                         // 文件大小（字节）
	ContentType string `json:"content_type" bson:"content_type"` // 文件MIME类型
}

// PostPasteReq 结构体表示创建分享请求的请求体
type PostPasteReq struct {
	Title       string      `form:"title" json:"title"`                               // 分享标题
	Description string      `form:"description" json:"description"`                   // 分享描述
	Snippets    []Snippet   `form:"-" json:"snippets"`                         // 多段代码内容，前后端都需要限制片段字符内容长度和数量
	Images      []ImageFile `form:"-" json:"images,omitempty"`         // 多张截图分享内容，前后端需要限制图片大小和数量（10M,5张）
	Password    string      `form:"password,omitempty" json:"password,omitempty"`     // 访问分享内容的可选密码（omitempty 表示如果为空则不序列化）
	ExpireDate  int64       `form:"expireDate,omitempty" json:"expireDate,omitempty"` // 过期时间戳，单位为秒（可选字段）
}

// PostPasteResp 结构体表示创建分享请求的响应体
type PostPasteResp struct {
	Code    int    `json:"code"`              // 状态码
	Key     string `json:"key"`     // 分享内容的唯一标识符
	Message string `json:"message,omitempty"` // 服务器返回的消息（可选）
}

// GetPasteResp 结构体表示获取分享请求的响应体
type GetPasteResp struct {
	Code     int         `json:"code"`               // 状态码
	Snippets []Snippet   `json:"snippets"` 			// 返回多个片段
	Images   []ImageFile `json:"images,omitempty"`   // 返回多张图片 (可选)
	Message  string      `json:"message,omitempty"`  // 服务器返回的消息（可选）
}
