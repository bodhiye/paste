package service

import (
	"fmt"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"

	"paste.org.cn/paste/server/db"
	"paste.org.cn/paste/server/proto"
	"paste.org.cn/paste/server/util"
)

// 创建分享内容
func (p *Paste) PostPaste(c *gin.Context) {
	var (
		ctx, log = util.EnsureWithLogger(c)
		req      proto.PostPasteReq
		length   int
	)

	err := c.ShouldBind(&req)
	if err != nil {
		log.Errorf("ShouldBind failed: %+v", err)
		c.JSON(http.StatusBadRequest, proto.PostPasteResp{
			Code:    http.StatusBadRequest,
			Message: proto.ErrInvalidArgs,
		})
		return
	}

	// 计算内容的数量
	if len(req.Snippets) > util.LimitConfig.SnippetCount() {
		log.Errorf("Count is too many: %d", len(req.Snippets))
		c.JSON(http.StatusBadRequest, proto.PostPasteResp{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf(proto.ErrTooManyCount, util.LimitConfig.SnippetCount()),
		})
		return
	}

	// 计算内容的字符数（而不是字节数，对多字节字符友好）
	for _, val := range req.Snippets {
		length = utf8.RuneCountInString(val.Content)
		if length > util.LimitConfig.SnippetLength() {
			log.Errorf("Content is too long: %d", length)
			c.JSON(http.StatusBadRequest, proto.PostPasteResp{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf(proto.ErrTooManyContent, util.LimitConfig.SnippetLength()),
			})
			return
		}
	}

	// 处理过期时间
	if req.ExpireIn != "" {
		req.ExpireIn = "24h"
	}
	duration, err := time.ParseDuration(req.ExpireIn)
	if err != nil {
		log.Errorf("Invalid expiration format: %v", err)
		c.JSON(http.StatusBadRequest, proto.PostPasteResp{
			Code:    http.StatusBadRequest,
			Message: proto.ErrInvalidArgs,
		})
		return
	}

	entry := db.PasteEntry{
		Title:       req.Title,
		Description: req.Description,
		Snippets:    req.Snippets,
		ClientIP:    c.ClientIP(),
		CreatedAt:   time.Now(),
		ExpireAt: time.Now().Add(duration),
	}
	if req.Password != "" {
		entry.Password = util.String2bcrypt(req.Password)
	}

	key, err := p.Paste.Set(ctx, entry)
	if err != nil {
		log.Errorf("Failed to insert entry into database: %+v", err)
		c.JSON(http.StatusBadRequest, proto.PostPasteResp{
			Code:    http.StatusBadRequest,
			Message: proto.ErrPasteFailed,
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

	err := c.ShouldBind(&req)
	if err != nil {
		log.Errorf("ShouldBind failed: %+v", err)
		c.JSON(http.StatusBadRequest, proto.PostPasteResp{
			Code:    http.StatusBadRequest,
			Message: proto.ErrInvalidArgs,
		})
		return
	}

	// 计算内容的数量
	if len(req.Snippets) > util.LimitConfig.SnippetCount() {
		log.Errorf("Count is too many: %d", len(req.Snippets))
		c.JSON(http.StatusBadRequest, proto.PostPasteResp{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf(proto.ErrTooManyCount, util.LimitConfig.SnippetCount()),
		})
		return
	}

	// 计算内容的字符数（而不是字节数，对多字节字符友好）
	for _, val := range req.Snippets {
		length = utf8.RuneCountInString(val.Content)
		if length > util.LimitConfig.SnippetLength() {
			log.Errorf("Content is too long: %d", length)
			c.JSON(http.StatusBadRequest, proto.PostPasteResp{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf(proto.ErrTooManyContent, util.LimitConfig.SnippetLength()),
			})
			return
		}
	}

	entry := db.PasteEntry{
		Title:       req.Title,
		Description: req.Description,
		Snippets:    req.Snippets,
		ClientIP:    c.ClientIP(),
		Once:        true,
		CreatedAt:   time.Now(),
	}

	if req.Password != "" {
		entry.Password = util.String2bcrypt(req.Password)
	}

	key, err := p.Paste.Set(ctx, entry)
	if err != nil {
		log.Errorf("Failed to insert entry into database: %+v", err)
		c.JSON(http.StatusBadRequest, proto.PostPasteResp{
			Code:    http.StatusBadRequest,
			Message: proto.ErrPasteFailed,
		})
		return
	}
	c.JSON(http.StatusCreated, proto.PostPasteResp{
		Code: http.StatusCreated,
		Key:  key,
	})
}

// 获取分享内容
func (p *Paste) GetPaste(c *gin.Context) {
	var (
		ctx, log = util.EnsureWithLogger(c)
		// 从URL路径参数获取key，从URL查询参数中获取password
		key, password = c.Param("key"), c.Query("password")
	)

	entry, err := p.Paste.Get(ctx, key, password)
	if err != nil {
		log.Errorf("Failed to get entry: %+v", err)
		if err.Error() == proto.ErrWrongPassword {
			c.JSON(http.StatusOK, proto.GetPasteResp{
				Code:    http.StatusUnauthorized,
				Message: proto.ErrWrongPassword,
			})
		} else if err.Error() == proto.ErrContentExpired {
			c.JSON(http.StatusOK, proto.GetPasteResp{
				Code:    http.StatusLocked,
				Message: proto.ErrContentExpired,
			})
		} else {
			c.JSON(http.StatusOK, proto.GetPasteResp{
				Code:    http.StatusBadRequest,
				Message: proto.ErrGetPasteFailed,
			})
		}
		return
	}

	c.JSON(http.StatusOK, proto.GetPasteResp{
		Code:        http.StatusOK,
		Snippets:    entry.Snippets,
		Title:       entry.Title,
		Description: entry.Description,
	})
}
