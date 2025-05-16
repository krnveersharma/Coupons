package helpers

import (
	"errors"
	"fmt"

	requestschemas "github.com/farmako/RequestSchemas"
	"github.com/farmako/models"
)

func IsCouponValid(requestData requestschemas.ValidaTeCouponData, coupon models.Coupon, totalPrice uint) error {
	var isValid bool

	if coupon.MinOrderValue > float64(totalPrice) {
		return fmt.Errorf("Minimum price is %s", coupon.MinOrderValue)
	}
	for i := 0; i < len(requestData.CartItems); i++ {
		if contains(coupon.ApplicableCategories, requestData.CartItems[i].Category) || contains(coupon.ApplicableMedicineIDs, requestData.CartItems[i].ID) {
			isInValidTime, err := IsWithinValidTimeWindow(requestData.TimeStamp, coupon.ValidTimeWindow)
			if err == nil {
				isValid = true && isInValidTime
			}
		}
	}
	if !isValid {
		return errors.New("Coupon not valid")
	}
	return nil
}

func GetDiscountedPrice(coupon models.Coupon, price uint) (uint, uint) {

	if coupon.DiscountTarget == "charges" {
		return 0, uint(coupon.DiscountValue)
	}
	if coupon.DiscountType == "percentage" {
		return price * uint(coupon.DiscountValue) / 100, 0
	}
	return uint(coupon.DiscountValue), 0
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
