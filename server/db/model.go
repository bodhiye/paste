package db

import (
	"time"
)

type PasteEntry struct {
	Key       string    `json:"key" bson:"key"`
	Langtype  string    `json:"langtype" bson:"langtype"`
	Content   string    `json:"content" bson:"content"`
	Password  string    `json:"-" bson:"password,omitempty"`
	ClientIP  string    `json:"-" bson:"client_ip"`
	Once      bool      `json:"-" bson:"once,omitempty"`
	CreatedAt time.Time `json:"-" bson:"created_at"`
	ExpireAt  time.Time `json:"-" bson:"expire_at,omitempty"`
}
