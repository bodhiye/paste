package proto

// PostPasteReq 结构体表示创建新的分享请求的请求体
type PostPasteReq struct {
	Langtype   string `json:"langtype"`             // 代码/文本类型（如 "go", "python", "text"）
	Content    string `json:"content"`              // 分享的内容，最大支持十万个字符
	Password   string `json:"password,omitempty"`   // 访问分享内容的可选密码（omitempty 表示如果为空则不序列化）
	ExpireDate int64  `json:"expireDate,omitempty"` // 过期时间戳，单位为秒（可选字段）
}

// PostPasteResp 结构体表示创建分享请求的响应体
type PostPasteResp struct {
	Code    int    `json:"code"`            // 状态码
	Key     string `json:"key,omitempty"`   // 分享内容的唯一标识符（可选）
	Message string `json:"message,omitempty"` // 服务器返回的消息（可选）
}

// GetPasteResp 结构体表示获取分享内容的响应体
type GetPasteResp struct {
	Code     int    `json:"code"`            // 状态码
	Langtype string `json:"langtype,omitempty"` // 代码/文本类型（可选）
	Content  string `json:"content,omitempty"`  // 分享的内容（可选）
	Message  string `json:"message,omitempty"`  // 服务器返回的消息（可选）
}