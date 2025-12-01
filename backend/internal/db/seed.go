package db

import (
	"fmt"

	"github.com/mokan/flame-crm-backend/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	if err := db.AutoMigrate(&models.Company{}, &models.User{}, &models.Customer{}, &models.Funnel{}); err != nil {
		fmt.Println("Error running migrations during seed:", err)
		return
	}

	var count int64
	if err := db.Model(&models.User{}).Count(&count).Error; err != nil {
		fmt.Println("Error checking user count:", err)
		return
	}

	if count > 0 {
		fmt.Println("Database already seeded with users.")
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return
	}

	admin := models.User{
		Name:     "Admin User",
		Email:    "admin@example.com",
		Password: string(password),
		Role:     models.RoleAdmin,
	}

	if err := db.Create(&admin).Error; err != nil {
		fmt.Printf("Failed to seed admin: %v\n", err)
	} else {
		fmt.Println("------------------------------------------------")
		fmt.Println("Seeding successful!")
		fmt.Println("Admin User created:")
		fmt.Println("Email:    admin@example.com")
		fmt.Println("Password: password123")
		fmt.Println("------------------------------------------------")
	}
}