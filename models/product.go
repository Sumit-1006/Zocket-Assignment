package models

type Product struct {
	ID                string   `json:"id"`
	UserID            string   `json:"user_id"`
	Name              string   `json:"product_name"`
	Description       string   `json:"product_description"`
	Images            []string `json:"product_images"`
	CompressedImages  []string `json:"compressed_product_images"`
	Price             float64  `json:"product_price"`
}
