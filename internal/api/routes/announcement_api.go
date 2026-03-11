package routes

import (
	"net/http"
	"pushnotification_services/internal/config"
	"pushnotification_services/internal/repositories"
	"pushnotification_services/internal/structure"
	"pushnotification_services/internal/utilities"

	"github.com/gin-gonic/gin"
)

func Announcement(router *gin.Engine) {
	public := router.Group(config.ANNOUNCEMENT)
	{
		public.POST(config.ANNOUNCEMENT_CREATE, CreateAnnouncement())
		public.DELETE(config.ANNOUNCEMENT_DELETE, DeleteAnnouncement())
		public.GET(config.ANNOUNCEMENT_LATEST, GetLatestAnnouncement())
		public.PUT(config.ANNOUNCEMENT_UPDATE, UpdateAnnouncement())
		public.GET(config.ANNOUNCEMENT_LIST_ALL, GetAllAnnouncements())
	}
}

// CreateAnnouncement 创建新公告
//
//	@Summary		创建新公告
//	@Description	创建一个新的公告并保存到数据库。公告用于向用户发布重要信息，包括系统维护、活动通知、节假日安排等。
//	@Description
//	@Description	必填字段：
//	@Description	- type: 公告类型，可选值：HOLIDAY（节假日）、EVENT（活动通知）
//	@Description	- message: 公告内容，支持多语言文本
//	@Description	- priority: 优先级，可选值：HIGH（高）、NORMAL（正常）、LOW（低）
//	@Description	- started_at: 开始时间，ISO 8601格式
//	@Description	- expires_at: 过期时间，必须晚于开始时间
//	@Description
//	@Description	示例请求：
//	@Description	```json
//	@Description	{
//	@Description	"id": "announcement_001",
//	@Description	"type": "EVENT",
//	@Description	"message": "系统维护通知：将于本周末进行系统升级维护",
//	@Description	"priority": "HIGH",
//	@Description	"created_at": "2024-01-15T10:00:00Z",
//	@Description	"started_at": "2024-01-20T09:00:00Z",
//	@Description	"expires_at": "2024-01-21T18:00:00Z"
//	@Description	}
//	@Description	```
//	@Description
//	@Description	注意事项：
//	@Description	- 公告ID必须唯一，建议使用有意义的标识符
//	@Description	- 时间格式必须符合ISO 8601标准
//	@Description	- 开始时间必须早于过期时间
//	@Description	- 消息内容建议简洁明了，不超过500字符
//	@Tags			公告管理
//	@Accept			json
//	@Produce		json
//	@Param			announcement	body		structure.Announcement	true	"公告信息"
//	@Success		200				{object}	map[string]interface{}	"成功响应"
//	@Failure		400				{object}	map[string]interface{}	"请求参数错误"
//	@Failure		500				{object}	map[string]interface{}	"服务器内部错误"
//	@Router			/announcement/create [post]
func CreateAnnouncement() gin.HandlerFunc {
	return func(c *gin.Context) {
		var announcement structure.Announcement
		if err := c.ShouldBindJSON(&announcement); err != nil {
			utilities.Log(utilities.ERROR, "解析请求参数失败: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid request parameters",
			})
			return
		}

		err := repositories.WriteAnnouncement(announcement)
		if err != nil {
			utilities.Log(utilities.ERROR, "创建公告失败: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "创建公告失败",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "公告创建成功",
			"data":    announcement,
		})
	}
}

// DeleteAnnouncement 删除公告
//
//	@Summary		删除公告
//	@Description	根据公告ID从数据库中永久删除指定的公告记录。此操作不可恢复，请谨慎使用。
//	@Description
//	@Description	参数说明：
//	@Description	- id: 公告的唯一标识符，必须存在于数据库中
//	@Description
//	@Description	删除成功后，该公告将不再对用户可见，且无法恢复。
//	@Description	如果指定的ID不存在，将返回404错误。
//	@Tags			公告管理
//	@Accept			json
//	@Produce		json
//	@Param			id	query		string					true	"公告 ID"
//	@Success		200	{object}	map[string]interface{}	"成功响应"
//	@Failure		400	{object}	map[string]interface{}	"请求参数错误"
//	@Failure		500	{object}	map[string]interface{}	"服务器内部错误"
//	@Router			/announcement/delete [delete]
func DeleteAnnouncement() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "ID 不能为空",
			})
			return
		}

		err := repositories.DeleteAnnouncement(id)
		if err != nil {
			utilities.Log(utilities.ERROR, "删除公告失败: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "删除公告失败",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "公告删除成功",
		})
	}
}

