package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mokan/flame-crm-backend/internal/db"
	"github.com/mokan/flame-crm-backend/internal/models"
)

func GetCompanies(c *gin.Context) {
	var companies []models.Company
	if err := db.DB.Preload("Users").Preload("Customers").Find(&companies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, companies)
}

func CreateCompany(c *gin.Context) {
	var input models.Company
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

func UpdateCompany(c *gin.Context) {
	id := c.Param("id")
	var company models.Company
	if err := db.DB.First(&company, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	var input models.Company
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.DB.Model(&company).Updates(input)
	c.JSON(http.StatusOK, company)
}
