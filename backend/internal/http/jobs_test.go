package http
import("net/http/httptest";"testing";"github.com/gin-gonic/gin")
func TestJobsList(t *testing.T){gin.SetMode(gin.TestMode);r:=NewRouter();q:=httptest.NewRequest("GET","/api/v1/projects/p/training-jobs",nil);q.Header.Set("X-User-ID","u");q.Header.Set("X-Project-ID","p");w:=httptest.NewRecorder();r.ServeHTTP(w,q);if w.Code!=200{t.Fatalf("got %d",w.Code)}}
