package controllers

import "gorm.io/gorm"

type ControllerSetup struct {
	DB *gorm.DB
}

func SetupController(db *gorm.DB) ControllerSetup {
	return ControllerSetup{
		DB: db,
	}
}
