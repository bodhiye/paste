package storage

import (
	"context"
	"fmt"
	"io"
)

const (
	ProviderTencent = "tencent"
	ProviderAliyun  = "aliyun"
	ProviderQiniu   = "qiniu"
	ProviderAWS     = "aws"
	ProviderAzure   = "azure"
)

type UploadOptions struct {
	ObjectKey   string            // 对象存储中的唯一标识符
	ContentType string            // 文件的内容类型
	Metadata    map[string]string // 文件的元数据信息
}

type OSS interface {
	Upload(ctx context.Context, content io.Reader, opts UploadOptions) error
	SetLifeCycle(ctx context.Context) error
	GetSignedURL(ctx context.Context, objectKey string) (string, error)
}

func NewOSSWithFactory(provider string) (OSS, error) {
	switch provider {
	case ProviderTencent:
		oss, err := NewTencentOSS()
		if err != nil {
			return nil, err
		}
		err = oss.SetLifeCycle(context.Background())
		if err != nil {
			return nil, err
		}
		return oss, nil
	default:
		return nil, fmt.Errorf("不支持的云存储提供商: %s", provider)
	}
}
