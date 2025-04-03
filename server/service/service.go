package service

import (
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"

	"paste.org.cn/paste/server/db"
	"paste.org.cn/paste/server/proto"
	"paste.org.cn/paste/server/util"
)

type Paste struct {
	db.Paste
}

// 创建分享内容
func (p *Paste) PostPaste(c *gin.Context) {
	var (
		ctx, log = util.EnsureWithLogger(c)
		req      proto.PostPasteReq
		length   int
	)

	err := c.BindJSON(&req)
	if err != nil {
		log.Errorf("BindJSON failed: %+v", err)
		c.JSON(http.StatusBadRequest, proto.PostPasteResp{
			Code:    http.StatusBadRequest,
			Message: proto.InvalidArgs,
		})
		return
	}

	// 计算内容的字符数（而不是字节数，对多字节字符友好）
	length = utf8.RuneCountInString(req.Content)
	if length > 100000 {
		log.Errorf("Content is too long: %d", length)
		c.JSON(http.StatusBadRequest, proto.PostPasteResp{
			Code:    http.StatusBadRequest,
			Message: proto.TooManyContent,
		})
		return
	}

	entry := db.PasteEntry{
		Langtype:  req.Langtype,
		Content:   req.Content,
		ClientIP:  c.ClientIP(),
		CreatedAt: time.Now(),
	}
	if req.Password != "" {
		entry.Password = util.String2bcrypt(req.Password)
	}
	if req.ExpireDate > 0 {
		entry.ExpireAt = time.Now().Add(time.Second * time.Duration(req.ExpireDate))
	}

	key, err := p.Paste.Set(ctx, entry)
	if err != nil {
		log.Errorf("Failed to insert entry into database: %+v", err)
		c.JSON(http.StatusBadRequest, proto.PostPasteResp{
			Code:    http.StatusBadRequest,
			Message: proto.PasteFailed,
		})
		return
	}
	c.JSON(http.StatusCreated, proto.PostPasteResp{
		Code: http.StatusCreated,
		Key:  key,
	})
}

// 创建一次性分享内容
func (p *Paste) PostPasteOnce(c *gin.Context) {
	var (
		ctx, log = util.EnsureWithLogger(c)
		req      proto.PostPasteReq
		length   int
	)

	err := c.BindJSON(&req)
	if err != nil {
		log.Errorf("BindJSON failed: %+v", err)
		c.JSON(http.StatusBadRequest, proto.PostPasteResp{
			Code:    http.StatusBadRequest,
			Message: proto.InvalidArgs,
		})
		return
	}

	length = utf8.RuneCountInString(req.Content)
	if length > 10000 {
		log.Errorf("Content is too long: %d", length)
		c.JSON(http.StatusBadRequest, proto.PostPasteResp{
			Code:    http.StatusBadRequest,
			Message: proto.TooManyContent,
		})
		return
	}

	entry := db.PasteEntry{
		Langtype:  req.Langtype,
		Content:   req.Content,
		ClientIP:  c.ClientIP(),
		Once:      true,
		CreatedAt: time.Now(),
	}

	if req.Password != "" {
		entry.Password = util.String2bcrypt(req.Password)
	}

	key, err := p.Paste.Set(ctx, entry)
	if err != nil {
		log.Errorf("Failed to insert entry into database: %+v", err)
		c.JSON(http.StatusBadRequest, proto.PostPasteResp{
			Code:    http.StatusBadRequest,
			Message: proto.PasteFailed,
		})
		return
	}
	c.JSON(http.StatusOK, proto.PostPasteResp{
		Code: http.StatusCreated,
		Key:  key,
	})
}

// 获取分享内容
func (p *Paste) GetPaste(c *gin.Context) {
	var (
		ctx, log = util.EnsureWithLogger(c)
		// 从URL路径参数获取key，从URL查询参数中获取passwor
		key, password = c.Param("key"), c.Query("password")
	)

	entry, err := p.Paste.Get(ctx, key, password)
	if err != nil {
		log.Errorf("Failed to get entry: %+v", err)
		if err.Error() == proto.WrongPassword {
			c.JSON(http.StatusOK, proto.GetPasteResp{
				Code:    http.StatusUnauthorized,
				Message: proto.WrongPassword,
			})
		} else if err.Error() == proto.ContentExpired {
			c.JSON(http.StatusOK, proto.GetPasteResp{
				Code:    http.StatusLocked,
				Message: proto.ContentExpired,
			})
		} else {
			c.JSON(http.StatusOK, proto.GetPasteResp{
				Code:    http.StatusBadRequest,
				Message: proto.GetPasteFailed,
			})
		}
		return
	}
	c.JSON(http.StatusOK, proto.GetPasteResp{
		Code:     http.StatusOK,
		Langtype: entry.Langtype,
		Content:  entry.Content,
	})
}
