package db

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/spf13/viper"                    // 用于读取配置文件
	"go.mongodb.org/mongo-driver/bson"          // 用于处理 BSON 格式数据
	"go.mongodb.org/mongo-driver/mongo"         // MongoDB 驱动
	"go.mongodb.org/mongo-driver/mongo/options" // 用于设置 MongoDB 操作的选项
	"golang.org/x/crypto/bcrypt"

	"paste.org.cn/paste/server/proto"

	"github.com/google/uuid"
)

// 保存全局MongoDB客户端实例
var (
	mongoClient *mongo.Client
	clientMutex sync.RWMutex
)

// Paste 接口定义了与 Paste 数据相关的操作
type Paste interface {
	Set(ctx context.Context, entry PasteEntry) (string, error)
	Get(ctx context.Context, key, password string) (PasteEntry, error)
	GetCollection() *mongo.Collection
}

// _Paste 结构体是 Paste 接口的实现，内嵌了 mongo.Collection，用于操作 MongoDB 的集合
type _Paste struct {
	*mongo.Collection
}

// 存储 MogoDB 的连接配置
type MongoConfig struct {
	Host string // 连接地址
	DB   string // 数据库名
	Coll string // 集合名
}

// GetMongoClient 返回全局MongoDB客户端实例
func GetMongoClient() *mongo.Client {
	clientMutex.RLock()
	defer clientMutex.RUnlock()
	return mongoClient
}

func NewPaste(ctx context.Context, viper_ *viper.Viper) (Paste, error) {
	var config MongoConfig
	// 将提取出来的viper实例映射到config
	if err := viper_.Unmarshal(&config); err != nil {
		return nil, err
	}

	// 连接 mogoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Host))
	if err != nil {
		return nil, err
	}

	// 保存MongoDB客户端实例到全局变量
	clientMutex.Lock()
	mongoClient = client
	clientMutex.Unlock()

	// 创建 _Paste 实例，并传入 MongoDB 的 Collection
	paste := _Paste{
		Collection: client.Database(config.DB).Collection(config.Coll),
	}
	// 初始化 Paste 实例，例如创建索引
	if err := paste.Init(ctx); err != nil {
		return nil, err // 如果初始化失败，则返回错误
	}
	return paste, nil // 返回创建成功的 Paste 实例
}

// Init 方法用于初始化 MongoDB 的集合，例如创建索引
func (p _Paste) Init(ctx context.Context) error {
	// 定义需要创建的索引模型，为集合创建索引，以优化数据存取效率
	models := []mongo.IndexModel{ // keys 用于指定索引的字段及其排序顺序
		{
			// 创建 key 的唯一索引，保证 key 的唯一性
			Keys: bson.D{
				{Key: "key", Value: 1}, // 1 升序，-1降序
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "created_at", Value: -1}, // 按创建时间降序排列
			},
		},
		{
			// 创建 expire_at 的过期索引，当文档的 expire_at 字段到达指定时间后，MongoDB 会自动删除该文档
			Keys:    bson.D{{Key: "expire_at", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(0), // 设置过期时间为到达 expire_at 字段立即过期
		},
	}

	// 设置创建索引的选项，例如最大执行时间
	opts := options.CreateIndexes().SetMaxTime(1 * time.Minute)
	// 在集合上创建多个索引
	_, err := p.Indexes().CreateMany(ctx, models, opts)
	return err // 返回创建索引过程中发生的错误
}

// GetCollection 返回 MongoDB 的集合
func (p _Paste) GetCollection() *mongo.Collection {
	return p.Collection
}

// Set 方法将新的 PasteEntry 存储到 MongoDB 中，并返回生成的唯一键
func (p _Paste) Set(ctx context.Context, entry PasteEntry) (key string, err error) {
	// 生成一个更长的随机键，减少碰撞概率
	entry.Key = uuid.NewString()[:16] //生成长度为16的随机字符串

	for {
		// 尝试将 entry 插入到集合中
		_, err = p.Collection.InsertOne(ctx, entry)
		if mongo.IsDuplicateKeyError(err) {
			// 如果遇到键重复错误，重新生成键并继续尝试插入
			entry.Key = uuid.NewString()[:16]
			continue
		}
		// 如果没有发生重复键错误，跳出循环
		break
	}

	// 返回生成的键
	key = entry.Key
	return
}

// Get 方法根据提供的键和密码检索相应的 PasteEntry
func (p _Paste) Get(ctx context.Context, key, password string) (entry PasteEntry, err error) {
	// 创建一个查询条件
	filter := bson.M{"key": key}

	// 如果是一次性文档，使用 FindOneAndDelete 原子操作
	if isOnce, err := p.isOnceDocument(ctx, key); err == nil && isOnce {
		// 使用 FindOneAndDelete 原子地获取并删除文档
		err = p.Collection.FindOneAndDelete(ctx, filter).Decode(&entry)
		if err != nil {
			return entry, err
		}
	} else {
		// 对于非一次性文档，使用普通的 FindOne
		err = p.Collection.FindOne(ctx, filter).Decode(&entry)
		if err != nil {
			return entry, err
		}
	}

	// 如果 entry 设置了密码，验证提供的密码是否匹配
	if entry.Password != "" && bcrypt.CompareHashAndPassword([]byte(entry.Password), []byte(password)) != nil {
		err = errors.New(proto.ErrWrongPassword) // 密码错误
		return
	}

	// 检查 entry 是否已过期
	if !entry.ExpireAt.IsZero() && time.Now().After(entry.ExpireAt) {
		err = errors.New(proto.ErrContentExpired) // 内容已过期
		return
	}

	return
}

// isOnceDocument 检查文档是否是一次性文档
func (p _Paste) isOnceDocument(ctx context.Context, key string) (bool, error) {
	var result struct {
		Once bool `bson:"once"`
	}
	err := p.Collection.FindOne(ctx, bson.M{"key": key}).Decode(&result)
	if err != nil {
		return false, err
	}
	return result.Once, nil
}
