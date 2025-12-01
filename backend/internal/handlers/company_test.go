package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateCompany_InvalidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jsonBody := `{"address": "123 St"}`
	req, _ := http.NewRequest("POST", "/companies", bytes.NewBufferString(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	CreateCompany(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Error:Field validation for 'Name' failed on the 'required' tag")
}

func TestCreateCompany_MalformedJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jsonBody := `{"name": "Acme", "address": `
	req, _ := http.NewRequest("POST", "/companies", bytes.NewBufferString(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	CreateCompany(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
