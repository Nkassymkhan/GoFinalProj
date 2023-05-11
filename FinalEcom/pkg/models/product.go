package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Cost        float32 `json:"cost"`
	Rating      float32 `json:"rating"`
}

type Read struct {
	Ord string `json:"ord"`
}
type Comment struct {
	gorm.Model
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	ItemID int    `json:"item_id"`
	Text   string `json:"text"`
}

type Purchase struct {
	gorm.Model
	ID     int `json:"id"`
	UserID int `json:"user_id"`
	ItemID int `json:"item_id"`
}
