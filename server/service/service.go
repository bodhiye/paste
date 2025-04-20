package service

import (
	"encoding/json"
	"fmt"
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
		err      error
	)

	// 使用 ShouldBind 绑定基本字段
	if err = c.ShouldBind(&req); err != nil {
		log.Errorf("绑定请求数据失败: %+v", err)
		c.JSON(http.StatusBadRequest, proto.PostPasteResp{
			Code:    http.StatusBadRequest,
			Message: proto.ErrInvalidArgs,
		})
		return
	}

	// 需要手动解析snippets
	raw := c.PostForm("snippets")
	if err = json.Unmarshal([]byte(raw), &req.Snippets); err != nil {
		log.Errorf("解析snippets失败: %+v", err)
		c.JSON(http.StatusBadRequest, proto.PostPasteResp{
			Code:    http.StatusBadRequest,
			Message: proto.ErrInvalidArgs,
		})
		return
	}

	if len(req.Snippets) == 0 && len(req.Images) == 0 {
		log.Errorf("内容为空")
		c.JSON(http.StatusBadRequest, proto.PostPasteResp{
			Code:    http.StatusBadRequest,
			Message: proto.ErrInvalidArgs,
		})
		return
	}

	// 验证代码片段数量
	if len(req.Snippets) > util.LimitConfig.SnippetsCount() {
		log.Errorf("代码片段数量过多: %d", len(req.Snippets))
		c.JSON(http.StatusBadRequest, proto.PostPasteResp{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf(proto.ErrTooManyCount, util.LimitConfig.SnippetsCount()),
		})
		return
	}

	// 验证代码片段内容
	for _, snippet := range req.Snippets {
		length := utf8.RuneCountInString(snippet.Content)
		if length > util.LimitConfig.SnippetsLength() {
			log.Errorf("内容过长: %d", length)
			c.JSON(http.StatusBadRequest, proto.PostPasteResp{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf(proto.ErrTooManyContent, util.LimitConfig.SnippetsLength()),
			})
			return
		}
	}

	req.Images, err = util.UploadImages(c, log)
	if err != nil {
		log.Errorf("获取图片失败: %+v", err)
		c.JSON(http.StatusBadRequest, proto.PostPasteResp{
			Code:    http.StatusBadRequest,
			Message: proto.ErrUploadFailed,
		})
		return
	}

	entry := db.PasteEntry{
		Title:       req.Title,
		Description: req.Description,
		Snippets:    req.Snippets,
		Images:      req.Images,
		ClientIP:    c.ClientIP(),
		CreatedAt:   time.Now(),
	}

	// 设置密码（如果有）
	if req.Password != "" {
		entry.Password = util.String2bcrypt(req.Password)
	}

	// 设置过期时间（如果有）
	if req.ExpireDate > 0 {
		entry.ExpireAt = time.Now().Add(time.Second * time.Duration(req.ExpireDate))
	}

	// 保存到数据库
	key, err := p.Paste.Set(ctx, entry)
	if err != nil {
		log.Errorf("插入数据库失败: %+v", err)
		c.JSON(http.StatusBadRequest, proto.PostPasteResp{
			Code:    http.StatusBadRequest,
			Message: proto.ErrPasteFailed,
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusCreated, proto.PostPasteResp{
		Code: http.StatusCreated,
		Key:  key,
	})
}

// 创建一次性分享内容
func (p *Paste) PostPasteOnce(c *gin.Context) {
	// 复用相同逻辑，仅添加一次性标记
	var (
		ctx, log = util.EnsureWithLogger(c)
		req      proto.PostPasteReq
		err      error
	)

	// 使用 ShouldBind 绑定基本字段
	if err = c.ShouldBind(&req); err != nil {
		log.Errorf("绑定请求数据失败: %+v", err)
		c.JSON(http.StatusBadRequest, proto.PostPasteResp{
			Code:    http.StatusBadRequest,
			Message: proto.ErrInvalidArgs,
		})
		return
	}

	// 需要手动解析snippets
	raw := c.PostForm("snippets")
	if err = json.Unmarshal([]byte(raw), &req.Snippets); err != nil {
		log.Errorf("解析snippets失败: %+v", err)
		c.JSON(http.StatusBadRequest, proto.PostPasteResp{
			Code:    http.StatusBadRequest,
			Message: proto.ErrInvalidArgs,
		})
		return
	}

	// 验证代码片段内容
	for _, snippet := range req.Snippets {
		length := utf8.RuneCountInString(snippet.Content)
		if length > 10000 { // 一次性分享内容限制更严格
			log.Errorf("内容过长: %d", length)
			c.JSON(http.StatusBadRequest, proto.PostPasteResp{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf(proto.ErrTooManyContent, 10000),
			})
			return
		}
	}

	// 获取图片
	req.Images, err = util.UploadImages(c, log)
	if err != nil {
		log.Errorf("获取图片失败: %+v", err)
		c.JSON(http.StatusBadRequest, proto.PostPasteResp{
			Code:    http.StatusBadRequest,
			Message: proto.ErrUploadFailed,
		})
		return
	}

	// 创建数据库记录
	entry := db.PasteEntry{
		Title:       req.Title,
		Description: req.Description,
		Snippets:    req.Snippets,
		Images:      req.Images,
		ClientIP:    c.ClientIP(),
		Once:        true, // 标记为一次性
		CreatedAt:   time.Now(),
	}

	// 设置密码（如果有）
	if req.Password != "" {
		entry.Password = util.String2bcrypt(req.Password)
	}

	// 保存到数据库
	key, err := p.Paste.Set(ctx, entry)
	if err != nil {
		log.Errorf("插入数据库失败: %+v", err)
		c.JSON(http.StatusBadRequest, proto.PostPasteResp{
			Code:    http.StatusBadRequest,
			Message: proto.ErrPasteFailed,
		})
		return
	}

	// 返回成功响应
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
		log.Errorf("获取分享内容失败: %+v", err)
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

	if entry.Images == nil {
		entry.Images = []proto.ImageFile{}
	}

	// 返回成功响应
	c.JSON(http.StatusOK, proto.GetPasteResp{
		Code:     http.StatusOK,
		Snippets: entry.Snippets,
		Images:   entry.Images,
	})
}
