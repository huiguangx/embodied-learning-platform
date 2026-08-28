package http

import (
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
)

func TestErrorEnvelope(t *testing.T) {
	gin.SetMode(gin.TestMode); c, rec := gin.CreateTestContext(httptest.NewRecorder()); Error(c, 400, "BAD_INPUT", "bad", map[string]string{"field":"name"})
	if rec.Code != 400 || rec.Header().Get("Content-Type") == "" { t.Fatalf("invalid response: %d", rec.Code) }
}
