package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/Nkassymkhan/GoFinalProj.git/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func New(db *gorm.DB) handler {
	return handler{db}
}

func (h *handler) Home(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Welcome to the Product store")
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
	}
}

func (h *handler) GetSortedProductsByCost(c *gin.Context) {
	var products []models.Product
	sort := c.Query("sort")
	parts := strings.Split(sort, "-")
	if parts[0] == "cost" {
		parts[0] = "price"
	}
	sorting := strings.Join(parts, " ")
	if sorting == "" {
		sorting = "cost asc"
	}
	if err := h.DB.Order(sorting).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

func (h *handler) GetSortedProductsByRating(c *gin.Context) {
	var products []models.Product
	sort := c.Query("sort")
	parts := strings.Split(sort, "-")

	sorting := strings.Join(parts, " ")
	if sorting == "" {
		sorting = "rating asc"
	}
	if err := h.DB.Order(sorting).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

func (h *handler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	readProduct := models.Product{}

	dbRresult := h.DB.Where("name = ?", id).First(&readProduct)
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

func (h *handler) CreateProduct(c *gin.Context) {
	var newproduct models.Product
	if err := c.BindJSON(&newproduct); err != nil {
		c.IndentedJSON(http.StatusOK, "Input is not correct")
		panic(err)
	} else {
		h.DB.Create(&newproduct)
		c.IndentedJSON(http.StatusOK, newproduct)
	}
}

func (h *handler) GiveRating(c *gin.Context) {
	id := c.Param("id")
	readProduct := &models.Product{}

	if dbResult := h.DB.Where("id = ?", id).First(&readProduct); dbResult.Error != nil {
		c.IndentedJSON(http.StatusNotFound, "product not found")
		return
	}

	var newProduct models.Product
	if err := c.BindJSON(&newProduct); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Input is not correct")
		log.Printf("Error binding JSON: %v", err)
		return
	}

	if newProduct.Rating != 0 {
		readProduct.Rating = newProduct.Rating
	}
	h.DB.Save(readProduct)
	c.IndentedJSON(http.StatusOK, readProduct)
}

func (h *handler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	var deleteproduct models.Product

	if dbRresult := h.DB.Where("id = ?", id).First(&deleteproduct); dbRresult.Error != nil {
		c.IndentedJSON(http.StatusOK, "product not found")
	} else {
		h.DB.Where("id = ?", id).Delete(&deleteproduct)
		c.IndentedJSON(http.StatusOK, "product deleted")
	}
}

func (h *handler) CommentItem(c *gin.Context) {
	var comment models.Comment
	c.BindJSON(&comment)

	if comment.UserID == 0 || comment.ItemID == 0 || comment.Text == "" {
		c.JSON(400, gin.H{"error": "Invalid comment data"})
		return
	}

	if err := h.DB.Create(&comment).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(200, gin.H{"message": "Comment created successfully"})
}

func (h *handler) PurchaseItem(c *gin.Context) {
	var purchase models.Purchase
	c.BindJSON(&purchase)

	if purchase.UserID == 0 || purchase.ItemID == 0 {
		c.JSON(400, gin.H{"error": "Invalid purchase data"})
		return
	}

	if err := h.DB.Create(&purchase).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create purchase"})
		return
	}

	c.JSON(200, gin.H{"message": "Item purchased successfully"})
}
