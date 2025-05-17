package requestschemas

import "time"

type CartItem struct {
	ID       string `json:"id"`
	Category string `json:"category"`
	Quantity uint   `json:"quantity"`
}

type ValidaTeCouponData struct {
	CouponCode string     `json:"coupon_code"`
	CartItems  []CartItem `json:"cart_items"`
	TimeStamp  time.Time  `json:"timestamp"`
}
