package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
)

func TestMiddlewareRequiresContextOutsideDevelopment(t *testing.T) {
	gin.SetMode(gin.TestMode); r := gin.New(); r.Use(Middleware()); r.GET("/", func(c *gin.Context) { c.Status(200) })
	rec := httptest.NewRecorder(); r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if rec.Code != http.StatusUnauthorized { t.Fatalf("got %d", rec.Code) }
}

func TestEngineerCannotWrite(t *testing.T) { if Engineer.CanWrite() { t.Fatal("engineer can write") }; if !Operator.CanWrite() { t.Fatal("operator cannot write") } }
