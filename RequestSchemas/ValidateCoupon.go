package requestschemas

import "time"

type ValidaTeCouponData struct {
	CouponCode string        `json:"coupon_code"`
	CartItems  []ProductInfo `json:"cart_items"`
	TimeStamp  time.Time     `json:"timestamp"`
}
