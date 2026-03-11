package handler

import (
	"net/http"
	"time"
	"pushnotification_services/internal/repositories"
	"pushnotification_services/internal/structure"
	"pushnotification_services/internal/utilities"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// CreateAdvertisement 创建新广告
// @Summary 创建新广告
// @Description 创建一个新的广告并保存到数据库。广告包含标题、内容、图片URL、开始时间、结束时间、语言区域等信息。
// @Description 
// @Description 必填字段：
// @Description - title: 广告标题，可为空字符串
// @Description - content: 广告内容，可为空字符串
// @Description - image_url: 广告图片URL，不能为空
// @Description - start_time: 广告开始时间，ISO 8601格式
// @Description - end_time: 广告结束时间，必须晚于开始时间
// @Description - locale: 语言区域，遵循 RFC 5646 标准，例如 zh、en、zh-CN、zh-TW
// @Description - action_url: 点击广告跳转URL（可选）
// @Description 
// @Description 注意事项：
// @Description - 广告将在开始时间和结束时间之间有效
// @Description - locale参数必须符合RFC 5646标准
// @Description - image_url为必填字段
// @Tags 广告管理
// @Accept json
// @Produce json
// @Param advertisement body structure.Advertisement true "广告信息"
// @Success 200 {object} map[string]interface{} "成功响应"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /advertisement/create [post]
func CreateAdvertisement() gin.HandlerFunc {
	return func(c *gin.Context) {
		var advertisement structure.Advertisement
		if err := c.ShouldBindJSON(&advertisement); err != nil {
			utilities.Log(utilities.ERROR, "解析请求参数失败: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid request parameters",
			})
			return
		}

		// 设置创建时间为当前时间
		if advertisement.CreatedAt.IsZero() {
			advertisement.CreatedAt = time.Now()
		}

		err := repositories.WriteAdvertisement(advertisement)
		if err != nil {
			utilities.Log(utilities.ERROR, "创建广告失败: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "创建广告失败",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "广告创建成功",
			"data":    advertisement,
		})
	}
}

// DeleteAdvertisement 删除广告
// @Summary 删除广告
// @Description 根据广告ID从数据库中永久删除指定的广告记录。此操作不可恢复，请谨慎使用。
// @Description 
// @Description 参数说明：
// @Description - id: 广告的唯一标识符（MongoDB ObjectID），必须存在于数据库中
// @Description 
// @Description 删除成功后，该广告将不再对用户可见，且无法恢复。
// @Description 如果指定的ID不存在，将返回404错误。
// @Tags 广告管理
// @Accept json
// @Produce json
// @Param id path string true "广告 ID"
// @Success 200 {object} map[string]interface{} "成功响应"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 404 {object} map[string]interface{} "广告不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /advertisement/delete/{id} [delete]
func DeleteAdvertisement() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "广告ID不能为空",
			})
			return
		}

		err := repositories.DeleteAdvertisement(id)
		if err != nil {
			utilities.Log(utilities.ERROR, "删除广告失败: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "删除广告失败",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "广告删除成功",
		})
	}
}

// GetLatestAdvertisement 获取最新广告
// @Summary 获取最新广告
// @Description 获取当前有效的、最新的广告信息。可以通过locale参数过滤特定语言的广告。
// @Description 
// @Description 功能说明：
// @Description - 返回当前时间在有效期内的最新广告
// @Description - 如果提供了locale参数，只返回匹配该语言的广告
// @Description - 如果没有找到广告，返回空数据
// @Description 
// @Description 参数说明：
// @Description - locale: 可选，语言区域过滤，支持的值：zh、en、zh-CN、zh-TW
// @Description   - zh: 简体中文
// @Description   - en: 英语
// @Description   - zh-CN: 中国大陆中文
// @Description   - zh-TW: 台湾中文
// @Description 
// @Description 使用示例：
// @Description - /advertisement/latest?locale=zh 获取最新的中文广告
// @Description - /advertisement/latest?locale=en 获取最新的英文广告
// @Description - /advertisement/latest 获取最新广告（不限语言）
// @Tags 广告管理
// @Accept json
// @Produce json
// @Param locale query string false "语言区域 (zh, en, zh-CN, zh-TW)"
// @Success 200 {object} map[string]interface{} "成功响应"
// @Failure 400 {object} map[string]interface{} "无效的locale参数"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /advertisement/latest [get]
func GetLatestAdvertisement() gin.HandlerFunc {
	return func(c *gin.Context) {
		locale := c.Query("locale")
		
		// 验证locale参数
		if locale != "" && locale != "zh" && locale != "en" && locale != "zh-CN" && locale != "zh-TW" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "无效的locale参数，只支持 'zh', 'en', 'zh-CN', 'zh-TW'",
			})
			return
		}

		advertisement, err := repositories.GetLatestAdvertisement(locale)
		if err != nil {
			utilities.Log(utilities.ERROR, "获取最新广告失败: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "获取最新广告失败",
			})
			return
		}

		if advertisement == nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  "success",
				"message": "未找到广告",
				"data":    nil,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "获取最新广告成功",
			"data":    advertisement,
		})
	}
}

