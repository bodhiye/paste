package util

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"

	"github.com/srwiley/oksvg"
	"golang.org/x/image/webp"

	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tencentyun/cos-go-sdk-v5"
)

// OSSClient 对象存储接口
type OSSClient interface {
	UploadImage(fileContent multipart.File, fileName string) (string, string, error)
	DownloadImage(fileName string) ([]byte, error)
	DeleteImage(objectKey string) error
	GetSignedURL(objectKey string, expiration time.Duration) (string, error)
}

// LifecycleRule 定义生命周期规则
type LifecycleRule struct {
	ID         string // 规则ID
	Prefix     string // 应用规则的对象前缀
	Status     string // 规则状态: Enabled 或 Disabled
	ExpireDays int    // 过期天数，从对象最后修改时间开始计算
}

// TencentOSSClient 腾讯云对象存储客户端
type TencentOSSClient struct {
	SecretID  string      `mapstructure:"secret_id"`
	SecretKey string      `mapstructure:"secret_key"`
	Region    string      `mapstructure:"region"`
	Bucket    string      `mapstructure:"bucket"`
	BaseURL   string      `mapstructure:"base_url"`
	client    *cos.Client // 腾讯云SDK客户端实例
}

var (
	// DefaultOSSClient 默认OSS客户端实例
	DefaultOSSClient OSSClient
)

// InitCOSLifecycle 初始化COS生命周期规则
func (t *TencentOSSClient) InitCOSLifecycle(days int) error {
	// 创建生命周期规则
	rules := []cos.BucketLifecycleRule{
		{
			ID:     "Expires1Hour",
			Status: "Enabled",
			Filter: &cos.BucketLifecycleFilter{
				Prefix: "expires/1h/", // 1小时过期目录
			},
			Expiration: &cos.BucketLifecycleExpiration{
				Days: 1, // 1天后删除对象（额外保留时间）
			},
		},
		{
			ID:     "Expires1Day",
			Status: "Enabled",
			Filter: &cos.BucketLifecycleFilter{
				Prefix: "expires/1d/", // 1天过期目录
			},
			Expiration: &cos.BucketLifecycleExpiration{
				Days: 2, // 2天后删除对象（额外保留时间）
			},
		},
		{
			ID:     "Expires1Week",
			Status: "Enabled",
			Filter: &cos.BucketLifecycleFilter{
				Prefix: "expires/1w/", // 1周过期目录
			},
			Expiration: &cos.BucketLifecycleExpiration{
				Days: 9, // 9天后删除对象（额外保留时间）
			},
		},
		{
			ID:     "Expires1Month",
			Status: "Enabled",
			Filter: &cos.BucketLifecycleFilter{
				Prefix: "expires/1m/", // 1个月过期目录
			},
			Expiration: &cos.BucketLifecycleExpiration{
				Days: 37, // 37天后删除对象（额外保留时间）
			},
		},
		{
			ID:     "Expires1Year",
			Status: "Enabled",
			Filter: &cos.BucketLifecycleFilter{
				Prefix: "expires/1y/", // 1年过期目录
			},
			Expiration: &cos.BucketLifecycleExpiration{
				Days: 372, // 372天后删除对象（额外保留时间）
			},
		},
	}

	// 设置生命周期规则
	_, err := t.client.Bucket.PutLifecycle(context.Background(), &cos.BucketPutLifecycleOptions{
		Rules: rules,
	})

	if err != nil {
		log.Printf("Failed to set COS lifecycle rules: %v", err)
		return err
	}

	log.Printf("Successfully set COS lifecycle rules, %d rules in total", len(rules))
	return nil
}

// InitOSSLifecycle 初始化OSS生命周期规则（全局函数）
func InitOSSLifecycle() error {
	if DefaultOSSClient == nil {
		return fmt.Errorf("OSS client not initialized")
	}

	// 从配置获取生命周期规则设置
	expireDays := viper.GetInt("oss.lifecycle.expire_days")
	if expireDays <= 0 {
		expireDays = 30 // 默认30天过期
	}

	// 调用具体实现
	client, ok := DefaultOSSClient.(*TencentOSSClient)
	if !ok {
		return fmt.Errorf("current OSS client is not a TencentOSSClient type")
	}

	return client.InitCOSLifecycle(expireDays)
}

