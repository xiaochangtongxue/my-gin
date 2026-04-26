package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuthSkipPathsUseNewAuthPrefix(t *testing.T) {
	gin.SetMode(gin.TestMode)
	InitAuth([]string{"/api/v1/auth/login"}, nil, nil)

	engine := gin.New()
	engine.Use(Auth())
	engine.POST("/api/v1/auth/login", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})
	engine.GET("/api/v1/private", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	loginReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	loginResp := httptest.NewRecorder()
	engine.ServeHTTP(loginResp, loginReq)
	if loginResp.Code != http.StatusNoContent {
		t.Fatalf("login status = %d, want %d", loginResp.Code, http.StatusNoContent)
	}

	privateReq := httptest.NewRequest(http.MethodGet, "/api/v1/private", nil)
	privateResp := httptest.NewRecorder()
	engine.ServeHTTP(privateResp, privateReq)
	if privateResp.Code != http.StatusUnauthorized {
		t.Fatalf("private status = %d, want %d", privateResp.Code, http.StatusUnauthorized)
	}
}
