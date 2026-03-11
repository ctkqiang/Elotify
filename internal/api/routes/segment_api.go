package routes

import (
	"pushnotification_services/internal/api/handler"
	"pushnotification_services/internal/config"

	"github.com/gin-gonic/gin"
)

// Segmentation 注册分段相关路由
func Segmentation(router *gin.Engine) {
	public := router.Group(config.SEGMENTATION)
	{
		public.GET(config.SEGMENT_LIST_ALL, GetAllSegments())
	}
}

// GetAllSegments 获取所有 OneSignal 分段
//
//	@Summary		获取所有分段
//	@Description	获取 OneSignal 中的所有分段信息
//	@Tags			分段管理
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}	"成功响应"
//	@Failure		500	{object}	map[string]interface{}	"服务器内部错误"
//	@Router			/segment/all [get]
func GetAllSegments() gin.HandlerFunc {
	return func(c *gin.Context) {
		segments := handler.ListAllSegments()
		if segments == nil {
			c.JSON(500, gin.H{
				"status":  "error",
				"message": "获取分段失败",
			})
			return
		}

		c.JSON(200, gin.H{
			"status":  "success",
			"message": "获取分段成功",
			"data":    segments,
		})
	}
}