// InitTencentOSS 初始化腾讯云OSS客户端
func InitTencentOSS() {
	// 从配置中读取腾讯云配置
	tencentConfig := viper.Sub("oss.tencent")
	if tencentConfig == nil {
		log.Error("Tencent OSS configuration does not exist")
		return
	}

	client := &TencentOSSClient{}
	err := tencentConfig.Unmarshal(client)
	if err != nil {
		log.Errorf("Failed to parse Tencent OSS configuration: %v", err)
		return
	}

	// 检查必要配置
	if client.SecretID == "" || client.SecretKey == "" || client.Bucket == "" || client.Region == "" {
		log.Error("Tencent OSS configuration is incomplete")
		return
	}

	// 创建腾讯云客户端
	bucketURL, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", client.Bucket, client.Region))
	serviceURL, _ := url.Parse(fmt.Sprintf("https://cos.%s.myqcloud.com", client.Region))

	b := &cos.BaseURL{
		BucketURL:  bucketURL,
		ServiceURL: serviceURL,
	}

	cosClient := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  client.SecretID,
			SecretKey: client.SecretKey,
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		},
	})

	client.client = cosClient
	DefaultOSSClient = client
	log.Info("Tencent OSS client initialized successfully")

	// 如果启用了生命周期管理，初始化生命周期规则
	if viper.GetBool("oss.lifecycle.enabled") {
		if err := InitOSSLifecycle(); err != nil {
			log.Warnf("Failed to initialize COS lifecycle rules: %v", err)
		} else {
			log.Info("Tencent COS lifecycle rules initialized successfully")
		}
	}
}

// 初始化OSS客户端
func InitOSS() {
	InitTencentOSS()
}

// UploadImageWithExpiration 上传图片到腾讯云COS，支持设置过期时间
func (t *TencentOSSClient) uploadImageWithExpiration(fileContent multipart.File, fileName string, duration time.Duration) (string, string, error) {
	// 根据过期时间决定存储路径
	var dirPrefix string
	if duration <= time.Hour {
		// 1小时内过期
		dirPrefix = "expires/1h/"
	} else if duration <= 24*time.Hour {
		// 1天内过期
		dirPrefix = "expires/1d/"
	} else if duration <= 7*24*time.Hour {
		// 1周内过期
		dirPrefix = "expires/1w/"
	} else if duration <= 30*24*time.Hour {
		// 1个月内过期
		dirPrefix = "expires/1m/"
	} else {
		// 1年内过期
		dirPrefix = "expires/1y/"
	}

	// 生成唯一文件名
	fileExt := filepath.Ext(fileName)
	objectKey := fmt.Sprintf("%s%d_%s%s", dirPrefix, time.Now().UnixNano(), RandString(8), fileExt)
	url, err := t.GetSignedURL(objectKey, duration)
	if err != nil {
		log.Errorf("Invalid expiration format: %v", err)
		return objectKey, "", err
	}

	// 获取文件的Content-Type
	contentType := getContentTypeByExt(fileExt)

	_, err = t.client.Object.Put(
		context.Background(),
		objectKey,
		fileContent,
		&cos.ObjectPutOptions{
			ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
				ContentType: contentType,
				// 添加额外的头部，确保文件可以被正确预览
				CacheControl: "max-age=31536000",
			},
		},
	)

	if err != nil {
		log.Printf("Failed to upload image to Tencent COS: %v", err)
		return objectKey, "", err
	}
	return objectKey, url, nil
}

// 根据文件扩展名获取Content-Type
func getContentTypeByExt(ext string) string {
	ext = strings.ToLower(ext)
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	default:
		return "application/octet-stream"
	}
}

