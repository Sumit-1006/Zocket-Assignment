package controllers

import (
	"net/http"
	"product-management-system/models"
	"product-management-system/services"

	"github.com/gin-gonic/gin"
)

func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save product in database
	if err := services.SaveProduct(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save product"})
		return
	}

	// Add image URLs to RabbitMQ queue
	go services.QueueProductImages(product.Images)

	c.JSON(http.StatusOK, gin.H{"message": "Product created successfully"})
}

func GetProductByID(c *gin.Context) {
	id := c.Param("id")
	product, err := services.GetProductByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func GetProducts(c *gin.Context) {
	// Optional filtering
	userID := c.Query("user_id")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")
	name := c.Query("name")

	products, err := services.GetProducts(userID, minPrice, maxPrice, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, products)
}
