package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	requestschemas "github.com/farmako/RequestSchemas"
	"github.com/farmako/models"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ControllerSetup struct {
	DB *gorm.DB
}

func SetupController(db *gorm.DB) ControllerSetup {
	return ControllerSetup{
		DB: db,
	}
}

func (c *ControllerSetup) GetApplicableCoupons(ctx *gin.Context) {
	var requestData requestschemas.RequestCoupons
	var result []requestschemas.CouponsResult

	if err := ctx.ShouldBindBodyWithJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Please input correct data",
		})
	}

	ids := getApplicableIds(requestData.CartItems)
	categories := GetApplicableCategories(requestData.CartItems)

	coupons, err := c.getCouponsFromDb(ids, categories, requestData.OrderTotal, requestData.TimeStamp)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	for i := 0; i < len(coupons); i++ {
		isInValidTime, err := isWithinValidTimeWindow(requestData.TimeStamp, coupons[i].ValidTimeWindow)
		if err == nil && isInValidTime {
			result = append(result, requestschemas.CouponsResult{
				CouponCode:    coupons[i].CouponCode,
				DiscountValue: coupons[i].DiscountValue,
				DiscountType:  string(coupons[i].DiscountType),
			})
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"applicable_coupons": result,
	})
}

func (c *ControllerSetup) ValidateCoupon(ctx *gin.Context) {
	var requestData requestschemas.ValidaTeCouponData
	if err := ctx.ShouldBindBodyWithJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Please input correct data",
		})
	}

	coupon, err := getCouponDataFromDb(requestData.CouponCode, c.DB)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Coupon",
		})
	}

	err = isCouponValid(requestData, coupon)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"isValid": false,
			"discount": gin.H{
				"items_discount":   0,
				"charges_discount": 0,
			},
		})
	}

	getDiscountedPrice := getDiscountedPrice(coupon, requestData.OrderTotal)

	ctx.JSON(http.StatusOK, gin.H{
		"is_valid": true,
		"discount": gin.H{
			"items_discount":   requestData.OrderTotal - getDiscountedPrice,
			"charges_discount": 0,
		},
	})
}

// all helpers are here

func getApplicableIds(cartItems []requestschemas.ProductInfo) []string {
	var ids []string
	for i := 0; i < len(cartItems); i++ {
		ids = append(ids, cartItems[i].ID)
	}
	return ids
}

func GetApplicableCategories(cartItems []requestschemas.ProductInfo) []string {
	var categories []string
	for i := 0; i < len(cartItems); i++ {
		categories = append(categories, cartItems[i].Category)
	}
	return categories
}

func isWithinValidTimeWindow(ts time.Time, validTimeWindow string) (bool, error) {
	if validTimeWindow == "" {
		return true, nil
	}

	parts := strings.Split(validTimeWindow, "-")
	if len(parts) != 2 {
		return false, errors.New("invalid validTimeWindow format")
	}

	layout := "15:04"

	startTime, err := time.Parse(layout, parts[0])
	if err != nil {
		return false, err
	}
	endTime, err := time.Parse(layout, parts[1])
	if err != nil {
		return false, err
	}

	tsTime, err := time.Parse(layout, ts.Format(layout))
	if err != nil {
		return false, err
	}

	if endTime.After(startTime) {
		return (tsTime.Equal(startTime) || tsTime.After(startTime)) && tsTime.Before(endTime), nil
	} else {
		return tsTime.Equal(startTime) || tsTime.After(startTime) || tsTime.Before(endTime), nil
	}
}

func isCouponValid(requestData requestschemas.ValidaTeCouponData, coupon models.Coupon) error {
	var isValid bool

	if coupon.MinOrderValue > float64(requestData.OrderTotal) {
		return fmt.Errorf("Minimum price is %s", coupon.MinOrderValue)
	}
	for i := 0; i < len(requestData.CartItems); i++ {
		if contains(coupon.ApplicableCategories, requestData.CartItems[i].Category) || contains(coupon.ApplicableMedicineIDs, requestData.CartItems[i].ID) {
			isInValidTime, err := isWithinValidTimeWindow(requestData.TimeStamp, coupon.ValidTimeWindow)
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

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// db functions
func (c *ControllerSetup) getCouponsFromDb(ids, categories []string, orderTotal uint, timestamp time.Time) ([]models.Coupon, error) {
	var coupons []models.Coupon
	query := "SELECT * from coupons where ( applicable_medicine_ids && ? OR applicable_categories && ? ) AND min_order_value <= ? AND expiry_date >= ?"
	err := c.DB.Raw(query, pq.Array(ids), pq.Array(categories), orderTotal, timestamp).Scan(&coupons).Error
	if err != nil {
		return []models.Coupon{}, err
	}
	return coupons, nil
}

func getDiscountedPrice(coupon models.Coupon, price uint) uint {
	if coupon.DiscountType == "percentage" {
		return price - price*uint(coupon.DiscountValue)/100
	}
	return price - uint(coupon.DiscountValue)
}

func getCouponDataFromDb(couponCode string, db *gorm.DB) (models.Coupon, error) {
	var coupon models.Coupon
	query := "SELECT * from coupons where coupon_code = ?"
	err := db.Raw(query, couponCode).Scan(&coupon).Error
	if err != nil {
		return models.Coupon{}, err
	}
	return coupon, nil
}
