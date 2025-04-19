package util

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"paste.org.cn/paste/server/proto"
)

// 保护图片操作的互斥锁
var imageMutex sync.Mutex

// DeleteImage 删除单个图片文件
// image 参数是图片的URL路径，例如: "/uploads/123456_abcd.jpg"
func DeleteImage(image string)error  {
	return DeleteFiles([]string{image})
}

// DeleteFiles 删除一个或多个图片文件
func DeleteFiles(files []string) (err error) {
	if len(files) == 0 {
		return nil
	}

	// 使用互斥锁保护文件操作
	imageMutex.Lock()
	defer imageMutex.Unlock()

	for _, file := range files {
		// 跳过空文件路径
		if file == "" {
			continue
		}

		filename := filepath.Base(file)
		filePath := GetImagePath(filename)

		// 删除文件
		if err1 := os.Remove(filePath); err1 != nil && !os.IsNotExist(err1) {
			// 记录错误但继续尝试删除其他文件
			err = err1
		}
	}
	return err
}

// 获取图片，使用简洁的实现保证原子性
func UploadImages(c *gin.Context, log *logrus.Entry) (images []proto.ImageFile, err error) {
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
	if len(files) > LimitConfig.ImagesCount() {
		log.Errorf("图片数量过多: %d", len(files))
		return nil, fmt.Errorf(proto.ErrTooManyCount, LimitConfig.ImagesCount())
	}

	// 确保上传目录存在
	uploadDir := GetUploadDir()

	// 存储已上传的文件路径，用于出错时清理
	var uploadedFiles []string

	// 处理每个图片
	for _, fileHeader := range files {
		// 检查文件大小
		fileSizeMB := fileHeader.Size / (1024 * 1024)
		if fileSizeMB > int64(LimitConfig.ImagesSize()) {
			log.Errorf("图片太大: %d MB", fileSizeMB)
			if err := DeleteFiles(uploadedFiles); err != nil {
				log.Errorf("删除图片失败: %v", err)
			}
			return nil, fmt.Errorf(proto.ErrOverMaxSize, LimitConfig.ImagesSize())
		}

		// 检查 MIME 类型
		contentType := fileHeader.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "image/") {
			log.Errorf("不支持的文件类型: %s", contentType)
			if err := DeleteFiles(uploadedFiles); err != nil {
				log.Errorf("删除图片失败: %v", err)
			}
			return nil, errors.New(proto.ErrInvalidFileType)
		}

		// 生成唯一文件名
		fileExt := filepath.Ext(fileHeader.Filename)
		newFilename := fmt.Sprintf("%d_%s%s",
			time.Now().UnixNano(),
			uuid.NewString(), // 使用完整UUID以确保唯一性
			fileExt)

		// 文件完整路径
		filePath := filepath.Join(uploadDir, newFilename)

		// 使用互斥锁保护文件保存操作
		imageMutex.Lock()
		err := c.SaveUploadedFile(fileHeader, filePath)
		imageMutex.Unlock()

		if err != nil {
			log.Errorf("保存文件失败: %+v", err)
			if err := DeleteFiles(uploadedFiles); err != nil {
				log.Errorf("删除图片失败: %v", err)
			}
			return nil, err
		}

		// 记录已上传文件路径
		uploadedFiles = append(uploadedFiles, filePath)

		// 添加到返回结果
		images = append(images, proto.ImageFile{
			Filename:    fileHeader.Filename,
			URL:         GetImageURL(newFilename),
			Size:        fileHeader.Size,
			ContentType: contentType,
		})
	}

	return images, nil
}
