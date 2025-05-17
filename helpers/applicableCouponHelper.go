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

type IdAndPrice struct {
	Id       string
	Quantity uint
}

type CategoryAndPrice struct {
	Category string
	Quantity uint
}

func GetApplicableIds(cartItems []requestschemas.CartItem) []string {
	var ids []string
	for i := 0; i < len(cartItems); i++ {
		ids = append(ids, cartItems[i].ID)
	}
	return ids
}

func GetApplicableIdsFromValidate(cartItems []requestschemas.CartItem) []IdAndPrice {
	var ids []IdAndPrice
	for i := 0; i < len(cartItems); i++ {
		ids = append(ids, IdAndPrice{Id: cartItems[i].ID, Quantity: cartItems[i].Quantity})
	}
	return ids
}

func GetApplicableCategoriessFromValidate(cartItems []requestschemas.CartItem) []CategoryAndPrice {
	var categories []CategoryAndPrice
	for i := 0; i < len(cartItems); i++ {
		categories = append(categories, CategoryAndPrice{Category: cartItems[i].Category, Quantity: cartItems[i].Quantity})
	}
	return categories
}

func GetApplicableCategories(cartItems []requestschemas.CartItem) []string {
	var categories []string
	for i := 0; i < len(cartItems); i++ {
		categories = append(categories, cartItems[i].Category)
	}
	return categories
}

func GetPricePerProduct(ids []string, db *gorm.DB) (uint, error) {
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

func GetPrice(ids []IdAndPrice, db *gorm.DB) (uint, error) {
	var totalPrice uint
	for i := range ids {

		price, err := GetProductPrice(ids[i].Id, db)
		if err != nil {
			return 0, err
		}
		totalPrice += price * ids[i].Quantity
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
