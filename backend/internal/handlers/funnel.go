package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mokan/flame-crm-backend/internal/db"
	"github.com/mokan/flame-crm-backend/internal/models"
)

type CreateFunnelInput struct {
	Name              string `json:"name" binding:"required"`
	NextFunnelIDs     []uint `json:"next_funnel_ids"`
	PreviousFunnelIDs []uint `json:"previous_funnel_ids"`
}

type UpdateFunnelInput struct {
	Name              string `json:"name"`
	NextFunnelIDs     []uint `json:"next_funnel_ids"`
	PreviousFunnelIDs []uint `json:"previous_funnel_ids"`
}

func GetFunnels(c *gin.Context) {
	var funnels []models.Funnel
	if err := db.DB.Preload("NextFunnels").Preload("PreviousFunnels").Find(&funnels).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, funnels)
}

func CreateFunnel(c *gin.Context) {
	var input CreateFunnelInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	funnel := models.Funnel{
		Name: input.Name,
	}

	if len(input.NextFunnelIDs) > 0 {
		var nextFunnels []*models.Funnel
		if err := db.DB.Where("id IN ?", input.NextFunnelIDs).Find(&nextFunnels).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid next funnel IDs"})
			return
		}
		funnel.NextFunnels = nextFunnels
	}

	if len(input.PreviousFunnelIDs) > 0 {
		var prevFunnels []*models.Funnel
		if err := db.DB.Where("id IN ?", input.PreviousFunnelIDs).Find(&prevFunnels).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid previous funnel IDs"})
			return
		}
		funnel.PreviousFunnels = prevFunnels
	}

	if err := db.DB.Create(&funnel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, funnel)
}

func UpdateFunnel(c *gin.Context) {
	id := c.Param("id")
	var funnel models.Funnel
	if err := db.DB.Preload("NextFunnels").Preload("PreviousFunnels").First(&funnel, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Funnel not found"})
		return
	}

	var input UpdateFunnelInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name != "" {
		funnel.Name = input.Name
	}

	if input.NextFunnelIDs != nil {
		var nextFunnels []*models.Funnel
		if len(input.NextFunnelIDs) > 0 {
			if err := db.DB.Where("id IN ?", input.NextFunnelIDs).Find(&nextFunnels).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid next funnel IDs"})
				return
			}
		}
		if err := db.DB.Model(&funnel).Association("NextFunnels").Replace(nextFunnels); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update next transitions"})
			return
		}
	}

	if input.PreviousFunnelIDs != nil {
		var prevFunnels []*models.Funnel
		if len(input.PreviousFunnelIDs) > 0 {
			if err := db.DB.Where("id IN ?", input.PreviousFunnelIDs).Find(&prevFunnels).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid previous funnel IDs"})
				return
			}
		}
		if err := db.DB.Model(&funnel).Association("PreviousFunnels").Replace(prevFunnels); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update previous transitions"})
			return
		}
	}

	db.DB.Save(&funnel)
	c.JSON(http.StatusOK, funnel)
}

func DeleteFunnel(c *gin.Context) {
	id := c.Param("id")
	var funnel models.Funnel
	if err := db.DB.First(&funnel, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Funnel not found"})
		return
	}

	if err := db.DB.Model(&funnel).Association("NextFunnels").Clear(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear next transitions"})
		return
	}
	if err := db.DB.Model(&funnel).Association("PreviousFunnels").Clear(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear previous transitions"})
		return
	}

	if err := db.DB.Delete(&funnel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Funnel deleted"})
}
