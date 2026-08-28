package http

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"eip-platform/backend/internal/repository"
)

func images(c *gin.Context) { page,_:=strconv.Atoi(c.DefaultQuery("page","1")); if page<1 {page=1}; items,total,ok:=repository.Images(c.Param("projectId"),c.Query("namespace"),c.Query("name"),page); if !ok {Error(c,http.StatusNotFound,"PROJECT_NOT_FOUND","project not found",nil);return}; c.JSON(http.StatusOK,gin.H{"items":items,"nextCursor":"","total":total}) }
func objects(c *gin.Context) { items,total,ok:=repository.Objects(c.Param("projectId"),c.Query("bucket"),c.Query("prefix")); if !ok {Error(c,http.StatusNotFound,"PROJECT_NOT_FOUND","project not found",nil);return}; c.JSON(http.StatusOK,gin.H{"items":items,"nextCursor":"","total":total}) }
