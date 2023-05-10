package handlers

import (
	"errors"
	"github.com/Nkassymkhan/GoFinalProj.git/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type handler struct {
	DB *gorm.DB
}

func New(db *gorm.DB) handler {
	return handler{db}
}

func (h *handler) GetProducts(c *gin.Context) {

	var ord models.Read
	if err := c.BindJSON(&ord); err != nil {
		c.IndentedJSON(http.StatusOK, "Input is not correct")
		panic(err)
	} else {
		var product []models.Product
		if ord.Ord != "" && ord.Ord == "desc" {
			h.DB.Order("id desc").Find(&product)
		} else {
			h.DB.Order("id asc").Find(&product)
		}
		c.IndentedJSON(http.StatusOK, &product)
		// c.IndentedJSON(http.StatusOK, ord)
	}
}

func (h *handler) Home(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Welcome to the productstore")
}

func (h *handler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	readProduct := models.Product{}

	dbRresult := h.DB.Where("title = ?", id).First(&readProduct)
	if errors.Is(dbRresult.Error, gorm.ErrRecordNotFound) {
		if dbRresult = h.DB.Where("id = ?", id).First(&readProduct); dbRresult.Error != nil {
			c.IndentedJSON(http.StatusOK, "product not found")
		} else {
			c.IndentedJSON(http.StatusOK, &readProduct)
		}
	} else {
		c.IndentedJSON(http.StatusOK, &readProduct)
	}

}

func (h *handler) Createproduct(c *gin.Context) {
	var newproduct models.Product
	if err := c.BindJSON(&newproduct); err != nil {
		c.IndentedJSON(http.StatusOK, "Input is not correct")
		panic(err)
	} else {
		h.DB.Create(&newproduct)
		c.IndentedJSON(http.StatusOK, newproduct)
	}
}

func (h *handler) Updateproduct(c *gin.Context) {
	id := c.Param("id")
	readProduct := &models.Product{}

	if dbRresult := h.DB.Where("id = ?", id).First(&readProduct); dbRresult.Error != nil {
		c.IndentedJSON(http.StatusOK, "product not found")
	} else {

		var newproduct models.Product
		if err := c.BindJSON(&newproduct); err != nil {
			c.IndentedJSON(http.StatusOK, "Input is not correct")
			panic(err)
		} else {
			if newproduct.Cost != 0 {
				readProduct.Cost = newproduct.Cost
			}
			if newproduct.Name != "" {
				readProduct.Name = newproduct.Name
			}
			if newproduct.Description != "" {
				readProduct.Description = newproduct.Description
			}
			h.DB.Save(readProduct)
			c.IndentedJSON(http.StatusOK, readProduct)
		}
	}

}

func (h *handler) Deleteproduct(c *gin.Context) {
	id := c.Param("id")
	var deleteproduct models.Product

	if dbRresult := h.DB.Where("id = ?", id).First(&deleteproduct); dbRresult.Error != nil {
		c.IndentedJSON(http.StatusOK, "product not found")
	} else {
		h.DB.Where("id = ?", id).Delete(&deleteproduct)
		c.IndentedJSON(http.StatusOK, "product deleted")
	}
}
