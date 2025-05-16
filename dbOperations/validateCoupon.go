package dboperations

import (
	"github.com/farmako/models"
	"gorm.io/gorm"
)

func GetCouponDataFromDb(couponCode string, db *gorm.DB) (models.Coupon, error) {
	var coupon models.Coupon
	query := "SELECT * from coupons where coupon_code = ?"
	err := db.Raw(query, couponCode).Scan(&coupon).Error
	if err != nil {
		return models.Coupon{}, err
	}
	return coupon, nil
}
