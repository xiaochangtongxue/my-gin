package auth

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/xiaochangtongxue/my-gin/internal/handler"
)

func TestModuleRegistersAuthRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	module := NewModule(handler.NewAuthHandler(nil), handler.NewCaptchaHandler(nil))

	module.RegisterRoutes(engine)

	routes := map[string]bool{}
	for _, route := range engine.Routes() {
		routes[route.Method+" "+route.Path] = true
	}

	want := []string{
		"GET /api/v1/auth/captcha",
		"POST /api/v1/auth/register",
		"POST /api/v1/auth/login",
		"POST /api/v1/auth/refresh",
		"POST /api/v1/auth/logout",
	}
	for _, route := range want {
		if !routes[route] {
			t.Fatalf("missing route %s", route)
		}
	}
	if routes["POST /api/v1/login"] {
		t.Fatal("old /api/v1/login route should not be registered by auth module")
	}
}
