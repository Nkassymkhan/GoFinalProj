package config

import (
	"github.com/Nkassymkhan/GoFinalProj.git/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	db, err := gorm.Open(postgres.Open("host=localhost dbname=go_store user=postgres password=1234"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.Product{})

	return db
}
