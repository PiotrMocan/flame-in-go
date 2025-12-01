package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegister_InvalidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		body         string
		expectedBody string
	}{
		{
			name: "Missing Name",
			body: `{"email": "test@example.com", "password": "password123"}`,
			expectedBody: "Error:Field validation for 'Name' failed on the 'required' tag",
		},
		{
			name: "Invalid Email",
			body: `{"name": "Test", "email": "not-an-email", "password": "password123"}`,
			expectedBody: "Error:Field validation for 'Email' failed on the 'email' tag",
		},
		{
			name: "Short Password",
			body: `{"name": "Test", "email": "test@example.com", "password": "123"}`,
			expectedBody: "Error:Field validation for 'Password' failed on the 'min' tag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/register", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			Register(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}

func TestLogin_InvalidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		body         string
		expectedBody string
	}{
		{
			name: "Missing Email",
			body: `{"password": "password123"}`,
			expectedBody: "Error:Field validation for 'Email' failed on the 'required' tag",
		},
		{
			name: "Invalid Email Format",
			body: `{"email": "bad-email", "password": "password123"}`,
			expectedBody: "Error:Field validation for 'Email' failed on the 'email' tag",
		},
		{
			name: "Missing Password",
			body: `{"email": "test@example.com"}`,
			expectedBody: "Error:Field validation for 'Password' failed on the 'required' tag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			Login(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}
