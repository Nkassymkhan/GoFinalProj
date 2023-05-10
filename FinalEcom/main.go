package main

import (
	"github.com/Nkassymkhan/GoFinalProj.git/pkg/config"
	"github.com/Nkassymkhan/GoFinalProj.git/pkg/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	db := config.Connect()
	h := handlers.New(db)
	router := gin.Default()
	router.GET("/", h.Home)
	router.POST("/books", h.GetProducts)
	router.GET("/book/:id", h.GetProduct)
	router.POST("/book", h.Createproduct)
	router.DELETE("/book/:id", h.Deleteproduct)
	router.PUT("/book/:id", h.Updateproduct)

	router.Run(":8080")

}
