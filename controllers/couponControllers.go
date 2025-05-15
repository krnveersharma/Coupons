package controllers

import (
	"net/http"

	requestschemas "github.com/farmako/RequestSchemas"
	"github.com/farmako/models"
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
	var ids []string
	var categories []string
	var coupons []models.Coupon
	var result []requestschemas.CouponsResult

	if err := ctx.ShouldBindBodyWithJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Please input correct data",
		})
	}

	for i := 0; i < len(requestData.CartItems); i++ {
		ids = append(ids, requestData.CartItems[i].ID)
		categories = append(categories, requestData.CartItems[i].Category)
	}

	query := "SELECT * from coupons where ( applicable_medicine_ids && ? OR applicable_categories && ? ) AND min_order_value <= 700 AND expiry_date >= ?"
	err := c.DB.Raw(query, ids, categories, requestData.OrderTotal, requestData.TimeStamp).Scan(&coupons).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error: " + err.Error(),
		})
		return
	}

	for i := 0; i < len(coupons); i++ {
		result = append(result, requestschemas.CouponsResult{
			CouponCode:    coupons[i].CouponCode,
			DiscountValue: coupons[i].DiscountValue,
			DiscountType:  string(coupons[i].DiscountType),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"applicable_coupons": result,
	})
}
