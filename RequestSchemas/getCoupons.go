package requestschemas

import "time"

type ProductInfo struct {
	ID       string `json:"id"`
	Category string `json:"category"`
}

type RequestCoupons struct {
	CartItems []CartItem `json:"cart_items"`
	TimeStamp time.Time  `json:"timestamp"`
}

type CouponsResult struct {
	CouponCode    string   `json:"coupon_code"`
	DiscountValue float64  `json:"discount_value"`
	DiscountType  string   `json:"discount_type"`
	Categories    []string `json:"categories"`
}
