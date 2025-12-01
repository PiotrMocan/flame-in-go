package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
	"github.com/mokan/flame-crm-backend/internal/db"
	"github.com/mokan/flame-crm-backend/internal/handlers"
	"github.com/mokan/flame-crm-backend/internal/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using OS env vars or defaults")
	}

	db.ConnectDatabase()

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/companies", handlers.GetCompanies)
		protected.POST("/companies", handlers.CreateCompany)
		protected.PUT("/companies/:id", handlers.UpdateCompany)

		protected.GET("/users", handlers.GetUsers)
		protected.POST("/users", handlers.CreateUser)
		protected.PUT("/users/:id", handlers.UpdateUser)

		protected.GET("/customers", handlers.GetCustomers)
		protected.POST("/customers", handlers.CreateCustomer)
		protected.PUT("/customers/:id", handlers.UpdateCustomer)
	}

	r.Run(":8080")
}