// GetAllAdvertisements 获取所有广告
// @Summary 获取所有广告
// @Description 获取当前有效的所有广告列表，按创建时间降序排序。可以通过locale参数过滤特定语言的广告。
// @Description 
// @Description 功能说明：
// @Description - 返回当前时间在有效期内的所有广告
// @Description - 如果提供了locale参数，只返回匹配该语言的广告
// @Description - 按创建时间降序排序（最新的在前）
// @Description 
// @Description 参数说明：
// @Description - locale: 可选，语言区域过滤，支持的值：zh、en、zh-CN、zh-TW
// @Description   - zh: 简体中文
// @Description   - en: 英语
// @Description   - zh-CN: 中国大陆中文
// @Description   - zh-TW: 台湾中文
// @Description 
// @Description 使用示例：
// @Description - /advertisement/all?locale=zh 获取所有中文广告
// @Description - /advertisement/all?locale=en 获取所有英文广告
// @Description - /advertisement/all 获取所有广告（不限语言）
// @Tags 广告管理
// @Accept json
// @Produce json
// @Param locale query string false "语言区域 (zh, en, zh-CN, zh-TW)"
// @Success 200 {object} map[string]interface{} "成功响应"
// @Failure 400 {object} map[string]interface{} "无效的locale参数"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /advertisement/all [get]
func GetAllAdvertisements() gin.HandlerFunc {
	return func(c *gin.Context) {
		locale := c.Query("locale")
		
		// 验证locale参数
		if locale != "" && locale != "zh" && locale != "en" && locale != "zh-CN" && locale != "zh-TW" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "无效的locale参数，只支持 'zh', 'en', 'zh-CN', 'zh-TW'",
			})
			return
		}

		advertisements, err := repositories.GetAllAdvertisements(locale)
		if err != nil {
			utilities.Log(utilities.ERROR, "获取广告列表失败: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "获取广告列表失败",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "获取广告列表成功",
			"data":    advertisements,
		})
	}
}

// UpdateAdvertisement 更新广告
// @Summary 更新广告
// @Description 根据广告ID更新广告信息。可以部分更新广告字段，只更新提供的字段。
// @Description 
// @Description 参数说明：
// @Description - id: 广告的唯一标识符（MongoDB ObjectID），必须存在于数据库中
// @Description - updateData: 要更新的字段键值对，支持部分更新
// @Description 
// @Description 可更新字段：
// @Description - title: 广告标题
// @Description - content: 广告内容
// @Description - image_url: 广告图片URL
// @Description - start_time: 广告开始时间
// @Description - end_time: 广告结束时间
// @Description - locale: 语言区域
// @Description - action_url: 点击广告跳转URL
// @Description 
// @Description 注意事项：
// @Description - 不能更新_id和id字段
// @Description - 只更新提供的字段，未提供的字段保持不变
// @Description - 时间格式必须符合ISO 8601标准
// @Tags 广告管理
// @Accept json
// @Produce json
// @Param id path string true "广告 ID"
// @Param updateData body map[string]interface{} true "更新字段键值对"
// @Success 200 {object} map[string]interface{} "成功响应"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 404 {object} map[string]interface{} "广告不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /advertisement/update/{id} [put]
func UpdateAdvertisement() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "广告ID不能为空",
			})
			return
		}

		var updateData map[string]interface{}
		if err := c.ShouldBindJSON(&updateData); err != nil {
			utilities.Log(utilities.ERROR, "解析更新参数失败: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid update parameters",
			})
			return
		}

		// 移除_id字段，防止更新
		delete(updateData, "_id")
		delete(updateData, "id")

		err := repositories.UpdateAdvertisement(id, bson.M(updateData))
		if err != nil {
			utilities.Log(utilities.ERROR, "更新广告失败: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "更新广告失败",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "广告更新成功",
		})
	}
}