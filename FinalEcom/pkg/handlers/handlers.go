package handlers

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Nkassymkhan/GoFinalProj.git/pkg/models"
	"github.com/dgrijalva/jwt-go"
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
	userID, err := h.getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var rating models.Rating
	if err := c.ShouldBindJSON(&rating); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rating data"})
		return
	}

	rating.UserID = int(userID)
	rating.ItemID = int(itemID)
	if err := h.DB.Where(models.Rating{UserID: rating.UserID, ItemID: rating.ItemID}).Assign(rating).FirstOrCreate(&rating).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create or update rating"})
		return
	}

	var avgRating float64
	if err := h.DB.Model(&models.Rating{}).Where("item_id = ?", itemID).Select("AVG(rating)").Row().Scan(&avgRating); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate average rating"})
		return
	}

	if err := h.DB.Model(&models.Product{}).Where("id = ?", itemID).Update("rating", avgRating).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ratings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ratings updated successfully"})
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

	userID, err := h.getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment data"})
		return
	}

	comment.UserID = int(userID)

	if err := h.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment created successfully"})
}

func (h *handler) getUserIDFromToken(c *gin.Context) (int, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return 0, errors.New("Authorization header missing")
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return 0, errors.New("Invalid token")
	}
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return 0, errors.New("Invalid token claims")
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return 0, errors.New("Invalid user ID")
	}

	return userID, nil
}

func (h *handler) PurchaseItem(c *gin.Context) {
	userID, err := h.getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	purchase := models.Purchase{
		UserID: int(userID),
		ItemID: int(itemID),
	}

	if err := h.DB.Create(&purchase).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create purchase"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item purchased successfully"})
}
