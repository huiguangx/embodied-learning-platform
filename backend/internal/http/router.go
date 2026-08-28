package http

import (
 "github.com/gin-gonic/gin"
 "eip-platform/backend/internal/auth"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(auth.Middleware())
	r.GET("/healthz", healthz)
	api := r.Group("/api/v1/projects/:projectId")
	api.GET("/dashboard", dashboard)
	api.GET("/assets/images", images)
	api.GET("/assets/objects", objects)
	registerResourceRoutes(api)
	registerJobRoutes(api)
	registerDeliveryRoutes(api)
	return r
}

func setupRouterForTest() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return NewRouter()
}

func healthz(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
