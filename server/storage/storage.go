package storage

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"paste.org.cn/paste/server/proto"
	"paste.org.cn/paste/server/util"
)

// 存储类型常量
const (
	StorageTypeBase64 = "base64" // Base64存储在MongoDB
	StorageTypeCloud  = "cloud"  // 第三方云存储
)

// 过期时间常量
const (
	ExpireAt1Hour     = "1"
	ExpireAt1Day      = "24"
	ExpireAt1Month    = "720"
	ExpireAt1Year     = "8760"
	ExpireAtPermanent = "0"
)

// 过期时间前缀映射
var ExpireAtPrefixMap = map[string]string{
	ExpireAt1Hour:     "1hour/",
	ExpireAt1Day:      "1day/",
	ExpireAt1Month:    "1month/",
	ExpireAt1Year:     "1year/",
	ExpireAtPermanent: "permanent/",
}

// ImageStorageConfig 存储策略配置
type ImageStorageConfig struct {
	Type string // 存储类型: base64, cloud
	OSS  OSS    // 云存储客户端
}

// 默认存储配置
var StorageConfig ImageStorageConfig

// InitializeStorage 初始化存储配置
func InitializeStorage() {
	// 从配置文件中读取存储配置
	StorageConfig.Type = viper.GetString("storage.type")
	if StorageConfig.Type == StorageTypeCloud {
		// 读取具体的云服务提供商
		provider := viper.GetString("storage.cloud.provider") // 从 storage.cloud.provider 读取
		if provider == "" {
			log.Errorf("未配置云存储提供商 (storage.cloud.provider)，回退到 base64 存储")
			StorageConfig.Type = StorageTypeBase64
			return
		}

		oss, err := NewOSSWithFactory(provider)
		if err != nil {
			log.Errorf("初始化云存储客户端失败 (%s): %+v", provider, err)
			// 初始化失败时，回退到 base64 存储
			StorageConfig.Type = StorageTypeBase64
			return
		}
		StorageConfig.OSS = oss
		log.Infof("云存储 (%s) 初始化成功", provider)
	} else {
		StorageConfig.Type = StorageTypeBase64
		log.Info("使用 base64 存储")
	}
}

// 将文件转换为base64编码
func FileToBase64(file io.Reader) (string, error) {
	// 读取文件内容
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	// 转换为base64编码
	return base64.StdEncoding.EncodeToString(fileBytes), nil
}

// 获取图片
func UploadImages(c *gin.Context, log *log.Entry) (images []proto.ImageFile, err error) {
	// 单独处理文件上传
	form, err := c.MultipartForm()
	if err != nil {
		if err != http.ErrNotMultipart {
			log.Errorf("解析 multipart form 失败: %+v", err)
		}
		return nil, nil
	}

	files := form.File["images"]
	if len(files) == 0 {
		return nil, nil
	}

	// 检查图片数量限制
	if len(files) > util.LimitConfig.ImagesCount() {
		log.Errorf("图片数量过多: %d", len(files))
		return nil, fmt.Errorf(proto.ErrTooManyCount, util.LimitConfig.ImagesCount())
	}

	// 处理每个图片
	for _, fileHeader := range files {
		// 检查文件大小
		fileSizeMB := fileHeader.Size / (1024 * 1024)
		if fileSizeMB > int64(util.LimitConfig.ImagesSize()) {
			log.Errorf("图片太大: %d MB", fileSizeMB)
			return nil, fmt.Errorf(proto.ErrOverMaxSize, util.LimitConfig.ImagesSize())
		}

		// 检查 MIME 类型
		contentType := fileHeader.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "image/") {
			log.Errorf("不支持的文件类型: %s", contentType)
			return nil, errors.New(proto.ErrInvalidFileType)
		}

		// 生成唯一文件名
		fileExt := filepath.Ext(fileHeader.Filename)
		newFilename := fmt.Sprintf("%d_%s%s",
			time.Now().UnixNano(),
			uuid.NewString(), // 使用完整UUID以确保唯一性
			fileExt)

		// 创建图片对象
		imageFile := proto.ImageFile{
			StorageType: StorageConfig.Type,
			Filename:    newFilename,
			Size:        fileHeader.Size,
			ContentType: contentType,
		}

		// 根据存储类型处理
		file, err := fileHeader.Open()
		if err != nil {
			log.Errorf("打开文件 '%s' 失败: %+v", fileHeader.Filename, err)
			return nil, fmt.Errorf("无法处理文件: %s", fileHeader.Filename)
		}
		defer file.Close()

		switch StorageConfig.Type {
		case StorageTypeBase64:
			// Base64编码存储
			imageFile.Base64Content, err = FileToBase64(file)
			if err != nil {
				log.Errorf("Base64编码失败: %+v", err)
				return nil, err
			}
		case StorageTypeCloud:
			// 上传图片到云存储
			var err error
			imageFile.ObjectKey, err = StorageInCloud(c, imageFile, file)
			if err != nil {
				log.Errorf("上传图片到云存储失败: %+v", err)
				return nil, err
			}
		}
		images = append(images, imageFile)
	}

	return images, nil
}

// StorageInCloud 将文件上传到云存储并返回 ObjectKey
func StorageInCloud(c *gin.Context, imageFile proto.ImageFile, file io.Reader) (string, error) {
	context, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// 获取过期时间
	expiredAt := c.PostForm("expire_at")
	var prefix string
	switch expiredAt {
	case ExpireAt1Hour:
		prefix = "1hour/"
	case ExpireAt1Day:
		prefix = "1day/"
	case ExpireAt1Month:
		prefix = "1month/"
	case ExpireAt1Year:
		prefix = "1year/"
	case ExpireAtPermanent:
		prefix = "permanent/"
	default:
		prefix = "1hour/"
	}

	opt := UploadOptions{
		ObjectKey:   prefix + imageFile.Filename,
		ContentType: imageFile.ContentType,
	}
	// 调用 Upload 方法，返回 ObjectKey
	err := StorageConfig.OSS.Upload(context, file, opt)
	if err != nil {
		log.Errorf("上传图片到云存储失败: %+v", err)
		return "", err
	}
	return opt.ObjectKey, nil
}
