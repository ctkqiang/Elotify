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
	announcementsCollection = config.COLLECTION_ANNOUNCEMENTS
)

func WriteAnnouncement(announcement structure.Announcement) error {
	client, err := service.GetMongoDatabaseConnection()
	if err != nil {
		return err
	}

	collection := client.Database("pushnotification").Collection(announcementsCollection)

	if announcement.ID == "" {
		announcement.ID = primitive.NewObjectID().Hex()
	}

	filter := bson.M{"_id": announcement.ID}
	update := bson.M{
		"$set": announcement,
	}
	upsert := true
	options := &options.UpdateOptions{Upsert: &upsert}

	_, err = collection.UpdateOne(context.Background(), filter, update, options)
	if err != nil {
		utilities.Log(utilities.ERROR, "写入公告失败: %s", err.Error())
		return err
	}

	utilities.Log(utilities.INFO, "公告写入成功: %s", announcement.ID)
	return nil
}

func DeleteAnnouncement(id string) error {
	client, err := service.GetMongoDatabaseConnection()
	if err != nil {
		return err
	}

	collection := client.Database("pushnotification").Collection(announcementsCollection)

	filter := bson.M{"_id": id}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		utilities.Log(utilities.ERROR, "删除公告失败: %s", err.Error())
		return err
	}

	if result.DeletedCount == 0 {
		utilities.Log(utilities.WARN, "未找到要删除的公告: %s", id)
		return nil
	}

	utilities.Log(utilities.INFO, "公告删除成功: %s", id)
	return nil
}

func GetLatestAnnouncement() (*structure.Announcement, error) {
	client, err := service.GetMongoDatabaseConnection()
	if err != nil {
		return nil, err
	}

	collection := client.Database("pushnotification").Collection(announcementsCollection)

	// 获取所有公告按创建时间倒序，查找第一个有效的公告
	var announcement structure.Announcement
	cursor, err := collection.Find(context.Background(), bson.M{}, &options.FindOptions{
		Sort:  bson.M{"created_at": -1},
	})
	if err != nil {
		utilities.Log(utilities.ERROR, "获取最新公告失败: %s", err.Error())
		return nil, err
	}
	defer cursor.Close(context.Background())

	now := time.Now()
	for cursor.Next(context.Background()) {
		if err := cursor.Decode(&announcement); err != nil {
			utilities.Log(utilities.ERROR, "解析公告失败: %s", err.Error())
			return nil, err
		}
		
		// 检查公告是否在有效期内
		if !announcement.StartedAt.IsZero() && !announcement.ExpiresAt.IsZero() {
			if announcement.StartedAt.Before(now) && announcement.ExpiresAt.After(now) {
				utilities.Log(utilities.INFO, "找到最新有效公告: ID=%s", announcement.ID)
				return &announcement, nil
			} else {
				utilities.Log(utilities.DEBUG, "跳过过期公告: ID=%s", announcement.ID)
				continue
			}
		} else {
			// 如果没有时间字段，直接返回
			utilities.Log(utilities.INFO, "找到最新公告（无时间限制）: ID=%s", announcement.ID)
			return &announcement, nil
		}
	}

	utilities.Log(utilities.INFO, "未找到当前有效的公告")
	return nil, nil
}

func GetAllAnnouncements() ([]structure.Announcement, error) {
	client, err := service.GetMongoDatabaseConnection()
	if err != nil {
		return nil, err
	}

	collection := client.Database("pushnotification").Collection(announcementsCollection)

	var announcements []structure.Announcement
	cursor, err := collection.Find(context.Background(), bson.M{}, &options.FindOptions{
		Sort: bson.M{"created_at": -1},
	})
	if err != nil {
		utilities.Log(utilities.ERROR, "获取所有公告失败: %s", err.Error())
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &announcements); err != nil {
		utilities.Log(utilities.ERROR, "解析公告列表失败: %s", err.Error())
		return nil, err
	}

	utilities.Log(utilities.INFO, "获取公告列表成功，共 %d 条记录", len(announcements))
	return announcements, nil
}

func UpdateAnnouncement(id string, announcement structure.Announcement) error {
	client, err := service.GetMongoDatabaseConnection()
	if err != nil {
		return err
	}

	collection := client.Database("pushnotification").Collection(announcementsCollection)

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": announcement,
	}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		utilities.Log(utilities.ERROR, "更新公告失败: %s", err.Error())
		return err
	}

	if result.MatchedCount == 0 {
		utilities.Log(utilities.WARN, "未找到要更新的公告: %s", id)
		return nil
	}

	utilities.Log(utilities.INFO, "公告更新成功: %s", id)
	return nil
}
