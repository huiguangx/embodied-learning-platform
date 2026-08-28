package http

import (
	"crypto/rand"
	"fmt"
	"github.com/gin-gonic/gin"
)

func RequestID(c *gin.Context) string { if id := c.GetHeader("X-Request-ID"); id != "" { return id }; b := make([]byte, 16); _, _ = rand.Read(b); id := fmt.Sprintf("%x", b); c.Header("X-Request-ID", id); return id }

func Error(c *gin.Context, status int, code, message string, details any) { c.JSON(status, gin.H{"code":code,"message":message,"requestId":RequestID(c),"details":details}) }