// GetLatestAnnouncement 获取最新公告
//
//	@Summary		获取最新公告
//	@Description	获取系统中最新发布的公告信息。该接口返回按创建时间排序的最新一条公告记录。
//	@Description
//	@Description	返回的公告信息包含完整的数据结构，包括：
//	@Description	- 公告ID、类型、内容、优先级等基本信息
//	@Description	- 创建时间、开始时间、过期时间等时间信息
//	@Description
//	@Description	如果没有找到任何公告，将返回404错误。
//	@Description	此接口常用于应用启动时获取最新的重要通知。
//	@Tags			公告管理
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}	"成功响应"
//	@Failure		404	{object}	map[string]interface{}	"没有找到公告"
//	@Failure		500	{object}	map[string]interface{}	"服务器内部错误"
//	@Router			/announcement/latest [get]
func GetLatestAnnouncement() gin.HandlerFunc {
	return func(c *gin.Context) {
		announcement, err := repositories.GetLatestAnnouncement()
		if err != nil {
			utilities.Log(utilities.ERROR, "获取最新公告失败: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "获取最新公告失败",
			})
			return
		}

		if announcement == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "没有找到公告",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "获取最新公告成功",
			"data":    announcement,
		})
	}
}

// UpdateAnnouncement 更新公告
//
//	@Summary		更新公告
//	@Description	根据公告ID更新现有的公告信息。可以修改公告的所有字段，包括类型、内容、优先级和时间信息。
//	@Description
//	@Description	参数说明：
//	@Description	- id: 要更新的公告ID，必须存在于数据库中
//	@Description	- announcement: 完整的公告对象，包含所有需要更新的字段
//	@Description
//	@Description	更新规则：
//	@Description	- 开始时间必须早于过期时间
//	@Description	- 消息内容不能为空
//	@Description	- 优先级必须在有效范围内（HIGH、NORMAL、LOW）
//	@Description	- 类型必须在有效范围内（HOLIDAY、EVENT）
//	@Description
//	@Description	更新成功后，将返回更新后的完整公告信息。
//	@Description	如果指定的ID不存在，将返回404错误。
//	@Tags			公告管理
//	@Accept			json
//	@Produce		json
//	@Param			id				query		string					true	"公告 ID"
//	@Param			announcement	body		structure.Announcement	true	"公告信息"
//	@Success		200				{object}	map[string]interface{}	"成功响应"
//	@Failure		400				{object}	map[string]interface{}	"请求参数错误"
//	@Failure		500				{object}	map[string]interface{}	"服务器内部错误"
//	@Router			/announcement/update [put]
func UpdateAnnouncement() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "ID 不能为空",
			})
			return
		}

		var announcement structure.Announcement
		if err := c.ShouldBindJSON(&announcement); err != nil {
			utilities.Log(utilities.ERROR, "解析请求参数失败: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid request parameters",
			})
			return
		}

		err := repositories.UpdateAnnouncement(id, announcement)
		if err != nil {
			utilities.Log(utilities.ERROR, "更新公告失败: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "更新公告失败",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "公告更新成功",
			"data":    announcement,
		})
	}
}

// GetAllAnnouncements 获取所有公告
//
//	@Summary		获取所有公告
//	@Description	获取系统中所有的公告信息，按创建时间降序排序。返回的公告列表包含每个公告的完整信息。
//	@Description
//	@Description	返回数据结构：
//	@Description	- 数组形式，每个元素是一个完整的公告对象
//	@Description	- 按创建时间从新到旧排序
//	@Description	- 包含公告的所有字段信息
//	@Description
//	@Description	使用场景：
//	@Description	- 管理后台查看所有公告
//	@Description	- 批量处理公告数据
//	@Description	- 数据分析和统计
//	@Description
//	@Description	如果系统中没有任何公告，将返回空数组。
//	@Tags			公告管理
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}	"成功响应"
//	@Failure		500	{object}	map[string]interface{}	"服务器内部错误"
//	@Router			/announcement/all [get]
func GetAllAnnouncements() gin.HandlerFunc {
	return func(c *gin.Context) {
		announcements, err := repositories.GetAllAnnouncements()
		if err != nil {
			utilities.Log(utilities.ERROR, "获取所有公告失败: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "获取所有公告失败",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "获取所有公告成功",
			"data":    announcements,
		})
	}
}
