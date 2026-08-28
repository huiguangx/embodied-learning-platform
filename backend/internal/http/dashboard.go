package http

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"eip-platform/backend/internal/repository"
)

func dashboard(c *gin.Context) { v, ok:=repository.GetDashboard(c.Param("projectId")); if !ok { Error(c,http.StatusNotFound,"PROJECT_NOT_FOUND","project not found",nil); return }; c.JSON(http.StatusOK,v) }
