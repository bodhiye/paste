package proto

// Snippet 结构体表示片段类型
type Snippet struct {
	Langtype string `json:"langtype" form:"langtype" bson:"langtype"` // 代码/文本（如 "go", "python", "text"，"markdown"等）
	Content  string `json:"content" form:"content" bson:"content"`  // 片段内容，最大支持十万个字符
}

// ImageFile 结构体表示图片类型
type ImageFile struct {
	Name    string `form:"name" json:"name"`       // 图片名称
	Content []byte `form:"content" json:"content"` // 图片文件 (form tag 用于标识文件上传字段)
	Width   int    `form:"width" json:"width"`     // 宽度
	Height  int    `form:"height" json:"height"`   // 高度
	SizeMB  int    `form:"size" json:"size"`       // 图片大小（单位：MB）
}

// PostPasteReq 结构体表示创建新的分享请求的请求体
type PostPasteReq struct {
	Title       string       `form:"title" json:"title"`                               // 分享标题
	Description string       `form:"description" json:"description"`                   // 分享描述
	Snippets    []Snippet `form:"code_snippets,omitempty" json:"code_snippets"`     // 多段代码内容，前后端都需要限制片段字符内容长度和数量
	Images      []ImageFile  `form:"image_urls,omitempty" json:"image_urls"`           // 多张截图分享内容，前后端需要限制图片大小和数量（10M,5张）
	Password    string       `form:"password,omitempty" json:"password,omitempty"`     // 访问分享内容的可选密码（omitempty 表示如果为空则不序列化）
	ExpireDate  int64        `form:"expireDate,omitempty" json:"expireDate,omitempty"` // 过期时间戳，单位为秒（可选字段）
}

// PostPasteResp 结构体表示创建分享请求的响应体
type PostPasteResp struct {
	Code    int    `json:"code"`              // 状态码
	Key     string `json:"key,omitempty"`     // 分享内容的唯一标识符（可选）
	Message string `json:"message,omitempty"` // 服务器返回的消息（可选）
}

// ImageInfo 结构体图片信息
type ImageInfo struct {
	Name   string  `json:"name" bson:"name"`   // 图片名称
	URL    string  `json:"url" bson:"url"`    // 图片文件
	Width  int     `json:"width" bson:"width"`  // 宽度
	Height int     `json:"height" bson:"height"` // 高度
	SizeMB   int     `json:"size_mb" bson:"size_mb"`   // 图片大小（单位：MB）
}

// GetPasteResp 结构体表示获取分享内容的响应体
type GetPasteResp struct {
	Code      int         `json:"code"`                    // 状态码
	Snippets  []Snippet   `json:"code_snippets,omitempty"` // 返回多个片段
	ImageURLs []ImageInfo `json:"image_urls,omitempty"`    // 返回多张图片
	Message   string      `json:"message,omitempty"`       // 服务器返回的消息（可选）
}
