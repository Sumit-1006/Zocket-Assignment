package main

import (
	"log"
	"product-management-system/config"
	"product-management-system/controllers"
	"product-management-system/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the configurations
	initConfigurations()

	// Setup the Gin router
	router := gin.Default()

	// Define the API routes
	setupRoutes(router)

	// Start the server
	log.Println("Starting the server on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// initConfigurations initializes all configurations for the app
func initConfigurations() {
	// Connect to the PostgreSQL database
	log.Println("Connecting to the database...")
	config.ConnectDB()

	// Initialize RabbitMQ
	log.Println("Initializing RabbitMQ...")
	services.InitRabbitMQ()

	// Initialize Redis
	log.Println("Initializing Redis...")
	config.ConnectRedis()
}

// setupRoutes sets up all the API routes
func setupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.POST("/products", controllers.CreateProduct)
		api.GET("/products/:id", controllers.GetProductByID)
		api.GET("/products", controllers.GetProducts)
	}
}
