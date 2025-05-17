package controllers

import (
	"net/http"

	dboperations "github.com/farmako/dbOperations"
	"github.com/farmako/models"
	"github.com/gin-gonic/gin"
)

// AddProduct godoc
// @Summary Add a new product
// @Description Add a product with id, name, category and price
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.Product true "Product Info"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /products [post]
func (c *ControllerSetup) GetProducts(ctx *gin.Context) {
	products, err := dboperations.GetProductsFromDb(c.DB)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}

// GetProducts godoc
// @Summary Get all products
// @Description Fetch all products from the database
// @Tags products
// @Produce json
// @Success 200 {object} map[string][]models.Product
// @Failure 500 {object} map[string]interface{}
// @Router /products [get]
func (c *ControllerSetup) AddProduct(ctx *gin.Context) {
	var requestData models.Product
	if err := ctx.ShouldBindBodyWithJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Please input correct data",
		})
	}

	err := dboperations.AddProduct(requestData, c.DB)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "product with same id already exists",
		})
		return
	}
	ctx.JSON(201, gin.H{
		"message": "Successfuly added product",
	})
}
