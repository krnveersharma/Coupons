package models

type Product struct {
	Id       string `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"not null"`
	Category string `json:"category" gorm:"not null"`
	Price    uint   `json:"price" gorm:"not null"`
}
