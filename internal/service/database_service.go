package service

import (
	"context"
	"fmt"
	"net/url"
	"pushnotification_services/internal/config"
	"pushnotification_services/internal/utilities"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
	mongoOnce   sync.Once
	mongoErr    error
)

// GetMongoDatabaseConnection 获取 MongoDB 连接（带重试）
func GetMongoDatabaseConnection() (*mongo.Client, error) {
	mongoOnce.Do(func() {
		creds := config.MongoDBCreds

		// URL 编码用户名和密码
		encodedUser := url.QueryEscape(creds.DatabaseUser)
		encodedPass := url.QueryEscape(creds.DatabasePassword)

		// 构建基础 URI
		baseURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s",
			encodedUser,
			encodedPass,
			creds.DatabaseHost,
			creds.DatabasePort,
			creds.DatabaseName,
		)

		// 添加认证参数
		params := url.Values{}
		authSource := creds.DatabaseAuthSource
		if authSource == "" {
			authSource = "admin" // 默认
		}
		params.Set("authSource", authSource)

		// 可选：指定认证机制（部分环境需要）
		// params.Set("authMechanism", "SCRAM-SHA-256")

		connectionString := baseURI + "?" + params.Encode()

		// 记录连接串（隐藏密码）
		masked := strings.Replace(connectionString, creds.DatabasePassword, "****", -1)
		utilities.Log(utilities.INFO, "MongoDB 连接字符串: %s", masked)

		clientOptions := options.Client().ApplyURI(connectionString)
		clientOptions.SetMinPoolSize(10)
		clientOptions.SetMaxPoolSize(100)
		clientOptions.SetConnectTimeout(10 * time.Second)
		clientOptions.SetServerSelectionTimeout(30 * time.Second)

		// 重试连接 (最多 5 次)
		var client *mongo.Client
		var err error
		for i := 0; i < 5; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			client, err = mongo.Connect(ctx, clientOptions)
			if err == nil {
				// 尝试 Ping
				pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer pingCancel()
				err = client.Ping(pingCtx, nil)
				if err == nil {
					break
				}
			}
			utilities.Log(utilities.WARN, "MongoDB 连接尝试 %d/5 失败: %s", i+1, err.Error())
			time.Sleep(2 * time.Second)
		}

		if err != nil {
			utilities.Log(utilities.ERROR, "MongoDB 最终连接失败: %s", err.Error())
			mongoErr = err
			return
		}

		utilities.Log(utilities.INFO, "MongoDB 连接成功，认证库: %s", authSource)
		mongoClient = client
	})

	return mongoClient, mongoErr
}

// GetMongoClient 返回已存在的客户端
func GetMongoClient() *mongo.Client {
	return mongoClient
}