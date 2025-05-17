package dboperations

import (
	"github.com/farmako/models"
	"gorm.io/gorm"
)

func GetProductsFromDb(db *gorm.DB) ([]models.Product, error) {
	var products []models.Product
	query := "SELECT * FROM products"
	err := db.Raw(query).Scan(&products).Error
	return products, err
}

func AddProduct(product models.Product, db *gorm.DB) error {
	return db.Create(&product).Error
}
