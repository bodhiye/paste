package db

import (
	"context"
	"errors"
	"time"

	"github.com/spf13/viper"                    // 用于读取配置文件
	"go.mongodb.org/mongo-driver/bson"          // 用于处理 BSON 格式数据
	"go.mongodb.org/mongo-driver/mongo"         // MongoDB 驱动
	"go.mongodb.org/mongo-driver/mongo/options" // 用于设置 MongoDB 操作的选项
	"golang.org/x/crypto/bcrypt"

	"paste.org.cn/paste/server/proto"
	"paste.org.cn/paste/server/util"
)

// Paste 接口定义了与 Paste 数据相关的操作
type Paste interface {
	Set(ctx context.Context, entry PasteEntry) (string, error)
	Get(ctx context.Context, key, password string) (PasteEntry, error)
	Delete(ctx context.Context, key string) error
}

// _Paste 结构体是 Paste 接口的实现，内嵌了 mongo.Collection，用于操作 MongoDB 的集合
type _Paste struct {
	*mongo.Collection
}

// 存储 MogoDB 的连接配置
type MongoConfig struct {
	Host string `mapstructure:"host"` // 连接地址
	DB   string `mapstructure:"db"`   // 数据库名
	Coll string `mapstructure:"coll"` // 集合名
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

// Set 方法将新的 PasteEntry 存储到 MongoDB 中，并返回生成的唯一键
func (p _Paste) Set(ctx context.Context, entry PasteEntry) (key string, err error) {
	entry.Key = util.RandString(10) //生成长度为10的随机字符串

	for {
		// 尝试将 entry 插入到集合中
		_, err = p.Collection.InsertOne(ctx, entry)
		if mongo.IsDuplicateKeyError(err) {
			// 如果遇到键重复错误，重新生成键并继续尝试插入
			entry.Key = util.RandString(10)
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
	// 在集合中查找匹配的键，并将结果解码到 entry 中
	err = p.Collection.FindOne(ctx, bson.M{"key": key}).Decode(&entry)
	if err != nil {
		// 如果查找失败，返回错误
		return
	}

	// 如果 entry 设置了密码，验证提供的密码是否匹配
	if entry.Password != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(entry.Password), []byte(password)); err != nil {
			return entry, errors.New(proto.ErrWrongPassword) // 密码错误
		}
	}

	// 检查 entry 是否已过期
	if !entry.ExpireAt.IsZero() && time.Now().After(entry.ExpireAt) {
		err = errors.New(proto.ErrContentExpired) // 内容已过期
		return
	}

	// 如果 entry 设置为一次性查看，删除该 entry
	if entry.Once {
		err = p.Delete(ctx, key)
	}
	return
}

// Delete 方法根据键删除对应的 PasteEntry
func (p _Paste) Delete(ctx context.Context, key string) (err error) {
	// 在集合中删除匹配的键的文档
	_, err = p.Collection.DeleteOne(ctx, bson.M{"key": key})
	return
}
