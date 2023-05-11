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

	return db
}
