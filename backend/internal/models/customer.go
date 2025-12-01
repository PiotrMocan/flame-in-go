package models

import "gorm.io/gorm"

type Customer struct {
	gorm.Model
	Name        string  `json:"name" binding:"required"`
	Email       string  `json:"email"`
	Phone       string  `json:"phone"`
	CompanyID   uint    `json:"company_id"`
	Company     Company `json:"-"`
	FunnelID    *uint   `json:"funnel_id"`
	FunnelStage string  `json:"funnel_stage"`
}
