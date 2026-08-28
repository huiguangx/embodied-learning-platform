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
	return r
}

func setupRouterForTest() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return NewRouter()
}

func healthz(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
