package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mokan/flame-crm-backend/internal/db"
	"github.com/mokan/flame-crm-backend/internal/models"
)

func GetCustomers(c *gin.Context) {
	var customers []models.Customer
	if err := db.DB.Preload("Company").Find(&customers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, customers)
}

func CreateCustomer(c *gin.Context) {
	var input models.Customer
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, input)
}

func UpdateCustomer(c *gin.Context) {
	id := c.Param("id")
	var customer models.Customer
	if err := db.DB.First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	var input models.UpdateCustomerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name != "" {
		customer.Name = input.Name
	}
	if input.Email != "" {
		customer.Email = input.Email
	}
	if input.Phone != "" {
		customer.Phone = input.Phone
	}

	if input.FunnelID != nil && customer.FunnelID != nil {
		if *input.FunnelID != *customer.FunnelID {
			var currentFunnel models.Funnel
			if err := db.DB.Preload("NextFunnels").First(&currentFunnel, *customer.FunnelID).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Current funnel state invalid"})
				return
			}

			valid := false
			for _, next := range currentFunnel.NextFunnels {
				if next.ID == *input.FunnelID {
					valid = true
					break
				}
			}
			if !valid {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid funnel transition"})
				return
			}
		}
	}

	customer.FunnelID = input.FunnelID

	if input.FunnelStage != "" {
		customer.FunnelStage = input.FunnelStage
	}

	db.DB.Save(&customer)
	c.JSON(http.StatusOK, customer)
}
