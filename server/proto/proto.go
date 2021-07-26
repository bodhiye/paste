package proto

type PostPasteReq struct {
	Langtype   string `json:"langtype"`
	Content    string `json:"content"` // 最大支持十万个字符
	Password   string `json:"password,omitempty"`
	ExpireDate int64  `json:"expireDate,omitempty"`
}

type PostPasteResp struct {
	Code    int    `json:"code"`
	Key     string `json:"key,omitempty"`
	Message string `json:"message,omitempty"`
}

type GetPasteResp struct {
	Code     int    `json:"code"`
	Langtype string `json:"langtype,omitempty"`
	Content  string `json:"content,omitempty"`
	Message  string `json:"message,omitempty"`
}
