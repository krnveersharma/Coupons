package helpers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	requestschemas "github.com/farmako/RequestSchemas"
	"github.com/farmako/cache"
	dboperations "github.com/farmako/dbOperations"
	"gorm.io/gorm"
)

func GetApplicableIds(cartItems []requestschemas.ProductInfo) []string {
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

func GetPrice(ids []string, db *gorm.DB) (uint, error) {
	var totalPrice uint
	for i := range ids {

		price, err := GetProductPrice(ids[i], db)
		if err != nil {
			return 0, err
		}
		totalPrice += price
	}

	fmt.Printf("total price is: %v", totalPrice)
	return totalPrice, nil
}

func GetProductPrice(id string, db *gorm.DB) (uint, error) {
	priceStr, err := cache.Get(id)
	if err == nil && priceStr != "" {
		priceUint64, err := strconv.ParseUint(priceStr, 10, 32)
		if err == nil {
			return uint(priceUint64), nil
		}
	}

	price, err := dboperations.GetPrice(id, db)
	if err != nil {
		return 0, err
	}

	err = cache.Set(id, strconv.FormatUint(uint64(price), 10))
	if err != nil {
		fmt.Printf("unable to set reis key: %v", err)
	}
	return price, nil
}

func IsWithinValidTimeWindow(ts time.Time, validTimeWindow string) (bool, error) {
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
