package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"pushnotification_services/internal/config"
	"pushnotification_services/internal/structure"
	"pushnotification_services/internal/utilities"
)

// OneSignalEndpoint 是 OneSignal API 的端点常量
const OneSignalEndpoint = "https://api.onesignal.com/notifications"

// OneSignalClient 是用于与 OneSignal API 交互的客户端结构
type OneSignalClient struct {
	ApplicationId string
	APIKey        string
	HTTPClient    *http.Client
}

// NewOneSignalClient 创建一个新的 OneSignal 客户端
func NewOneSignalClient() *OneSignalClient {
	var (
		appID  = config.OneSignalCreds.AppID
		apiKey = config.OneSignalCreds.APIKey
	)

	if appID == "" {
		utilities.Log(utilities.ERROR, "环境变量中缺少 OneSignal AppID")
	}

	if apiKey == "" {
		utilities.Log(utilities.ERROR, "环境变量中缺少 OneSignal APIKey")
	}

	return &OneSignalClient{
		ApplicationId: appID,
		APIKey:        apiKey,
		HTTPClient:    &http.Client{},
	}
}

// SendNotification 发送通知到 OneSignal API
func (c *OneSignalClient) SendNotification(requestBody *structure.OneSignalNotificationRequest) (*structure.OneSignalNotificationResponse, error) {
	// 转换为 JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		utilities.Log(utilities.ERROR, "JSON 序列化失败: %s", err.Error())
		return nil, err
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", OneSignalEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		utilities.Log(utilities.ERROR, "创建 HTTP 请求失败: %s", err.Error())
		return nil, err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+c.APIKey)

	// 发送请求
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		utilities.Log(utilities.ERROR, "发送 HTTP 请求失败: %s", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	// 解析响应
	var apiResponse structure.OneSignalNotificationResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		utilities.Log(utilities.ERROR, "解析 API 响应失败: %s", err.Error())
		return nil, err
	}

	return &apiResponse, nil
}
