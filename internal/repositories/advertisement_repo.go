package repositories

import (
	"context"
	"pushnotification_services/internal/config"
	"pushnotification_services/internal/service"
	"pushnotification_services/internal/structure"
	"pushnotification_services/internal/utilities"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	advertisementsCollection = config.COLLECTION_ADVERTISEMENTS
)

func WriteAdvertisement(advertisement structure.Advertisement) error {
	client, err := service.GetMongoDatabaseConnection()
	if err != nil {
		return err
	}

	collection := client.Database("pushnotification").Collection(advertisementsCollection)

	filter := bson.M{"_id": advertisement.ID}
	update := bson.M{
		"$set": advertisement,
	}
	upsert := true
	opts := &options.UpdateOptions{Upsert: &upsert}

	_, err = collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		utilities.Log(utilities.ERROR, "写入广告失败: %s", err.Error())
		return err
	}

	utilities.Log(utilities.INFO, "广告写入成功: %s", advertisement.ID.Hex())
	return nil
}

func DeleteAdvertisement(id string) error {
	client, err := service.GetMongoDatabaseConnection()
	if err != nil {
		return err
	}

	collection := client.Database("pushnotification").Collection(advertisementsCollection)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utilities.Log(utilities.ERROR, "无效的广告ID格式: %s", err.Error())
		return err
	}

	filter := bson.M{"_id": objID}
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		utilities.Log(utilities.ERROR, "删除广告失败: %s", err.Error())
		return err
	}

	utilities.Log(utilities.INFO, "广告删除成功: %s", id)
	return nil
}

func GetLatestAdvertisement(locale string) (*structure.Advertisement, error) {
	client, err := service.GetMongoDatabaseConnection()
	if err != nil {
		return nil, err
	}

	collection := client.Database("pushnotification").Collection(advertisementsCollection)

	// 构建查询过滤器
	filter := bson.M{}
	if locale != "" {
		filter["locale"] = locale
	}

	// 只获取当前时间在有效期内的广告
	now := time.Now()
	filter["start_time"] = bson.M{"$lte": now}
	filter["end_time"] = bson.M{"$gte": now}

	var advertisement structure.Advertisement
	limit := int64(1)
	cursor, err := collection.Find(context.Background(), filter, &options.FindOptions{
		Sort:  bson.M{"created_at": -1},
		Limit: &limit,
	})
	if err != nil {
		utilities.Log(utilities.ERROR, "获取最新广告失败: %s", err.Error())
		return nil, err
	}
	defer cursor.Close(context.Background())

	if cursor.Next(context.Background()) {
		if err := cursor.Decode(&advertisement); err != nil {
			utilities.Log(utilities.ERROR, "解析广告失败: %s", err.Error())
			return nil, err
		}
		return &advertisement, nil
	}

	return nil, nil
}

func GetAllAdvertisements(locale string) ([]structure.Advertisement, error) {
	client, err := service.GetMongoDatabaseConnection()
	if err != nil {
		return nil, err
	}

	collection := client.Database("pushnotification").Collection(advertisementsCollection)

	// 构建查询过滤器
	filter := bson.M{}
	if locale != "" {
		filter["locale"] = locale
	}

	// 只获取当前时间在有效期内的广告
	now := time.Now()
	filter["start_time"] = bson.M{"$lte": now}
	filter["end_time"] = bson.M{"$gte": now}

	cursor, err := collection.Find(context.Background(), filter, &options.FindOptions{
		Sort: bson.M{"created_at": -1},
	})
	if err != nil {
		utilities.Log(utilities.ERROR, "获取广告列表失败: %s", err.Error())
		return nil, err
	}
	defer cursor.Close(context.Background())

	var advertisements []structure.Advertisement
	if err := cursor.All(context.Background(), &advertisements); err != nil {
		utilities.Log(utilities.ERROR, "解析广告列表失败: %s", err.Error())
		return nil, err
	}

	return advertisements, nil
}

func UpdateAdvertisement(id string, updateData bson.M) error {
	client, err := service.GetMongoDatabaseConnection()
	if err != nil {
		return err
	}

	collection := client.Database("pushnotification").Collection(advertisementsCollection)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utilities.Log(utilities.ERROR, "无效的广告ID格式: %s", err.Error())
		return err
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": updateData}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		utilities.Log(utilities.ERROR, "更新广告失败: %s", err.Error())
		return err
	}

	utilities.Log(utilities.INFO, "广告更新成功: %s", id)
	return nil
}
