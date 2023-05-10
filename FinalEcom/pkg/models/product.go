package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Cost        float32 `json:"cost"`
}

type Read struct {
	Ord string `json:"ord"`
}
