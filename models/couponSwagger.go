package models

import (
	"time"
)

// CouponSwagger is a Swagger-compatible representation of Coupon
// swagger:model
type CouponSwagger struct {
	// Coupon code (must be unique)
	CouponCode string `json:"coupon_code" gorm:"primaryKey"`
	// Expiry date of the coupon
	ExpiryDate time.Time `json:"expiry_date"`
	// Usage type of the coupon (one_time, multi_use, time_based)
	UsageType UsageType `json:"usage_type"`
	// IDs of applicable medicines
	ApplicableMedicineIDs []string `json:"applicable_medicine_ids" gorm:"type:text[]"`
	// Categories the coupon applies to
	ApplicableCategories []string `json:"applicable_categories" gorm:"type:text[]"`
	// Minimum order value required for the coupon
	MinOrderValue float64 `json:"min_order_value"`
	// Valid time window for the coupon (e.g. "10:00-18:00")
	ValidTimeWindow string `json:"valid_time_window"`
	// Terms and conditions of the coupon
	TermsAndConditions string `json:"terms_and_conditions"`
	// Type of discount (percentage, flat)
	DiscountType DiscountType `json:"discount_type"`
	// Value of the discount
	DiscountValue float64 `json:"discount_value"`
	// Maximum usage allowed per user
	MaxUsagePerUser int `json:"max_usage_per_user"`
	// Discount target (inventory, charges)
	DiscountTarget DiscountTarget `json:"discount_target"`
}
