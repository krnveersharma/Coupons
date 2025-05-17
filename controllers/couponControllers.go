package controllers

import (
	"net/http"

	requestschemas "github.com/farmako/RequestSchemas"
	dboperations "github.com/farmako/dbOperations"
	"github.com/farmako/helpers"
	"github.com/farmako/models"
	"github.com/gin-gonic/gin"
)

// @Summary Add a new coupon
// @Description Adds a coupon to the system
// @Tags Coupon
// @Accept json
// @Produce json
// @Param coupon body models.CouponSwagger true "Coupon data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /coupon/add [post]
func (c *ControllerSetup) AddCoupon(ctx *gin.Context) {
	var requestData models.Coupon
	if err := ctx.ShouldBindBodyWithJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Please input correct data",
		})
	}

	err := dboperations.AddCouponToDb(requestData, c.DB)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(201, gin.H{
		"message": "Successfuly added coupon",
	})
}

// GetApplicableCoupons godoc
// @Summary Get applicable coupons
// @Description Returns applicable coupons based on cart data
// @Tags coupons
// @Accept json
// @Produce json
// @Param request body requestschemas.RequestCoupons true "Request data"
// @Success 200 {object} map[string][]requestschemas.CouponsResult
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /coupon/applicable [post]
func (c *ControllerSetup) GetApplicableCoupons(ctx *gin.Context) {
	var requestData requestschemas.RequestCoupons
	var result []requestschemas.CouponsResult

	if err := ctx.ShouldBindBodyWithJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Please input correct data",
		})
		return
	}

	ids := helpers.GetApplicableIds(requestData.CartItems)
	categories := helpers.GetApplicableCategories(requestData.CartItems)

	totalPrice, err := helpers.GetPrice(ids, c.DB)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
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
				Categories:    coupons[i].ApplicableCategories,
			})
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"applicable_coupons": result,
	})
}

// ValidateCoupon godoc
// @Summary Validate a coupon
// @Description Validates coupon for cart and returns applicable discounts
// @Tags coupons
// @Accept json
// @Produce json
// @Param request body requestschemas.ValidaTeCouponData true "Validation data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /coupon/validate [post]
func (c *ControllerSetup) ValidateCoupon(ctx *gin.Context) {
	var requestData requestschemas.ValidaTeCouponData
	if err := ctx.ShouldBindBodyWithJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Please input correct data",
		})
		return
	}

	coupon, err := dboperations.GetCouponDataFromDb(requestData.CouponCode, c.DB)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Coupon",
		})
		return
	}

	ids := helpers.GetApplicableIds(requestData.CartItems)
	totalPrice, err := helpers.GetPrice(ids, c.DB)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
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
		return
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
