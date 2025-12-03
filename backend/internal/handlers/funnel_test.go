package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mokan/flame-crm-backend/internal/db"
	"github.com/mokan/flame-crm-backend/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	var err error
	testDB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		fmt.Printf("Failed to connect to test database: %v\n", err)
		os.Exit(1)
	}

	err = testDB.AutoMigrate(&models.Company{}, &models.User{}, &models.Customer{}, &models.Funnel{})
	if err != nil {
		fmt.Printf("Failed to migrate test database: %v\n", err)
		os.Exit(1)
	}

	db.DB = testDB

	code := m.Run()

	sqlDB, _ := testDB.DB()
	sqlDB.Close()

	os.Exit(code)
}

func clearTable(t *testing.T) {
	if err := testDB.Exec("DELETE FROM funnel_transitions;").Error; err != nil {
		t.Fatalf("Failed to clear funnel_transitions: %v", err)
	}
	if err := testDB.Exec("DELETE FROM customers;").Error; err != nil {
		t.Fatalf("Failed to clear customers: %v", err)
	}
	if err := testDB.Exec("DELETE FROM funnels;").Error; err != nil {
		t.Fatalf("Failed to clear funnels: %v", err)
	}
	if err := testDB.Exec("DELETE FROM users;").Error; err != nil {
		t.Fatalf("Failed to clear users: %v", err)
	}
	if err := testDB.Exec("DELETE FROM companies;").Error; err != nil {
		t.Fatalf("Failed to clear companies: %v", err)
	}
}

func createTestCompanyAndUser(t *testing.T) (models.Company, models.User) {
	clearTable(t)

	company := models.Company{Name: "Test Company"}
	assert.NoError(t, testDB.Create(&company).Error)

	user := models.User{
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  "password123",
		CompanyID: &company.ID,
		Role:      models.RoleAdmin,
	}
	assert.NoError(t, testDB.Create(&user).Error)
	return company, user
}

func performRequest(r http.Handler, method, path string, body interface{}) *httptest.ResponseRecorder {
	var req *http.Request
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		req, _ = http.NewRequest(method, path, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(method, path, nil)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/funnels", GetFunnels)
	r.POST("/funnels", CreateFunnel)
	r.PUT("/funnels/:id", UpdateFunnel)
	r.PUT("/customers/:id", UpdateCustomer)
	return r
}

func TestFunnelCreationAndRetrieval(t *testing.T) {
	r := setupRouter()
	createTestCompanyAndUser(t)

	createInput := CreateFunnelInput{Name: "Prospect"}
	w := performRequest(r, "POST", "/funnels", createInput)
	assert.Equal(t, http.StatusOK, w.Code)

	var createdFunnel models.Funnel
	err := json.Unmarshal(w.Body.Bytes(), &createdFunnel)
	assert.NoError(t, err)
	assert.Equal(t, "Prospect", createdFunnel.Name)
	assert.Greater(t, createdFunnel.ID, uint(0))

	w = performRequest(r, "GET", "/funnels", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var funnels []models.Funnel
	err = json.Unmarshal(w.Body.Bytes(), &funnels)
	assert.NoError(t, err)
	assert.Len(t, funnels, 1)
	assert.Equal(t, "Prospect", funnels[0].Name)
}

func TestFunnelUpdate(t *testing.T) {
	r := setupRouter()
	createTestCompanyAndUser(t)

	funnel1 := models.Funnel{Name: "Prospect"}
	funnel2 := models.Funnel{Name: "Contacted"}
	funnel3 := models.Funnel{Name: "Qualified"}
	assert.NoError(t, testDB.Create(&funnel1).Error)
	assert.NoError(t, testDB.Create(&funnel2).Error)
	assert.NoError(t, testDB.Create(&funnel3).Error)

	updateNameInput := UpdateFunnelInput{Name: "Initial Prospect"}
	w := performRequest(r, "PUT", fmt.Sprintf("/funnels/%d", funnel1.ID), updateNameInput)
	assert.Equal(t, http.StatusOK, w.Code)

	var updatedFunnel models.Funnel
	err := json.Unmarshal(w.Body.Bytes(), &updatedFunnel)
	assert.NoError(t, err)
	assert.Equal(t, "Initial Prospect", updatedFunnel.Name)
	assert.Equal(t, funnel1.ID, updatedFunnel.ID)

	updateTransitionsInput := UpdateFunnelInput{
		Name:          "Initial Prospect",
		NextFunnelIDs: []uint{funnel2.ID, funnel3.ID},
	}
	w = performRequest(r, "PUT", fmt.Sprintf("/funnels/%d", funnel1.ID), updateTransitionsInput)
	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &updatedFunnel)
	assert.NoError(t, err)
	assert.Equal(t, "Initial Prospect", updatedFunnel.Name)
	assert.Len(t, updatedFunnel.NextFunnels, 2)

	var dbFunnel models.Funnel
	err = testDB.Preload("NextFunnels").First(&dbFunnel, funnel1.ID).Error
	assert.NoError(t, err)
	assert.Len(t, dbFunnel.NextFunnels, 2)
	assert.Contains(t, []uint{dbFunnel.NextFunnels[0].ID, dbFunnel.NextFunnels[1].ID}, funnel2.ID)
	assert.Contains(t, []uint{dbFunnel.NextFunnels[0].ID, dbFunnel.NextFunnels[1].ID}, funnel3.ID)

	updateRemoveTransitionsInput := UpdateFunnelInput{
		Name:          "Initial Prospect",
		NextFunnelIDs: []uint{funnel2.ID},
	}
	w = performRequest(r, "PUT", fmt.Sprintf("/funnels/%d", funnel1.ID), updateRemoveTransitionsInput)
	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &updatedFunnel)
	assert.NoError(t, err)
	assert.Len(t, updatedFunnel.NextFunnels, 1)

	err = testDB.Preload("NextFunnels").First(&dbFunnel, funnel1.ID).Error
	assert.NoError(t, err)
	assert.Len(t, dbFunnel.NextFunnels, 1)
	assert.Equal(t, funnel2.ID, dbFunnel.NextFunnels[0].ID)
}

