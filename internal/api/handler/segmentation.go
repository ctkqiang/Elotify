package handler

import (
	"encoding/json"
	"net/http"
	"pushnotification_services/internal/config"
	"pushnotification_services/internal/structure"
	"pushnotification_services/internal/utilities"
)

// ListAllSegments 获取所有 OneSignal 分段
func ListAllSegments() *[]structure.Segments {
	url := "https://api.onesignal.com/apps/" + config.OneSignalCreds.AppID + "/segments?limit=302"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		utilities.Log(utilities.ERROR, "创建请求失败: %s", err.Error())
		return nil
	}

	req.Header.Set("Authorization", config.OneSignalCreds.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utilities.Log(utilities.ERROR, "发送请求失败: %s", err.Error())
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		utilities.Log(utilities.ERROR, "API 请求失败，状态码: %d", resp.StatusCode)
		return nil
	}

	var response struct {
		Segments []structure.Segments `json:"segments"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		utilities.Log(utilities.ERROR, "解析响应失败: %s", err.Error())
		return nil
	}

	return &response.Segments
}
