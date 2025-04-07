package service

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"paste.org.cn/paste/server/db"
	"paste.org.cn/paste/server/proto"
	"paste.org.cn/paste/server/util"
)

// 创建图片分享内容
func (p *Paste) PostImage(c *gin.Context) {
	var (
		ctx, log  = util.EnsureWithLogger(c)
		req       proto.PostPasteReq
		images    = []proto.ImageInfo{}
		objectKey string
		url       string
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

	// 处理过期时间
	if req.ExpireIn == "" || strings.TrimSpace(req.ExpireIn) == "" {
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

	// 处理图片文件上传
	form, err := c.MultipartForm()
	if err != nil {
		log.Errorf("Failed to get multipart form: %+v", err)
		c.JSON(http.StatusBadRequest, proto.PostPasteResp{
			Code:    http.StatusBadRequest,
			Message: proto.ErrInvalidArgs,
		})
		return
	}

	files := form.File["images"]

	// 检查图片数量
	if len(files) == 0 {
		log.Error("No images uploaded")
		c.JSON(http.StatusBadRequest, proto.PostPasteResp{
			Code:    http.StatusBadRequest,
			Message: proto.ErrInvalidArgs,
		})
		return
	}

	if len(files) > util.LimitConfig.ImageCount() {
		log.Errorf("Image count exceeds the allowed limit: %d", len(files))
		c.JSON(http.StatusBadRequest, proto.PostPasteResp{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf(proto.ErrTooManyCount, util.LimitConfig.ImageCount()),
		})
		return
	}

	for _, file := range files {
		// 1.检查文件大小
		if file.Size > util.LimitConfig.ImageSize() {
			log.Errorf("Image size exceeds the allowed limit: %d bytes", file.Size)
			c.JSON(http.StatusBadRequest, proto.PostPasteResp{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf(proto.ErrTooManyContent, util.LimitConfig.ImageSize()),
			})
			return
		}

		// 2.检查文件类型
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" && ext != ".svg" {
			log.Errorf("Unsupported image format: %s", ext)
			c.JSON(http.StatusBadRequest, proto.PostPasteResp{
				Code:    http.StatusBadRequest,
				Message: proto.ErrUnspportType,
			})
			return
		}

		// 3.读取文件
		fileContent, err := file.Open()
		if err != nil {
			log.Errorf("Failed to open file: %+v", err)
			c.JSON(http.StatusInternalServerError, proto.PostPasteResp{
				Code:    http.StatusInternalServerError,
				Message: proto.ErrInvalidArgs,
			})
			return
		}
		defer fileContent.Close()

		// 4.获取图片尺寸
		width, height, err := util.GetImageDimensions(fileContent, ext)
		if err != nil {
			log.Errorf("Failed to get image dimensions: %+v", err)
			c.JSON(http.StatusBadRequest, proto.PostPasteResp{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf(proto.ErrTooManyContent, util.LimitConfig.ImageSize()),
			})
			return
		}

		// 重新打开文件流用于上传
		fileContent.Seek(0, 0)

		// 5.上传到对象存储
		objectKey, url, err = util.UploadToOSSWithExpiration(fileContent, file.Filename, duration)

		if err != nil {
			log.Errorf("Failed to upload image: %+v", err)
			c.JSON(http.StatusInternalServerError, proto.PostPasteResp{
				Code:    http.StatusInternalServerError,
				Message: proto.ErrInvalidArgs,
			})
			return
		}

		// 6.记录图片信息
		images = append(images, proto.ImageInfo{
			FileName:    file.Filename,
			ObjectKey:   objectKey,
			MimeType:    "image/" + ext,
			Width:       width,
			Height:      height,
			SizeBytes:   file.Size,
			URL:         url,
			UploadedAt:  time.Now(),
			AccessCount: 0,
		})
	}

	entry := db.PasteEntry{
		Title:       req.Title,
		Description: req.Description,
		Images:      images,
		ClientIP:    c.ClientIP(),
		CreatedAt:   time.Now(),
		ExpireAt:    time.Now().Add(duration),
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

// PostImageOnce 创建一次性图片分享内容
func (p *Paste) PostImageOnce(c *gin.Context) {
	var (
		ctx, log  = util.EnsureWithLogger(c)
		req       proto.PostPasteReq
		images    = []proto.ImageInfo{}
		objectKey string
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

	// 处理图片文件上传
	form, err := c.MultipartForm()
	if err != nil {
		log.Errorf("Failed to get multipart form: %+v", err)
		c.JSON(http.StatusBadRequest, proto.PostPasteResp{
			Code:    http.StatusBadRequest,
			Message: proto.ErrInvalidArgs,
		})
		return
	}

	files := form.File["images"]

	// 检查图片数量
	if len(files) == 0 {
		log.Error("No images uploaded")
		c.JSON(http.StatusBadRequest, proto.PostPasteResp{
			Code:    http.StatusBadRequest,
			Message: proto.ErrInvalidArgs,
		})
		return
	}

	if len(files) > util.LimitConfig.ImageCount() {
		log.Errorf("Image count exceeds the allowed limit: %d", len(files))
		c.JSON(http.StatusBadRequest, proto.PostPasteResp{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf(proto.ErrTooManyCount, util.LimitConfig.ImageCount()),
		})
		return
	}

	for _, file := range files {
		// 1.检查文件大小
		if file.Size > util.LimitConfig.ImageSize() {
			log.Errorf("Image size exceeds the allowed limit: %d bytes", file.Size)
			c.JSON(http.StatusBadRequest, proto.PostPasteResp{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf(proto.ErrTooManyContent, util.LimitConfig.ImageSize()),
			})
			return
		}

		// 2.检查文件类型
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" && ext != ".svg" {
			log.Errorf("Unsupported image format: %s", ext)
			c.JSON(http.StatusBadRequest, proto.PostPasteResp{
				Code:    http.StatusBadRequest,
				Message: proto.ErrUnspportType,
			})
			return
		}

		// 3.读取文件
		fileContent, err := file.Open()
		if err != nil {
			log.Errorf("Failed to open file: %+v", err)
			c.JSON(http.StatusInternalServerError, proto.PostPasteResp{
				Code:    http.StatusInternalServerError,
				Message: proto.ErrInvalidArgs,
			})
			return
		}
		defer fileContent.Close()

		// 4.获取图片尺寸
		width, height, err := util.GetImageDimensions(fileContent, ext)
		if err != nil {
			log.Errorf("Failed to get image dimensions: %+v", err)
			c.JSON(http.StatusBadRequest, proto.PostPasteResp{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf(proto.ErrTooManyContent, util.LimitConfig.ImageSize()),
			})
			return
		}

		// 重新打开文件流用于上传
		fileContent.Seek(0, 0)

		// 上传到对象存储，使用365天的过期时间（实际访问时会立即删除）
		objectKey, _, err = util.UploadToOSSWithExpiration(fileContent, file.Filename, time.Hour*24*365)
		if err != nil {
			log.Errorf("Failed to upload image: %+v", err)
			c.JSON(http.StatusInternalServerError, proto.PostPasteResp{
				Code:    http.StatusInternalServerError,
				Message: proto.ErrInvalidArgs,
			})
			return
		}

		// 6.记录图片信息，不生成URL（将在访问时生成）
		images = append(images, proto.ImageInfo{
			FileName:    file.Filename,
			ObjectKey:   objectKey,
			MimeType:    "image/" + ext,
			Width:       width,
			Height:      height,
			SizeBytes:   file.Size,
			UploadedAt:  time.Now(),
			AccessCount: 0,
		})
	}

	entry := db.PasteEntry{
		Title:       req.Title,
		Description: req.Description,
		Images:      images,
		ClientIP:    c.ClientIP(),
		CreatedAt:   time.Now(),
		Once:        true,
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

// GetImage 获取图片分享内容
func (p *Paste) GetImage(c *gin.Context) {
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

	// 获取图片配置
	if len(entry.Images) == 0 {
		c.JSON(http.StatusOK, proto.GetPasteResp{
			Code:    http.StatusBadRequest,
			Message: proto.ErrInvalidArgs,
		})
		return
	}

	// 如果是一次性访问的图片，生成短期URL并删除对象
	if entry.Once {
		for i := range entry.Images {
			// 生成5分钟有效的URL
			url, err := util.DefaultOSSClient.(*util.TencentOSSClient).GetSignedURL(entry.Images[i].ObjectKey, time.Minute*5)
			if err != nil {
				log.Errorf("Failed to generate signed URL: %+v", err)
				continue
			}
			entry.Images[i].URL = url
			entry.Images[i].AccessCount++

			// 删除对象
			go func(objectKey string) {
				time.Sleep(time.Minute)
				if err := util.DefaultOSSClient.(*util.TencentOSSClient).DeleteImage(objectKey); err != nil {
					log.Errorf("Failed to delete object: %+v", err)
				}
			}(entry.Images[i].ObjectKey)
		}
	}

	for i := range entry.Images {
		entry.Images[i].AccessCount++
	}
	c.JSON(http.StatusOK, proto.GetPasteResp{
		Code:        http.StatusOK,
		Images:      entry.Images,
		Title:       entry.Title,
		Description: entry.Description,
	})
}
