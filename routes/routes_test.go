package routes

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/gin-gonic/gin"
)

func TestReportPageFallback(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	index := []byte(`<!doctype html><html><body><div id="root"></div></body></html>`)
	SetRoutes(router, index, http.FS(fstest.MapFS{}))
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/review-mining/example", nil))
	if response.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", response.Code)
	}
	if !strings.HasPrefix(response.Result().Header.Get("Content-Type"), "text/html") {
		t.Errorf("Content-Type = %q, want text/html", response.Result().Header.Get("Content-Type"))
	}
	if response.Body.String() != string(index) {
		t.Errorf("response must contain exactly one copy of the application HTML: %q", response.Body.String())
	}
}
