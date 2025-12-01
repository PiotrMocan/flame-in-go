package models

import "gorm.io/gorm"

type Role string

const (
	RoleAdmin       Role = "admin"
	RoleSales       Role = "sales"
	RoleHeadOfSales Role = "head_of_sales"
)

type User struct {
	gorm.Model
	Name      string  `json:"name" binding:"required"`
	Email     string  `json:"email" gorm:"uniqueIndex" binding:"required,email"`
	Password  string  `json:"-"`
	Role      Role    `json:"role" binding:"required"`
	CompanyID *uint   `json:"company_id"`
	Company   Company `json:"company,omitempty"`
}
