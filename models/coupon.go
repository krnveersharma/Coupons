package models

import (
	"time"
)

type UsageType string
type DiscountType string

const (
	UsageTypeOneTime   UsageType = "one_time"
	UsageTypeMultiUse  UsageType = "multi_use"
	UsageTypeTimeBased UsageType = "time_based"
)
const (
	DiscountPercentage DiscountType = "percentage"
	DiscountFlat       DiscountType = "flat"
)

type Coupon struct {
	ID                    uint         `json:"id" gorm:"primaryKey;autoIncrement"`
	CouponCode            string       `json:"coupon_code" gorm:"unique;not null"`
	ExpiryDate            time.Time    `json:"expiry_date"`
	UsageType             UsageType    `json:"usage_type"`
	ApplicableMedicineIDs []uint       `json:"applicable_medicine_ids" gorm:"-"`
	ApplicableCategories  []string     `json:"applicable_categories" gorm:"-"`
	MinOrderValue         float64      `json:"min_order_value"`
	ValidTimeWindow       string       `json:"valid_time_window"`
	TermsAndConditions    string       `json:"terms_and_conditions"`
	DiscountType          DiscountType `json:"discount_type"`
	DiscountValue         float64      `json:"discount_value"`
	MaxUsagePerUser       int          `json:"max_usage_per_user"`
}
