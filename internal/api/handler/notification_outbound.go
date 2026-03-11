package handler

import (
	"pushnotification_services/internal/database"
	"pushnotification_services/internal/repositories"
	"pushnotification_services/internal/service"
	"pushnotification_services/internal/structure"
	"pushnotification_services/internal/utilities"
)

func CreateBaseNotification(appID, title, message string) *structure.OneSignalNotificationRequest {
	return &structure.OneSignalNotificationRequest{
		AppID: appID,
		Contents: map[string]string{
			"en": message,
		},
		Headings: map[string]string{
			"en": title,
		},
	}
}

// SendGeneralNotification 发送通用通知
// 此函数仅用于发送通用通知，只包含标题和消息文本，不包含图片
// 通知将广播给所有用户，不针对特定分段
// 参数:
//
//	client: OneSignalClient 实例，用于与 OneSignal API 交互
//	content: NotificationContent 实例，包含通知的标题和消息内容
//
// 返回值:
//
//	database.NotificationResponse: 通知发送的响应信息，包含状态和内容
//	error: 发送过程中遇到的错误，如果成功则为 nil
func SendGeneralNotification(client *structure.OneSignalClient, content *structure.NotificationContent) (database.NotificationResponse, error) {
	utilities.Log(utilities.INFO, "正在发送通用通知")

	// 创建服务层客户端
	onesignalClient := service.NewOneSignalClient()

	requestBody := &structure.OneSignalNotificationRequest{
		AppID: client.ApplicationId,
		Contents: map[string]string{
			"en": content.Message,
		},
		Headings: map[string]string{
			"en": content.Title,
		},
		IncludedSegments: []string{"All"},
	}

	// 添加图片支持
	if content.ImageUrl != nil && *content.ImageUrl != "" {
		requestBody.BigPicture = *content.ImageUrl
	}

	// 发送通知
	apiResponse, err := onesignalClient.SendNotification(requestBody)
	if err != nil {
		utilities.Log(utilities.ERROR, "发送通用通知失败: %s", err.Error())
		return database.NotificationResponse{
			Status: database.StatusFailed,
		}, err
	}

	utilities.Log(utilities.INFO, "OneSignal API 响应: %+v", apiResponse)
	utilities.Log(utilities.INFO, "通用通知发送成功")

	if err := repositories.SaveRecord(&database.NotificationResponse{
		Status:  database.StatusSuccess,
		Content: content,
	}); err != nil {
		utilities.Log(utilities.ERROR, "保存通知记录失败: %s", err.Error())
		return database.NotificationResponse{
			Status: database.StatusFailed,
		}, err
	}

	return database.NotificationResponse{
		Status:  database.StatusSuccess,
		Content: content,
	}, nil
}
