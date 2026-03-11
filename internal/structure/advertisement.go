package structure

import (
	"time"

	"github.com/kamva/mgm/v3"
)

// Advertisement 广告数据结构
//
//	@Description	广告数据结构，包含广告的所有信息
type Advertisement struct {
	mgm.DefaultModel `bson:",inline" swaggerignore:"true"`

	// Title 广告标题，可选
	//	@example	"春节促销"
	Title *string `json:"title" bson:"title"`

	// Content 广告内容，可选
	//	@example	"春节期间全场商品8折优惠"
	Content *string `json:"content" bson:"content"`

	// ImageURL 广告图片URL，必填
	//	@example	"https://example.com/advertisement.jpg"
	ImageURL string `json:"image_url" bson:"image_url"`

	// StartTime 广告开始时间，ISO 8601格式
	//	@example	"2024-02-01T00:00:00Z"
	StartTime time.Time `json:"start_time" bson:"start_time"`

	// EndTime 广告结束时间，ISO 8601格式，必须晚于开始时间
	//	@example	"2024-02-15T23:59:59Z"
	EndTime time.Time `json:"end_time" bson:"end_time"`

	// Locale 语言区域，遵循 RFC 5646 标准
	//	@example	"zh"
	Locale string `json:"locale" bson:"locale"` // 遵循 RFC 5646 标准，例如 en、zh-CN

	// ActionUrl 点击广告跳转URL，可选
	//	@example	"https://example.com/promo"
	ActionUrl *string `json:"action_url,omitempty" bson:"action_url,omitempty"`

	// CreatedAt 广告创建时间，ISO 8601格式
	//	@example	"2024-01-15T10:00:00Z"
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
