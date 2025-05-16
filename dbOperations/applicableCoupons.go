package dboperations

import (
	"time"

	"github.com/farmako/models"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

func GetCouponsFromDb(ids, categories []string, orderTotal uint, timestamp time.Time, db *gorm.DB) ([]models.Coupon, error) {
	var coupons []models.Coupon
	query := "SELECT * from coupons where ( applicable_medicine_ids && ? OR applicable_categories && ? ) AND min_order_value <= ? AND expiry_date >= ?"
	err := db.Raw(query, pq.Array(ids), pq.Array(categories), orderTotal, timestamp).Scan(&coupons).Error
	if err != nil {
		return []models.Coupon{}, err
	}
	return coupons, nil
}

func GetPrice(id string, db *gorm.DB) (uint, error) {
	var price uint
	query := "SELECT price FROM products where id = ?"
	err := db.Raw(query, id).Scan(&price).Error
	if err != nil {
		return 0, err
	}
	return price, err

}
