package models

import (
	"time"

	"github.com/lib/pq"
)

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

type Coupon struct {
	CouponCode            string         `json:"coupon_code" gorm:"primaryKey"`
	ExpiryDate            time.Time      `json:"expiry_date"`
	UsageType             UsageType      `json:"usage_type"`
	ApplicableMedicineIDs pq.StringArray `json:"applicable_medicine_ids" gorm:"type:text[]"`
	ApplicableCategories  pq.StringArray `json:"applicable_categories" gorm:"type:text[]"`
	MinOrderValue         float64        `json:"min_order_value"`
	ValidTimeWindow       string         `json:"valid_time_window"`
	TermsAndConditions    string         `json:"terms_and_conditions"`
	DiscountType          DiscountType   `json:"discount_type"`
	DiscountValue         float64        `json:"discount_value"`
	MaxUsagePerUser       int            `json:"max_usage_per_user"`
	DiscountTarget        DiscountTarget `json:"discount_target"`
}
