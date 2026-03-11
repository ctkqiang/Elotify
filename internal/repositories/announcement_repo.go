package repositories

import (
	"context"
	"pushnotification_services/internal/config"
	"pushnotification_services/internal/service"
	"pushnotification_services/internal/structure"
	"pushnotification_services/internal/utilities"

	"go.mongodb.org/mongo-driver/bson"
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
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		utilities.Log(utilities.ERROR, "删除公告失败: %s", err.Error())
		return err
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

	var announcement structure.Announcement
	limit := int64(1)
	cursor, err := collection.Find(context.Background(), bson.M{}, &options.FindOptions{
		Sort:  bson.M{"created_at": -1},
		Limit: &limit,
	})
	if err != nil {
		utilities.Log(utilities.ERROR, "获取最新公告失败: %s", err.Error())
		return nil, err
	}
	defer cursor.Close(context.Background())

	if cursor.Next(context.Background()) {
		if err := cursor.Decode(&announcement); err != nil {
			utilities.Log(utilities.ERROR, "解析公告失败: %s", err.Error())
			return nil, err
		}
		return &announcement, nil
	}

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

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		utilities.Log(utilities.ERROR, "更新公告失败: %s", err.Error())
		return err
	}

	utilities.Log(utilities.INFO, "公告更新成功: %s", id)
	return nil
}
