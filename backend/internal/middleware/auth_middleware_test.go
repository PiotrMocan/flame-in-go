package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mokan/flame-crm-backend/internal/auth"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupAuth      func() string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "No Authorization Header",
			setupAuth: func() string {
				return ""
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Authorization header is required",
		},
		{
			name: "Invalid Token Format",
			setupAuth: func() string {
				return "Bearer"
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid token format",
		},
		{
			name: "Invalid Token",
			setupAuth: func() string {
				return "Bearer invalidtoken123"
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid or expired token",
		},
		{
			name: "Valid Token",
			setupAuth: func() string {
				token, _ := auth.GenerateToken(1, "admin")
				return "Bearer " + token
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "Success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req, _ := http.NewRequest("GET", "/", nil)
			authHeader := tt.setupAuth()
			if authHeader != "" {
				req.Header.Set("Authorization", authHeader)
			}
			c.Request = req

			middleware := AuthMiddleware()

			handler := func(c *gin.Context) {
				middleware(c)
				if !c.IsAborted() {
					c.String(http.StatusOK, "Success")
				}
			}

			handler(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, w.Body.String(), tt.expectedBody)
			}

			if tt.expectedStatus == http.StatusOK {
				userID, exists := c.Get("user_id")
				assert.True(t, exists)
				assert.Equal(t, uint(1), userID)

				role, exists := c.Get("role")
				assert.True(t, exists)
				assert.Equal(t, "admin", role)
			}
		})
	}
}