func TestCustomerFunnelTransition(t *testing.T) {
	r := setupRouter()
	company, _ := createTestCompanyAndUser(t)

	funnelA := models.Funnel{Name: "Funnel A"}
	funnelB := models.Funnel{Name: "Funnel B"}
	funnelC := models.Funnel{Name: "Funnel C"}
	assert.NoError(t, testDB.Create(&funnelA).Error)
	assert.NoError(t, testDB.Create(&funnelB).Error)
	assert.NoError(t, testDB.Create(&funnelC).Error)

	updateA := UpdateFunnelInput{NextFunnelIDs: []uint{funnelB.ID}}
	w := performRequest(r, "PUT", fmt.Sprintf("/funnels/%d", funnelA.ID), updateA)
	assert.Equal(t, http.StatusOK, w.Code)

	updateB := UpdateFunnelInput{NextFunnelIDs: []uint{funnelC.ID}}
	w = performRequest(r, "PUT", fmt.Sprintf("/funnels/%d", funnelB.ID), updateB)
	assert.Equal(t, http.StatusOK, w.Code)

	updateC := UpdateFunnelInput{NextFunnelIDs: []uint{funnelA.ID}}
	w = performRequest(r, "PUT", fmt.Sprintf("/funnels/%d", funnelC.ID), updateC)
	assert.Equal(t, http.StatusOK, w.Code)

	customerFunnelAID := &funnelA.ID
	customer := models.Customer{Name: "Test Customer", CompanyID: company.ID, FunnelID: customerFunnelAID}
	assert.NoError(t, testDB.Create(&customer).Error)

	customerFunnelBID := &funnelB.ID
	updateCustomerInput := models.UpdateCustomerInput{Name: customer.Name, FunnelID: customerFunnelBID}
	w = performRequest(r, "PUT", fmt.Sprintf("/customers/%d", customer.ID), updateCustomerInput)
	assert.Equal(t, http.StatusOK, w.Code)

	var updatedCustomer models.Customer
	err := json.Unmarshal(w.Body.Bytes(), &updatedCustomer)
	assert.NoError(t, err)
	assert.Equal(t, *customerFunnelBID, *updatedCustomer.FunnelID)

	customerFunnelAID = &funnelA.ID
	updateCustomerInput = models.UpdateCustomerInput{Name: updatedCustomer.Name, FunnelID: customerFunnelAID}
	w = performRequest(r, "PUT", fmt.Sprintf("/customers/%d", updatedCustomer.ID), updateCustomerInput)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid funnel transition")

	customerFunnelCID := &funnelC.ID
	updateCustomerInput = models.UpdateCustomerInput{Name: updatedCustomer.Name, FunnelID: customerFunnelCID}
	w = performRequest(r, "PUT", fmt.Sprintf("/customers/%d", updatedCustomer.ID), updateCustomerInput)
	assert.Equal(t, http.StatusOK, w.Code)
	err = json.Unmarshal(w.Body.Bytes(), &updatedCustomer)
	assert.NoError(t, err)
	assert.Equal(t, *customerFunnelCID, *updatedCustomer.FunnelID)

	clearTable(t)
	company, _ = createTestCompanyAndUser(t)
	funnelX := models.Funnel{Name: "Funnel X"}
	assert.NoError(t, testDB.Create(&funnelX).Error)
	customerNoFunnel := models.Customer{Name: "Customer No Funnel", CompanyID: company.ID, FunnelID: nil}
	assert.NoError(t, testDB.Create(&customerNoFunnel).Error)

	customerFunnelXID := &funnelX.ID
	updateCustomerInput = models.UpdateCustomerInput{Name: customerNoFunnel.Name, FunnelID: customerFunnelXID}
	w = performRequest(r, "PUT", fmt.Sprintf("/customers/%d", customerNoFunnel.ID), updateCustomerInput)
	assert.Equal(t, http.StatusOK, w.Code)
	err = json.Unmarshal(w.Body.Bytes(), &updatedCustomer)
	assert.NoError(t, err)
	assert.Equal(t, *customerFunnelXID, *updatedCustomer.FunnelID)

	customerNoFunnel2 := models.Customer{Name: "Customer No Funnel 2", CompanyID: company.ID, FunnelID: nil}
	assert.NoError(t, testDB.Create(&customerNoFunnel2).Error)
	updateCustomerInput = models.UpdateCustomerInput{Name: customerNoFunnel2.Name, FunnelID: nil}
	w = performRequest(r, "PUT", fmt.Sprintf("/customers/%d", customerNoFunnel2.ID), updateCustomerInput)
	assert.Equal(t, http.StatusOK, w.Code)
	err = json.Unmarshal(w.Body.Bytes(), &updatedCustomer)
	assert.NoError(t, err)
	assert.Nil(t, updatedCustomer.FunnelID)

	customerInFunnel := models.Customer{Name: "Customer In Funnel", CompanyID: company.ID, FunnelID: customerFunnelXID}
	assert.NoError(t, testDB.Create(&customerInFunnel).Error)
	updateCustomerInput = models.UpdateCustomerInput{Name: customerInFunnel.Name, FunnelID: nil}
	w = performRequest(r, "PUT", fmt.Sprintf("/customers/%d", customerInFunnel.ID), updateCustomerInput)
	assert.Equal(t, http.StatusOK, w.Code)
	err = json.Unmarshal(w.Body.Bytes(), &updatedCustomer)
	assert.NoError(t, err)
	assert.Nil(t, updatedCustomer.FunnelID)
}

func TestFunnelNotFound(t *testing.T) {
	r := setupRouter()
	createTestCompanyAndUser(t)

	updateInput := UpdateFunnelInput{Name: "Non Existent"}
	w := performRequest(r, "PUT", "/funnels/9999", updateInput)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Funnel not found")
}

func TestInvalidNextFunnelIDs(t *testing.T) {
	r := setupRouter()
	createTestCompanyAndUser(t)

	funnel1 := models.Funnel{Name: "Prospect"}
	assert.NoError(t, testDB.Create(&funnel1).Error)

	createInput := CreateFunnelInput{Name: "Invalid Funnel", NextFunnelIDs: []uint{9999}}
	w := performRequest(r, "POST", "/funnels", createInput)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid next funnel IDs")

	updateInput := UpdateFunnelInput{NextFunnelIDs: []uint{9999}}
	w = performRequest(r, "PUT", fmt.Sprintf("/funnels/%d", funnel1.ID), updateInput)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid next funnel IDs")
}

func generateDummyToken(userID uint) (string, error) {
	return fmt.Sprintf("Bearer dummy-token-for-user-%d", userID), nil
}

func init() {
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		fmt.Printf("Failed to load UTC location: %v\n", err)
	}
	time.Local = loc
}
