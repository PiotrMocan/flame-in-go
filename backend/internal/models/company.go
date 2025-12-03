package models

import "gorm.io/gorm"

type Company struct {
	gorm.Model
	Name      string     `json:"name" binding:"required"`
	Address   string     `json:"address"`
	Users     []User     `json:"users,omitempty"`
	Customers []Customer `json:"customers,omitempty"`
	FunnelID  *uint      `json:"funnel_id"`
	Funnel    *Funnel    `json:"funnel,omitempty"`
}
