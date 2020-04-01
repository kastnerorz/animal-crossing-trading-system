package pkg

import (
	"context"
	"github.com/kastnerorz/animal-crossing-trading-system/backend/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func InitiateMongoClient() *mongo.Client {
	var err error
	var client *mongo.Client
	opts := options.Client()
	opts.ApplyURI(configs.MongoURI)
	opts.SetMaxPoolSize(5)
	if client, err = mongo.Connect(context.Background(), opts); err != nil {
		log.Println(err.Error())
	}
	return client
}

func GetMongoContext(collection string) (context.Context, *mongo.Collection) {
	client := InitiateMongoClient()
	db := client.Database(configs.MongoCollection)
	col := db.Collection(collection)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return ctx, col
}
