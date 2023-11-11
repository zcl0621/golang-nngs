package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"higo-game-node/config"
	"log"
)

var mgoCli *mongo.Client

func InitMongoDB() {
	var err error
	clientOptions := options.Client().ApplyURI(
		fmt.Sprintf("mongodb://%s:%s@%s:%s/",
			config.Conf.MongoDB.User,
			config.Conf.MongoDB.Password,
			config.Conf.MongoDB.Host,
			config.Conf.MongoDB.Port,
		))

	// 连接到MongoDB
	mgoCli, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// 检查连接
	err = mgoCli.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func GetMgoCli() *mongo.Client {
	if mgoCli == nil {
		InitMongoDB()
	}
	return mgoCli
}