// UploadToOSSWithExpiration 上传图片到OSS并返回URL，支持设置过期时间
func UploadToOSSWithExpiration(fileContent multipart.File, fileName string, duration time.Duration) (string, string, error) {
	if DefaultOSSClient == nil {
		return "", "", fmt.Errorf("OSS client not initialized")
	}

	client, ok := DefaultOSSClient.(*TencentOSSClient)
	if !ok {
		return "", "", fmt.Errorf("current OSS client is not a TencentOSSClient type")
	}

	return client.uploadImageWithExpiration(fileContent, fileName, duration)
}

// GetImageDimensions 获取图片的宽度、高度和格式
func GetImageDimensions(file io.Reader, ext string) (int, int, error) {
	var img image.Image
	var err error

	// 根据扩展名选择解码器
	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	case ".gif":
		img, err = gif.Decode(file)
	case ".webp":
		img, err = webp.Decode(file)
	case ".svg":
		data, err := io.ReadAll(file)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to read file: %v", err)
		}
		icon, err := oksvg.ReadIconStream(bytes.NewReader(data))
		if err != nil {
			return 0, 0, fmt.Errorf("failed to parse SVG image: %v", err)
		}
		width := int(icon.ViewBox.W)
		height := int(icon.ViewBox.H)
		return width, height, nil
	default:
		return 0, 0, fmt.Errorf("unsupported image format: %s", ext)
	}

	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse image: %v", err)
	}

	bounds := img.Bounds()
	return bounds.Dx(), bounds.Dy(), nil
}

// UploadImage 上传图片到腾讯云OSS
func (t *TencentOSSClient) UploadImage(fileContent multipart.File, fileName string) (string, string, error) {
	// 调用带过期时间的方法，默认不设置过期时间（0表示永不过期）
	return t.uploadImageWithExpiration(fileContent, fileName, 0)
}

// DownloadImage 从腾讯云OSS下载图片
func (t *TencentOSSClient) DownloadImage(objectKey string) ([]byte, error) {
	// 获取对象
	resp, err := t.client.Object.Get(context.Background(), objectKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to download file from Tencent COS: %v", err)
	}
	defer resp.Body.Close()

	// 读取内容
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %v", err)
	}

	return data, nil
}

// GetSignedURL 生成临时签名URL
func (t *TencentOSSClient) GetSignedURL(objectKey string, expiration time.Duration) (string, error) {
	// 如果objectKey是完整URL，提取实际的对象键
	if strings.HasPrefix(objectKey, "http") {
		// 解析URL
		parsedURL, err := url.Parse(objectKey)
		if err != nil {
			return "", fmt.Errorf("failed to parse object URL: %v", err)
		}
		// 提取路径部分作为对象键
		objectKey = strings.TrimPrefix(parsedURL.Path, "/") // 移除前缀/
	}

	// 创建签名URL选项
	ctx := context.Background()
	opt := &cos.PresignedURLOptions{
		Query:  &url.Values{},
		Header: &http.Header{},
	}

	// 添加必要的请求头
	opt.Header.Set("Content-Type", "application/octet-stream")
	opt.Header.Set("Content-Disposition", "inline")
	opt.Header.Set("Cache-Control", "no-cache")

	// 添加必要的查询参数
	opt.Query.Set("response-content-type", "application/octet-stream")
	opt.Query.Set("response-cache-control", "no-cache")

	// 生成预签名URL
	signedURL, err := t.client.Object.GetPresignedURL(
		ctx, http.MethodGet, objectKey,
		t.SecretID, t.SecretKey,
		expiration, opt,
	)

	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %v", err)
	}

	return signedURL.String(), nil
}

// DeleteImage 从腾讯云COS删除对象
func (t *TencentOSSClient) DeleteImage(objectKey string) error {
	_, err := t.client.Object.Delete(context.Background(), objectKey)
	if err != nil {	
		return fmt.Errorf("failed to delete object from COS: %v", err)
	}
	log.Infof("Deleted object from COS: %s", objectKey)
	return nil
}
