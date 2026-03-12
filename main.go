package main

import (

	"fmt"

	"pushnotification_services/internal/api/middleware"
	"pushnotification_services/internal/api/routes"
	"pushnotification_services/internal/config"
	"pushnotification_services/internal/utilities"

	_ "pushnotification_services/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	Addr = ":8082"
	Port = 8082
)

// @title			推送通知服务
// @version		1.0
// @author			钟智强
// @description	提供推送通知相关的 API 接口，支持发送文本和图片通知，集成 OneSignal 推送服务，
// @description	提供 WebSocket 实时通知功能，支持向所有用户广播通知，
// @description	同时提供通知记录存储和管理功能，为应用提供完整的推送通知解决方案
// @BasePath		/
func main() {
	// gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.Use(gin.Recovery())
	router.Use(middleware.SecurityMiddleware())

	routes.General(router)
	router.GET(config.SWAGGER_DOCS, ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.InitWebSocketManager()
	routes.StandardPushNotification(router)
	routes.WebSocketRoutes(router)
	routes.Segmentation(router)
	routes.Announcement(router)
	routes.Advertisement(router)

	if err := router.Run(fmt.Sprintf(":%d", Port)); err != nil {
		utilities.Log(utilities.ERROR, "HTTP 服务启动失败: %v", err)
	}
}
