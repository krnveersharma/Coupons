package requestschemas

import "time"

type UsageType string
type DiscountType string
type DiscountTarget string

const (
	UsageTypeOneTime   UsageType = "one_time"
	UsageTypeMultiUse  UsageType = "multi_use"
	UsageTypeTimeBased UsageType = "time_based"
)
const (
	DiscountPercentage DiscountType = "percentage"
	DiscountFlat       DiscountType = "flat"
)

const (
	DiscountInventory DiscountTarget = "inventory"
	DiscountCharges   DiscountTarget = "charges"
)

// CouponSwagger is a Swagger-compatible representation of Coupon
// swagger:model
type CouponSwagger struct {
	// Coupon code (must be unique)
	CouponCode string `json:"coupon_code"`
	// Expiry date of the coupon
	ExpiryDate time.Time `json:"expiry_date"`
	// Usage type of the coupon
	UsageType UsageType `json:"usage_type"`
	// IDs of applicable medicines
	ApplicableMedicineIDs []string `json:"applicable_medicine_ids"`
	// Applicable product categories
	ApplicableCategories []string `json:"applicable_categories"`
	// Minimum order value
	MinOrderValue float64 `json:"min_order_value"`
	// Valid time window (e.g., "09:00-18:00")
	ValidTimeWindow string `json:"valid_time_window"`
	// Terms and conditions
	TermsAndConditions string `json:"terms_and_conditions"`
	// Discount type (percentage or flat)
	DiscountType DiscountType `json:"discount_type"`
	// Discount value
	DiscountValue float64 `json:"discount_value"`
	// Max usage per user
	MaxUsagePerUser int `json:"max_usage_per_user"`
	// Target of discount
	DiscountTarget DiscountTarget `json:"discount_target"`
}
