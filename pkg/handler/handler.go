package handler

import (
	"banner/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("")
	{
		auth.POST("/register", h.register)
		auth.GET("/login", h.login)
	}

	banner := router.Group("", h.userIdentity)
	{
		banner.POST("/banner", h.createBanner)
		banner.PATCH("/banner/:id", h.updateBanner)
		banner.DELETE("/banner/:id", h.deleteBanner)
		banner.GET("/banner", h.getAllBanners)
		banner.GET("/user_banner", h.getUserBanner)

	}

	return router
}
