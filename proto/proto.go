package proto

import "time"

type PostPasteReq struct {
	Langtype   string        `json:"langtype"`
	Content    string        `json:"content"` // Support up to 10000 characters.
	Password   string        `json:"password,omitempty"`
	ExpireDate time.Duration `json:"expireDate,omitempty"`
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
