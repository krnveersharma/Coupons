package controllers

import (
	"net/http"

	requestschemas "github.com/farmako/RequestSchemas"
	dboperations "github.com/farmako/dbOperations"
	"github.com/farmako/helpers"
	"github.com/gin-gonic/gin"
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

	ids := helpers.GetApplicableIds(requestData.CartItems)
	categories := helpers.GetApplicableCategories(requestData.CartItems)

	totalPrice, err := helpers.GetPrice(ids, c.DB)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

	coupons, err := dboperations.GetCouponsFromDb(ids, categories, totalPrice, requestData.TimeStamp, c.DB)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	for i := 0; i < len(coupons); i++ {
		isInValidTime, err := helpers.IsWithinValidTimeWindow(requestData.TimeStamp, coupons[i].ValidTimeWindow)
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

	coupon, err := dboperations.GetCouponDataFromDb(requestData.CouponCode, c.DB)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Coupon",
		})
	}

	ids := helpers.GetApplicableIds(requestData.CartItems)
	totalPrice, err := helpers.GetPrice(ids, c.DB)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

	err = helpers.IsCouponValid(requestData, coupon, totalPrice)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"isValid": false,
			"discount": gin.H{
				"items_discount":   0,
				"charges_discount": 0,
			},
		})
	}

	itemsDiscount, chargesDiscount := helpers.GetDiscountedPrice(coupon, totalPrice)

	ctx.JSON(http.StatusOK, gin.H{
		"is_valid": true,
		"discount": gin.H{
			"items_discount":   itemsDiscount,
			"charges_discount": chargesDiscount,
		},
	})
}
