package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/viper"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

// TencentOSS 腾讯云对象存储
type TencentOSS struct {
	OSS    *cos.Client // 腾讯云对象存储客户端
	Config OSSConfig   // 腾讯云配置
}

// OSSConfig 腾讯云存储配置
type OSSConfig struct {
	Provider    string `mapstructure:"provider"`
	Region      string `mapstructure:"region"`
	Bucket      string `mapstructure:"bucket"`
	SecretID    string `mapstructure:"secret_id"`
	SecretKey   string `mapstructure:"secret_key"`
	URLExpireAt int    `mapstructure:"url_expire_at"`
}

func NewTencentOSS() (*TencentOSS, error) {
	newViper := viper.Sub("storage.cloud")
	if newViper == nil {
		return nil, errors.New("未找到 'storage.cloud' 配置")
	}

	// 解析腾讯云配置
	var config OSSConfig
	err := newViper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	// 构建腾讯云COS的存储桶URL
	// 格式: https://{bucket}.cos.{region}.myqcloud.com
	url, err := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", config.Bucket, config.Region))
	if err != nil {
		return nil, err
	}

	// 创建腾讯云COS客户端
	oss := cos.NewClient(
		&cos.BaseURL{
			BucketURL: url,
		},
		&http.Client{
			Timeout: 100 * time.Second,
			Transport: &cos.AuthorizationTransport{
				SecretID:  config.SecretID,
				SecretKey: config.SecretKey,
				Transport: &debug.DebugRequestTransport{
					RequestHeader:  true,
					RequestBody:    true,
					ResponseHeader: true,
					ResponseBody:   true,
				},
			},
		})

	return &TencentOSS{
		OSS:    oss,
		Config: config,
	}, nil
}

func (t *TencentOSS) Upload(ctx context.Context, content io.Reader, opts UploadOptions) error {
	// 准备自定义元数据
	header := &http.Header{}
	if opts.Metadata != nil {
		for k, v := range opts.Metadata {
			header.Add(fmt.Sprintf("x-cos-meta-%s", k), v)
		}
	}

	// 上传文件到私有存储桶
	// 默认权限即为私有读写，无需显式设置 ACL
	_, err := t.OSS.Object.Put(ctx, opts.ObjectKey, content, &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: opts.ContentType,
			XCosMetaXXX: header, // 使用转换后的 Header
		},
	})
	if err != nil {
		return fmt.Errorf("腾讯云COS上传失败: %w", err)
	}

	return nil
}

func (t *TencentOSS) SetLifeCycle(ctx context.Context) error {
	lc := &cos.BucketPutLifecycleOptions{
		Rules: []cos.BucketLifecycleRule{
			// 1小时后过期
			{
				ID:     "1h_rule",
				Filter: &cos.BucketLifecycleFilter{Prefix: ExpireAtPrefixMap[ExpireAt1Hour]},
				Status: "Enabled",
				Expiration: &cos.BucketLifecycleExpiration{
					Days: 1,
				},
			},
			// 1天后过期
			{
				ID:     "1d_rule",
				Filter: &cos.BucketLifecycleFilter{Prefix: ExpireAtPrefixMap[ExpireAt1Day]},
				Status: "Enabled",
				Expiration: &cos.BucketLifecycleExpiration{
					Days: 8,
				},
			},
			// 1个月后过期
			{
				ID:     "1m_rule",
				Filter: &cos.BucketLifecycleFilter{Prefix: ExpireAtPrefixMap[ExpireAt1Month]},
				Status: "Enabled",
				Expiration: &cos.BucketLifecycleExpiration{
					Days: 37,
				},
			},
			// 1年后过期
			{
				ID:     "1y_rule",
				Filter: &cos.BucketLifecycleFilter{Prefix: ExpireAtPrefixMap[ExpireAt1Year]},
				Status: "Enabled",
				Expiration: &cos.BucketLifecycleExpiration{
					Days: 372,
				},
			},
		},
	}
	_, err := t.OSS.Bucket.PutLifecycle(ctx, lc)
	if err != nil {
		return err
	}
	return nil
}

// GetSignedURL 生成带签名的临时访问 URL
func (t *TencentOSS) GetSignedURL(ctx context.Context, objectKey string) (string, error) {
	// 获取预签名 URL
	presignedURL, err := t.OSS.Object.GetPresignedURL(
		ctx,
		http.MethodGet,     // 使用 GET 方法访问
		objectKey,          // 对象键
		t.Config.SecretID,  // 使用配置中的 SecretID
		t.Config.SecretKey, // 使用配置中的 SecretKey
		time.Duration(t.Config.URLExpireAt)*time.Minute, // URL 有效期
		nil, // 不需要额外请求头
	)
	if err != nil {
		return "", fmt.Errorf("生成腾讯云COS预签名URL失败: %w", err)
	}

	return presignedURL.String(), nil
}
