package models

import "gorm.io/gorm"

type Company struct {
	gorm.Model
	Name      string     `json:"name" binding:"required"`
	Address   string     `json:"address"`
	Users     []User     `json:"users,omitempty"`
	Customers []Customer `json:"customers,omitempty"`
	Funnels   []Funnel   `json:"funnels,omitempty"`
}
