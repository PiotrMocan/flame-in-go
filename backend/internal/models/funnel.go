package models

import "gorm.io/gorm"

type Funnel struct {
	gorm.Model
	Name      string `json:"name" binding:"required"`
	CompanyID uint   `json:"company_id"`
}
