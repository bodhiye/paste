package util

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"paste.org.cn/paste/server/proto"
)

// ImageCleaner 图片清理器
type ImageCleaner struct {
	db       *mongo.Collection // MongoDB集合
	interval time.Duration     // 清理间隔
	stopChan chan struct{}     // 停止信号
	mutex    sync.Mutex        // 清理操作的互斥锁
}

// NewImageCleaner 创建一个新的图片清理器
func NewImageCleaner(db *mongo.Collection) *ImageCleaner {
	return &ImageCleaner{
		db:       db,
		interval: time.Duration(viper.GetInt("cleaner.interval"))*time.Minute,
		stopChan: make(chan struct{}),
	}
}

// Start 启动清理器
func (ic *ImageCleaner) Start() {
	go func() {
		ticker := time.NewTicker(ic.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := ic.Clean(); err != nil {
					log.Errorf("清理图片失败: %v", err)
				}
			case <-ic.stopChan:
				return
			}
		}
	}()
}

// Stop 停止清理器
func (ic *ImageCleaner) Stop() {
	close(ic.stopChan)
}

// Clean 执行清理操作
func (ic *ImageCleaner) Clean() error {
	// 获取锁以确保只有一个清理操作在进行
	ic.mutex.Lock()
	defer ic.mutex.Unlock()

	// 获取所有图片文件
	files, err := os.ReadDir(GetUploadDir())
	if err != nil {
		return err
	}

	// 创建文件名到文件的映射
	fileMap := make(map[string]bool)
	for _, file := range files {
		if !file.IsDir() {
			fileMap[file.Name()] = true
		}
	}

	// 从数据库获取所有有效的图片URL
	ctx := context.Background()
	cursor, err := ic.db.Find(ctx, bson.M{
		"images": bson.M{"$exists": true, "$ne": []proto.ImageFile{}},
		"$or": []bson.M{
			{"expire_at": bson.M{"$exists": false}},  // 没有过期时间的
			{"expire_at": bson.M{"$gt": time.Now()}}, // 未过期的
		},
	})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	// 遍历数据库中的所有图片记录
	validFiles := make(map[string]bool)
	var entry struct {
		Images []proto.ImageFile `bson:"images"`
	}
	for cursor.Next(ctx) {
		if err := cursor.Decode(&entry); err != nil {
			log.Errorf("解析数据库记录失败: %v", err)
			continue
		}
		for _, img := range entry.Images {
			filename := filepath.Base(img.URL)
			validFiles[filename] = true
		}
	}

	// 删除没有对应数据库记录的文件
	for filename := range fileMap {
		if !validFiles[filename] {
			// 使用DeleteImage函数而不是直接删除文件
			if err := DeleteImage(GetImageURL(filename)); err != nil {
				log.Errorf("删除孤儿图片失败: %v", err)
			}
			log.Infof("已删除孤儿图片: %s", GetImagePath(filename))
		}
	}

	return nil
}
