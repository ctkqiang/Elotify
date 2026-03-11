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
	databaseName             = config.MongoDBCreds.DatabaseName
)

func WriteAdvertisement(advertisement *structure.Advertisement) error {
	client, err := service.GetMongoDatabaseConnection()
	if err != nil {
		return err
	}

	collection := client.Database(databaseName).Collection(advertisementsCollection)

	// 使用mgm的ID字段，避免重复ID问题
	if advertisement.ID.IsZero() {
		// 为新文档生成ObjectID
		advertisement.SetID(primitive.NewObjectID())
		utilities.Log(utilities.INFO, "手动生成广告ID: %s", advertisement.ID.Hex())
		
		// 默认设置开始时间为当前时间，结束时间为24小时后
		if advertisement.StartTime.IsZero() {
			advertisement.StartTime = time.Now()
		}
		if advertisement.EndTime.IsZero() {
			advertisement.EndTime = time.Now().Add(24 * time.Hour)
		}
		
		// 记录插入前的广告信息
		utilities.Log(utilities.INFO, "插入前广告信息: ID=%s, title=%v, locale=%s, start=%v, end=%v", 
			advertisement.ID.Hex(), advertisement.Title, advertisement.Locale, advertisement.StartTime, advertisement.EndTime)
		
		// 使用InsertOne创建新文档
		result, err := collection.InsertOne(context.Background(), advertisement)
		if err != nil {
			utilities.Log(utilities.ERROR, "创建广告失败: %s", err.Error())
			return err
		}
		
		// 记录InsertedID的类型和值
		utilities.Log(utilities.INFO, "InsertOne返回的InsertedID: %v, 类型: %T", result.InsertedID, result.InsertedID)
		
		// 验证插入结果
		count, err := collection.CountDocuments(context.Background(), bson.M{"_id": advertisement.ID})
		if err != nil {
			utilities.Log(utilities.ERROR, "验证广告插入失败: %s", err.Error())
		} else {
			utilities.Log(utilities.INFO, "广告插入验证: 找到 %d 条记录", count)
		}
		
		utilities.Log(utilities.INFO, "广告创建成功，最终ID: %s", advertisement.ID.Hex())
	} else {
		// 使用UpdateOne更新现有文档
		filter := bson.M{"_id": advertisement.ID}
		update := bson.M{"$set": advertisement}
		upsert := true
		opts := &options.UpdateOptions{Upsert: &upsert}

		_, err = collection.UpdateOne(context.Background(), filter, update, opts)
		if err != nil {
			utilities.Log(utilities.ERROR, "更新广告失败: %s", err.Error())
			return err
		}

		utilities.Log(utilities.INFO, "广告更新成功: %s", advertisement.ID.Hex())
	}
	
	return nil
}

func DeleteAdvertisement(id string) error {
	client, err := service.GetMongoDatabaseConnection()
	if err != nil {
		return err
	}

	collection := client.Database(databaseName).Collection(advertisementsCollection)

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

	collection := client.Database(databaseName).Collection(advertisementsCollection)

	filter := bson.M{}
	if locale != "" {
		filter["locale"] = locale
	}
	
	utilities.Log(utilities.INFO, "广告查询条件: locale=%s, filter=%+v", locale, filter)

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
		utilities.Log(utilities.INFO, "找到最新广告: ID=%s, start=%v, end=%v, locale=%s", 
			advertisement.ID.Hex(), advertisement.StartTime, advertisement.EndTime, advertisement.Locale)
		return &advertisement, nil
	}

	utilities.Log(utilities.INFO, "未找到符合条件的最新广告")
	return nil, nil
}

func GetAllAdvertisements(locale string, showAll bool) ([]structure.Advertisement, error) {
	client, err := service.GetMongoDatabaseConnection()
	if err != nil {
		utilities.Log(utilities.ERROR, "MongoDB连接失败: %s", err.Error())
		return nil, err
	}

	collection := client.Database(databaseName).Collection(advertisementsCollection)

	// Build filter: only locale if provided (time filtering removed as requested)
	filter := bson.M{}
	if locale != "" {
		filter["locale"] = locale
	} else {
		filter["locale"] = "en"
	}

	// Debug: log the filter being used
	utilities.Log(utilities.INFO, "广告查询过滤器: %+v", filter)

	// Find with sort by created_at descending
	cursor, err := collection.Find(context.Background(), filter, &options.FindOptions{
		Sort: bson.M{"created_at": -1},
	})

	if err != nil {
		utilities.Log(utilities.ERROR, "获取广告列表失败: %s", err.Error())
		return nil, err
	}

	defer cursor.Close(context.Background())

	// Decode cursor results into advertisements slice
	var advertisements []structure.Advertisement
	if err := cursor.All(context.Background(), &advertisements); err != nil {
		utilities.Log(utilities.ERROR, "解析广告列表失败: %s", err.Error())
		return nil, err
	}

	utilities.Log(utilities.INFO, "广告列表查询完成，找到 %d 条记录", len(advertisements))
	return advertisements, nil
}

func UpdateAdvertisement(id string, updateData bson.M) error {
	client, err := service.GetMongoDatabaseConnection()
	if err != nil {
		return err
	}

	collection := client.Database(databaseName).Collection(advertisementsCollection)

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
