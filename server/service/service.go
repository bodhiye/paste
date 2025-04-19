package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"paste.org.cn/paste/server/db"
	"paste.org.cn/paste/server/proto"
	"paste.org.cn/paste/server/util"
)

type Paste struct {
	db.Paste
}

type tempImage struct {
	TempPath  string
	FinalPath string
	ImageFile proto.ImageFile
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

	req.Images, err = uploadImages(c, log)
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
	req.Images, err = uploadImages(c, log)
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

// 获取图片
func uploadImages(c *gin.Context, log *logrus.Entry) (images []proto.ImageFile, err error) {
	// 单独处理文件上传
	form, err := c.MultipartForm()
	if err != nil || form == nil {
		return nil, nil // 没有文件上传，不是错误
	}

	files := form.File["images"]
	if len(files) == 0 {
		return nil, nil // 没有图片，不是错误
	}

	// 检查图片数量限制
	if len(files) > util.LimitConfig.ImagesCount() {
		log.Errorf("图片数量过多: %d", len(files))
		return nil, fmt.Errorf(proto.ErrTooManyCount, util.LimitConfig.ImagesCount())
	}

	// 确保上传目录存在 - 使用storage中的函数确保线程安全
	uploadDir := util.GetUploadDir()

	// 临时图片列表，用于跟踪暂存的图片
	var tempImages []tempImage

	// 处理每个图片
	for _, fileHeader := range files {
		// 检查文件大小
		fileSizeMB := fileHeader.Size / (1024 * 1024)
		if fileSizeMB > int64(util.LimitConfig.ImagesSize()) {
			log.Errorf("图片太大: %d MB", fileSizeMB)
			// 清理已上传的临时文件
			cleanupTempImages(tempImages)
			return nil, fmt.Errorf(proto.ErrOverMaxSize, util.LimitConfig.ImagesSize())
		}

		// 检查 MIME 类型
		contentType := fileHeader.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "image/") {
			log.Errorf("不支持的文件类型: %s", contentType)
			// 清理已上传的临时文件
			cleanupTempImages(tempImages)
			return nil, errors.New(proto.ErrInvalidFileType)
		}

		// 生成一个更独特的唯一文件名，减少冲突可能性
		fileExt := filepath.Ext(fileHeader.Filename)
		newFilename := fmt.Sprintf("%d_%s%s",
			time.Now().UnixNano(),
			uuid.NewString(), // 使用完整的UUID
			fileExt)

		// 生成临时文件路径和最终文件路径
		tempFilename := "temp_" + newFilename
		tempPath := filepath.Join(uploadDir, tempFilename)
		finalPath := filepath.Join(uploadDir, newFilename)

		// 保存到临时文件
		if err := c.SaveUploadedFile(fileHeader, tempPath); err != nil {
			log.Errorf("保存文件失败: %+v", err)
			// 清理已上传的临时文件
			cleanupTempImages(tempImages)
			return nil, err
		}

		// 创建图片记录
		imageFile := proto.ImageFile{
			Filename:    fileHeader.Filename,
			URL:         util.GetImageURL(newFilename),
			Size:        fileHeader.Size,
			ContentType: contentType,
		}

		// 添加到临时列表
		tempImages = append(tempImages, tempImage{
			TempPath:  tempPath,
			FinalPath: finalPath,
			ImageFile: imageFile,
		})
	}

	// 创建同步等待组和错误通道，以便并发移动文件
	var wg sync.WaitGroup
	errorCh := make(chan error, len(tempImages))
	imagesMutex := sync.Mutex{}

	// 所有图片都已成功上传到临时位置，现在并发移动到最终位置
	for _, img := range tempImages {
		wg.Add(1)
		go func(img tempImage) {
			defer wg.Done()

			// 移动文件
			if err := os.Rename(img.TempPath, img.FinalPath); err != nil {
				log.Errorf("移动临时文件失败: %+v", err)
				errorCh <- err
				return
			}

			// 线程安全地添加到最终结果
			imagesMutex.Lock()
			images = append(images, img.ImageFile)
			imagesMutex.Unlock()
		}(img)
	}

	// 等待所有移动操作完成
	wg.Wait()
	close(errorCh)

	// 检查是否有错误
	select {
	case err := <-errorCh:
		// 有错误发生，清理所有文件
		for _, img := range tempImages {
			os.Remove(img.TempPath)  // 删除可能仍存在的临时文件
			os.Remove(img.FinalPath) // 删除可能已移动的文件
		}
		return nil, err
	default:
		// 没有错误，继续处理
	}

	return images, nil
}

// 辅助函数，清理临时图片文件
func cleanupTempImages(tempImages []tempImage) {
	for _, img := range tempImages {
		os.Remove(img.TempPath)
	}
}