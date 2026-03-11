package routes

import (
	"net/http"
	"pushnotification_services/internal/api/handler"
	"pushnotification_services/internal/config"
	"pushnotification_services/internal/structure"
	"pushnotification_services/internal/utilities"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	request structure.NotificationContent
)

func StandardPushNotification(router *gin.Engine) {
	public := router.Group(config.SEND_TEXT_PUSH_NOTIFICATION_HEAD)
	{
		public.POST(config.SEND_TEXT_PUSH_NOTIFICATION, SendTextPushNotification())
		public.POST(config.SEND_TEXT_AND_IMAGE_PUSH_NOTIFICATION, SendTextAndImagePushNotification())
	}
}

// SendTextPushNotification 注册标准推送通知路由
//
//	@Summary		发送文本推送通知
//	@Description	发送通用文本推送通知到所有用户。此端点用于发送面向全体用户的通知，如系统公告、重要更新等。
//	@Description	警告：请勿使用此端点发送用户特定的通知，因为它会广播给所有用户。对于用户特定的通知，请使用专门的用户定向通知端点。
//	@Tags			通知
//	@Accept			json
//	@Produce		json
//	@Param			notification	body		structure.NotificationContent	true	"通知内容"
//	@Success		200				{object}	map[string]interface{}			"成功响应"
//	@Failure		400				{object}	map[string]interface{}			"请求参数错误"
//	@Failure		500				{object}	map[string]interface{}			"服务器内部错误"
//	@Router			/push/text [post]
//	@Example		curl -X POST http://localhost:8080/push/text \
//	@Example		-H "Content-Type: application/json" \
//	@Example		-d '{"title": "测试通知", "message": "这是一条测试通知"}'
func SendTextPushNotification() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request structure.NotificationContent
		if err := c.ShouldBindJSON(&request); err != nil {
			utilities.Log(utilities.ERROR, "解析请求参数失败: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid request parameters",
			})
			return
		}

		// 设置推送时间为当前时间
		request.AuditTrail.PushedAt = time.Now().Format(time.RFC3339)

		client := &structure.OneSignalClient{
			ApplicationId: config.OneSignalCreds.AppID,
			APIKey:        config.OneSignalCreds.APIKey,
		}

		response, err := handler.SendGeneralNotification(client, &request)
		if err != nil {
			utilities.Log(utilities.ERROR, "发送通知失败: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "发送通知失败",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "通知发送成功",
			"data":    response,
		})
	}
}

// SendTextAndImagePushNotification 发送文本和图片推送通知
//
//	@Summary		发送文本和图片推送通知
//	@Description	发送包含图片的推送通知到所有用户。此端点用于发送面向全体用户的通知，如系统公告、重要更新等，支持添加图片。
//	@Description	警告：请勿使用此端点发送用户特定的通知，因为它会广播给所有用户。对于用户特定的通知，请使用专门的用户定向通知端点。
//	@Tags			通知
//	@Accept			json
//	@Produce		json
//	@Param			notification	body		structure.NotificationContent	true	"通知内容"
//	@Success		200				{object}	map[string]interface{}			"成功响应"
//	@Failure		400				{object}	map[string]interface{}			"请求参数错误"
//	@Failure		500				{object}	map[string]interface{}			"服务器内部错误"
//	@Router			/push/text-image [post]
//	@Example		curl -X POST http://localhost:8080/push/text-image \
//	@Example		-H "Content-Type: application/json" \
//	@Example		-d '{"title": "测试通知", "message": "这是一条包含图片的测试通知", "image_url": "https://example.com/image.jpg"}'
func SendTextAndImagePushNotification() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request structure.NotificationContent
		if err := c.ShouldBindJSON(&request); err != nil {
			utilities.Log(utilities.ERROR, "解析请求参数失败: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid request parameters",
			})
			return
		}

		// 设置推送时间为当前时间
		request.AuditTrail.PushedAt = time.Now().Format(time.RFC3339)

		client := &structure.OneSignalClient{
			ApplicationId: config.OneSignalCreds.AppID,
			APIKey:        config.OneSignalCreds.APIKey,
		}

		response, err := handler.SendGeneralNotification(client, &request)
		if err != nil {
			utilities.Log(utilities.ERROR, "发送通知失败: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "发送通知失败",
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "通知发送成功",
			"data":    response,
		})
	}
}
