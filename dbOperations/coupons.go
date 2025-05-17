package dboperations

import (
	"errors"
	"sync"
	"time"

	"github.com/farmako/models"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

var mut sync.Mutex

func AddCouponToDb(coupon models.Coupon, db *gorm.DB) error {
	errCh := make(chan error)

	go func() {
		for _, id := range coupon.ApplicableMedicineIDs {
			var product models.Product
			err := db.First(&product, "id = ?", id).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				errCh <- errors.New("product with id " + id + " does not exist")
				close(errCh)
				return
			}
			if err != nil {
				errCh <- err
				close(errCh)
				return
			}
		}

		mut.Lock()
		defer mut.Unlock()

		var existingCoupon models.Coupon
		err := db.First(&existingCoupon, "coupon_code = ?", coupon.CouponCode).Error
		if err == nil {
			errCh <- errors.New("coupon with this code already exists")
			return
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			errCh <- err
			return
		}

		err = db.Create(&coupon).Error
		errCh <- err
	}()
	return <-errCh
}

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

func GetCouponDataFromDb(couponCode string, db *gorm.DB) (models.Coupon, error) {
	var coupon models.Coupon
	query := "SELECT * from coupons where coupon_code = ?"
	err := db.Raw(query, couponCode).Scan(&coupon).Error
	if err != nil {
		return models.Coupon{}, err
	}
	return coupon, nil
}
