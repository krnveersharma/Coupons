package requestschemas

import "time"

type ValidaTeCouponData struct {
	CouponCode string        `json:"coupon_code"`
	CartItems  []ProductInfo `json:"cart_items"`
	OrderTotal uint          `json:"order_total"`
	TimeStamp  time.Time     `json:"timestamp"`
}
