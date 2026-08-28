package auth

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const ContextKey = "auth"

type AuthContext struct { UserID, ProjectID string; Role Role }

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/healthz" { c.Next(); return }
		user, project := c.GetHeader("X-User-ID"), c.GetHeader("X-Project-ID")
		role := Role(c.GetHeader("X-Role"))
		if os.Getenv("APP_ENV") == "development" {
			if user == "" { user = "demo-user" }
			if project == "" { project = "00000000-0000-0000-0000-000000000001" }
			if role == "" { role = Engineer }
		}
		if user == "" || project == "" { c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code":"UNAUTHENTICATED","message":"user and project context are required"}); return }
		c.Set(ContextKey, AuthContext{UserID:user, ProjectID:project, Role:role})
		if requested := c.Param("projectId"); requested != "" && requested != project { c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code":"PROJECT_SCOPE_DENIED","message":"project access denied"}); return }
		c.Next()
	}
}

func Context(c *gin.Context) (AuthContext, bool) { v, ok := c.Get(ContextKey); if !ok { return AuthContext{}, false }; a, ok := v.(AuthContext); return a, ok }

func RequireWrite() gin.HandlerFunc {
	return func(c *gin.Context) { a, ok := Context(c); if !ok || !a.Role.CanWrite() { c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code":"FORBIDDEN","message":"write permission required"}); return }; c.Next() }
}
