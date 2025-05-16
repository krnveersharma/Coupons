package api

import (
	"fmt"
	"log"
	"time"

	"github.com/farmako/cache"
	"github.com/farmako/config"
	"github.com/farmako/controllers"
	"github.com/farmako/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func StartServer(config config.AppConfig) {
	app := gin.Default()

	fmt.Printf("config is: %v", config)

	db, err := gorm.Open(postgres.Open(config.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error in connecting DB: %v", err.Error())
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	db.AutoMigrate(&models.Coupon{})
	db.AutoMigrate(&models.Product{})

	setUpRoutes(app, db)
	cache.SetupRedis(config.RedisAddress, config.RedisPassword, 3600*time.Second)

	app.Run(fmt.Sprintf(":%v", config.ServerPort))
}

func setUpRoutes(app *gin.Engine, db *gorm.DB) {
	setupController := controllers.SetupController(db)
	app.POST("/coupons/applicable", setupController.GetApplicableCoupons)
	app.POST("/coupons/validate", setupController.ValidateCoupon)
}
