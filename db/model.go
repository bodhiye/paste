package db

import (
	"time"
)

type PasteEntry struct {
	Key       string    `json:"key" bson:"key"`
	Langtype  string    `json:"langtype" bson:"langtype"`
	Content   string    `json:"content" bson:"content"`
	Password  string    `json:"password,omitempty" bson:"password,omitempty"`
	ClientIP  string    `json:"client_ip" bson:"client_ip"`
	Once      bool      `json:"once,omitempty" bson:"once,omitempty"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	ExpireAt  time.Time `json:"expire_at,omitempty" bson:"expire_at,omitempty"`
}
