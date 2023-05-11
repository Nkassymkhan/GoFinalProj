package config

import (
	"github.com/Nkassymkhan/GoFinalProj.git/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	db, err := gorm.Open(postgres.Open("host=localhost dbname=store_db user=postgres password=xzsawq21"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.Product{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Comment{})
	db.AutoMigrate(&models.Purchase{})
	db.AutoMigrate(&models.Rating{})
	return db
}
