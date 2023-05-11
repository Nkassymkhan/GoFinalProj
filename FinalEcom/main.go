package main

import (
	"github.com/Nkassymkhan/GoFinalProj.git/pkg/authorization"
	"github.com/Nkassymkhan/GoFinalProj.git/pkg/config"
	"github.com/Nkassymkhan/GoFinalProj.git/pkg/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	db := config.Connect()
	h := handlers.New(db)
	a := authorization.New(db)
	router := gin.Default()
	router.GET("/", h.Home)
	router.POST("/products", h.GetProducts)
	router.GET("/product/:id", h.GetProduct)
	router.POST("/product", h.CreateProduct)
	router.DELETE("/product/:id", h.DeleteProduct)
	router.POST("/product/:id/rating", h.GiveRating)
	router.GET("/products/sorted", h.GetSortedProductsByCost)
	router.GET("/products/sortedRatings", h.GetSortedProductsByRating)
	router.POST("/register", a.Register)
	router.POST("/products/:id/comment", h.CommentItem)
	router.POST("/products/:id/purchase", h.PurchaseItem)
	router.POST("/login", a.Login)

	router.Run(":8080")

}
