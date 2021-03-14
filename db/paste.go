package db

import (
	"context"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"paste.org.cn/paste/util"
)

type Paste interface {
	Set(ctx context.Context, entry PasteEntry) (string, error)
	Get(ctx context.Context, key string) (PasteEntry, error)
}

type _Paste struct {
	*mongo.Collection
}

type MongoConfig struct {
	Host string
	DB   string
	Coll string
}

func NewPaste(ctx context.Context, viper_ *viper.Viper) (Paste, error) {
	var config MongoConfig
	if err := viper_.Unmarshal(&config); err != nil {
		return nil, err
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Host))
	if err != nil {
		return nil, err
	}

	paste := _Paste{
		Collection: client.Database(config.DB).Collection(config.Coll),
	}
	if err := paste.Init(ctx); err != nil {
		return nil, err
	}
	return paste, nil
}

func (p _Paste) Init(ctx context.Context) error {
	models := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "key", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "expire_at", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(0),
		},
	}

	opts := options.CreateIndexes().SetMaxTime(1 * time.Minute)
	_, err := p.Indexes().CreateMany(ctx, models, opts)
	return err
}

func (p _Paste) Set(ctx context.Context, entry PasteEntry) (key string, err error) {
	entry.CreatedAt = time.Now()
	entry.Key = util.RandString(10)

	for {
		_, err = p.Collection.InsertOne(ctx, entry)
		if mongo.IsDuplicateKeyError(err) {
			entry.Key = util.RandString(10)
			continue
		}
		break
	}

	key = entry.Key
	return
}

func (p _Paste) Get(ctx context.Context, key string) (entry PasteEntry, err error) {
	err = p.Collection.FindOne(ctx, bson.M{"key": key}).Decode(&entry)
	return
}
