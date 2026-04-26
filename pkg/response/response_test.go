package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	apperrors "github.com/xiaochangtongxue/my-gin/pkg/errors"
)

func TestErrorMapsBusinessCodesToHTTPStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   int
	}{
		{"validation", apperrors.New(CodeInvalidParam, "bad input"), http.StatusBadRequest, CodeInvalidParam},
		{"unauthorized", apperrors.New(CodeUnauthorized, "login required"), http.StatusUnauthorized, CodeUnauthorized},
		{"forbidden", apperrors.New(CodePermissionDenied, "denied"), http.StatusForbidden, CodePermissionDenied},
		{"rate limited", apperrors.New(CodeRateLimitExceeded, "slow down"), http.StatusTooManyRequests, CodeRateLimitExceeded},
		{"captcha", apperrors.New(CodeCaptchaError, "captcha failed"), http.StatusBadRequest, CodeCaptchaError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			Error(c, tt.err)

			if w.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", w.Code, tt.wantStatus)
			}

			var body Response
			if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}
			if body.Code != tt.wantCode {
				t.Fatalf("body.Code = %d, want %d", body.Code, tt.wantCode)
			}
			if body.Timestamp == 0 {
				t.Fatal("Timestamp = 0, want populated")
			}
		})
	}
}
