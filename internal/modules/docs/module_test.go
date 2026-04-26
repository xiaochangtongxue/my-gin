package docs

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestModuleRegistersSwaggerRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	NewModule().RegisterRoutes(engine)

	for _, route := range engine.Routes() {
		if route.Method == "GET" && route.Path == "/swagger/*any" {
			return
		}
	}
	t.Fatal("missing GET /swagger/*any route")
}
