package repositories

import (
	"context"
	"pushnotification_services/internal/config"
	"pushnotification_services/internal/database"
	"pushnotification_services/internal/service"
	"pushnotification_services/internal/utilities"
)

func SaveRecord(record *database.NotificationResponse) error {

	mongoClient, err := service.GetMongoDatabaseConnection()
	if err != nil {
		utilities.Log(utilities.ERROR, "MongoDB 连接失败: %s", err.Error())
		return err
	}

	mongodb := mongoClient.Database(config.MongoDBCreds.DatabaseName)
	collection := mongodb.Collection(config.COLLECTION_NOTIFICATIONS)

	_, err = collection.InsertOne(context.Background(), record)
	if err != nil {
		utilities.Log(utilities.ERROR, "保存通知记录失败: %s", err.Error())
		return err
	}

	utilities.Log(utilities.INFO, "通知记录保存成功")
	return nil
}
