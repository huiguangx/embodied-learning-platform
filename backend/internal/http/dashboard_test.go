package http

import("net/http/httptest"; "testing"; "github.com/gin-gonic/gin")
func TestDashboard(t *testing.T){ gin.SetMode(gin.TestMode); r:=NewRouter(); q:=httptest.NewRequest("GET","/api/v1/projects/00000000-0000-0000-0000-000000000001/dashboard",nil); q.Header.Set("X-User-ID","u"); q.Header.Set("X-Project-ID","00000000-0000-0000-0000-000000000001"); w:=httptest.NewRecorder(); r.ServeHTTP(w,q); if w.Code!=200 {t.Fatalf("got %d",w.Code)} }
