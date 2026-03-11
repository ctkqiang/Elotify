package structure

// OneSignalNotificationRequest 定义 OneSignal API 请求结构
type OneSignalNotificationRequest struct {
	AppID            string            `json:"app_id"`
	Contents         map[string]string `json:"contents"`
	Headings         map[string]string `json:"headings"`
	IncludedSegments []string          `json:"included_segments"`
	BigPicture       string            `json:"big_picture,omitempty"`
}

// OneSignalNotificationResponse 定义 OneSignal API 响应结构
type OneSignalNotificationResponse struct {
	ID         string   `json:"id"`
	Recipients int      `json:"recipients"`
	ExternalID string   `json:"external_id,omitempty"`
	Errors     []string `json:"errors,omitempty"`
}

// OneSignalClient 是用于与 OneSignal API 交互的客户端结构
type OneSignalClient struct {
	ApplicationId string
	APIKey        string
}
