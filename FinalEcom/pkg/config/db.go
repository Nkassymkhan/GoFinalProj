package config

import (
	"github.com/Nkassymkhan/GoFinalProj.git/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	db, err := gorm.Open(postgres.Open("host=5432 dbname=postgres user=postgres password=postgres"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.Product{})

	return db
}
