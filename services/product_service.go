package services

import (
	"database/sql"
	"errors"
	"product-management-system/config"
	"product-management-system/models"
)

// SaveProduct saves a product into the PostgreSQL database
func SaveProduct(product *models.Product) error {
	query := `
		INSERT INTO products (user_id, product_name, product_description, product_images, product_price)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := config.DB.Exec(query, product.UserID, product.Name, product.Description, pq.Array(product.Images), product.Price)
	if err != nil {
		return err
	}
	return nil
}

// GetProductByID retrieves a product by its ID from the database
func GetProductByID(productID string) (*models.Product, error) {
	query := `
		SELECT id, user_id, product_name, product_description, product_images, compressed_product_images, product_price
		FROM products
		WHERE id = $1
	`

	var product models.Product
	err := config.DB.QueryRow(query, productID).Scan(
		&product.ID, &product.UserID, &product.Name, &product.Description,
		pq.Array(&product.Images), pq.Array(&product.CompressedImages), &product.Price,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("product not found")
	}
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// GetProducts retrieves all products for a specific user ID, with optional filters
func GetProducts(userID, minPrice, maxPrice, name string) ([]models.Product, error) {
	query := `
		SELECT id, user_id, product_name, product_description, product_images, compressed_product_images, product_price
		FROM products
		WHERE user_id = $1
	`
	var filters []interface{}
	filters = append(filters, userID)

	// Add filters dynamically
	if minPrice != "" {
		query += " AND product_price >= $2"
		filters = append(filters, minPrice)
	}
	if maxPrice != "" {
		query += " AND product_price <= $3"
		filters = append(filters, maxPrice)
	}
	if name != "" {
		query += " AND product_name ILIKE $4"
		filters = append(filters, "%"+name+"%")
	}

	rows, err := config.DB.Query(query, filters...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ID, &product.UserID, &product.Name, &product.Description,
			pq.Array(&product.Images), pq.Array(&product.CompressedImages), &product.Price,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
import (
	"database/sql"
	"errors"
	"product-management-system/config"
	"product-management-system/models"

	"github.com/lib/pq"
)

