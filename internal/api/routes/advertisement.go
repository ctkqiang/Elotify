package routes

import (
	"pushnotification_services/internal/api/handler"
	"pushnotification_services/internal/config"

	"github.com/gin-gonic/gin"
)

func Advertisement(router *gin.Engine) {
	public := router.Group(config.ADVERTISEMENT)
	{
		public.POST(config.ADVERTISEMENT_CREATE, handler.CreateAdvertisement())
		public.DELETE(config.ADVERTISEMENT_DELETE+"/:id", handler.DeleteAdvertisement())
		public.GET(config.ADVERTISEMENT_LATEST, handler.GetLatestAdvertisement())
		public.PUT(config.ADVERTISEMENT_UPDATE+"/:id", handler.UpdateAdvertisement())
		public.GET(config.ADVERTISEMENT_LIST_ALL + "&show_all=true", handler.GetAllAdvertisements())
	}
}
